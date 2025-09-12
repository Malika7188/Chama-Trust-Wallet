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

func GetNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	fmt.Printf("🔍 Getting notifications for user: %s\n", user.ID)

	notifications, err := services.GetUserNotifications(user.ID)
	if err != nil {
		fmt.Printf("❌ Error getting notifications: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("✅ Found %d notifications for user %s\n", len(notifications), user.ID)
	return c.JSON(notifications)
}

func MarkNotificationRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Verify notification belongs to user
	var notification models.Notification
	if err := database.DB.Where("id = ? AND user_id = ?", notificationID, user.ID).First(&notification).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Notification not found"})
	}

	