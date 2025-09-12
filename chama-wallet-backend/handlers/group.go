package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/config"
	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
	"chama-wallet-backend/utils"
)

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateGroup(c *fiber.Ctx) error {
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Get authenticated user
	user := c.Locals("user").(models.User)

	// Generate wallet for the group
	wallet, err := utils.GenerateStellarWallet()
	if err != nil {
		fmt.Printf("❌ Failed to generate group wallet: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate wallet"})
	}

	fmt.Printf("✅ Generated group wallet: %s on %s\n", wallet.PublicKey, config.Config.Network)

	// Fund wallet only on testnet
	if !config.Config.IsMainnet {
		err = services.FundTestAccount(wallet.PublicKey)
		if err != nil {
			fmt.Printf("⚠️ Warning: Failed to fund group wallet: %v\n", err)
			// Don't fail the group creation, just log the warning
		} else {
			fmt.Printf("✅ Group wallet funded successfully on testnet\n")
		}
	}

	// Get contract ID from configuration
	contractID := config.Config.ContractID
	if contractID == "" {
		if config.Config.IsMainnet {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Mainnet contract ID not configured. Please set SOROBAN_CONTRACT_ID in environment.",
			})
		} else {
			// Deploy contract for testnet
			contractID, err = services.DeployChamaContract()
			if err != nil {
				fmt.Printf("❌ Failed to deploy contract: %v\n", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deploy contract"})
			}
		}
	}

	fmt.Printf("✅ Using contract ID: %s on %s\n", contractID, config.Config.Network)

	// Save group in DB with the contract ID
	group := models.Group{
		ID:          uuid.NewString(),
		Name:        payload.Name,
		Description: payload.Description,
		Wallet:      wallet.PublicKey,
		CreatorID:   user.ID,
		ContractID:  contractID,
		Status:      "pending",
		SecretKey:   wallet.SecretKey,
	}

	if err := database.DB.Create(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Add the creator as the first member with creator role and approved status
	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  group.ID,
		UserID:   user.ID,
		Wallet:   user.Wallet,
		Role:     "creator",
		Status:   "approved",
		JoinedAt: time.Now(),
	}

	err = database.DB.Create(&member).Error
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to add creator as member: %v\n", err)
		// Don't fail the group creation
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Group created successfully",
		"group": fiber.Map{
			"id":          group.ID,
			"name":        group.Name,
			"description": group.Description,
			"wallet":      group.Wallet,
			"secret_key":  group.SecretKey,
			"status":      group.Status,
			"contract_id": contractID,
			"network":     config.Config.Network,
		},
	})
}

func AddMember(c *fiber.Ctx) error {
	groupID := c.Params("id")

	var body struct {
		Wallet string `json:"wallet"`
		UserID string `json:"user_id"` // Add this field
	}

	if err := c.BodyParser(&body); err != nil || body.Wallet == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request. Wallet is required.",
		})
	}

	group, err := services.AddMemberToGroup(groupID, body.UserID, body.Wallet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(group)
}
func DepositToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")

	var body struct {
		FromWallet string `json:"from_wallet"`
		Secret     string `json:"secret"` // sender's secret key
		Amount     string `json:"amount"` // XLM to deposit
	}

	if err := c.BodyParser(&body); err != nil || body.FromWallet == "" || body.Secret == "" || body.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields.",
		})
	}

	