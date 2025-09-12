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

	