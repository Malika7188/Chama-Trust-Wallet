package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

func NominateAdmin(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		NomineeID string `json:"nominee_id"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Check if nominator is a member
	var nominator models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&nominator).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not a group member"})
	}

	// Check if nominee is a member
	var nominee models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, payload.NomineeID, "approved").First(&nominee).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nominee is not a group member"})
	}

	// Check if already nominated
	var existing models.AdminNomination
	if database.DB.Where("group_id = ? AND nominee_id = ? AND status = ?",
		groupID, payload.NomineeID, "pending").First(&existing).Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already nominated"})
	}

	nomination := models.AdminNomination{
		ID:          uuid.NewString(),
		GroupID:     groupID,
		NominatorID: user.ID,
		NomineeID:   payload.NomineeID,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&nomination).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if nominee has 2 nominations, auto-approve as admin
	var nominationCount int64
	database.DB.Model(&models.AdminNomination{}).
		Where("group_id = ? AND nominee_id = ? AND status = ?", groupID, payload.NomineeID, "pending").
		Count(&nominationCount)

	