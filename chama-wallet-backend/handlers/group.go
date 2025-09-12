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

	