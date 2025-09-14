package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/config"
	"chama-wallet-backend/handlers"
	"chama-wallet-backend/middleware"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🚀 Community Wallet API is running",
			"network": config.Config.Network,
			"version": "1.0.0",
		})
	})

	