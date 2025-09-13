package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stellar/go/keypair"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

func ContributeToRound(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Round  int     `json:"round"`
		Amount float64 `json:"amount"`
		Secret string  `json:"secret"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// ✅ Validate the secret key belongs to the user
	kp, err := keypair.ParseFull(payload.Secret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid secret key format",
		})
	}

	if kp.Address() != user.Wallet {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Secret key does not match your wallet address",
		})
	}

	// Verify user is a member of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not a member of this group"})
	}

	// Get group details
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	if group.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Group is not active"})
	}

	