package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/services"
)

// AuthMiddleware validates JWT tokens and loads the user into context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		