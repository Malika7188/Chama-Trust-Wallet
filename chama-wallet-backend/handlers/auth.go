package handlers

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
