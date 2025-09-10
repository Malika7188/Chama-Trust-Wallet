package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
	"chama-wallet-backend/config"
)

func ContributeToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		From   string `json:"from"`
		Secret string `json:"secret"`
		Amount string `json:"amount"`
	}
	
	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("❌ Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if payload.From == "" || payload.Secret == "" || payload.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: from, secret, and amount are required",
		})
	}

	// Validate amount limits for mainnet
	if config.Config.IsMainnet {
		amount, err := strconv.ParseFloat(payload.Amount, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid amount format",
			})
		}
		
		minAmount := 0.0000001 // Minimum XLM amount
		if amount < minAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount below minimum of %f XLM", minAmount),
			})
		}
	}
	// Verify user is a member of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not an approved member of this group",
		})
	}

	// Get group details
	group, err := services.GetGroupByID(groupID)
	if err != nil {
		fmt.Printf("❌ Group not found: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	// Validate group status
	if group.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Group is not active (current status: %s)", group.Status),
		})
	}

	// Validate contract ID
	if group.ContractID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Group contract ID is not set",
		})
	}

	// Validate user's wallet address matches the 'from' field
	if payload.From != user.Wallet {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "From address must match your wallet address",
		})
	}

	fmt.Printf("🔄 Processing contribution: %s XLM from %s to group %s (contract: %s) on %s\n", 
		payload.Amount, payload.From, group.Name, group.ContractID, config.Config.Network)

	