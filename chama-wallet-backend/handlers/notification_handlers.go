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

	