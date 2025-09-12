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

	if err := services.MarkNotificationAsRead(notificationID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Notification marked as read"})
}

func GetUserInvitations(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	fmt.Printf("🔍 Getting invitations for user: %s (email: %s)\n", user.ID, user.Email)

	var invitations []models.GroupInvitation
	err := database.DB.Where("email = ? AND status = ?", user.Email, "pending").
		Preload("Group").
		Preload("Inviter").
		Find(&invitations).Error

	if err != nil {
		fmt.Printf("❌ Error getting invitations: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("✅ Found %d invitations for user %s\n", len(invitations), user.Email)
	return c.JSON(invitations)
}

// DeleteNotification deletes a notification by ID for the authenticated user
func DeleteNotification(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	user := c.Locals("user").(models.User)

	