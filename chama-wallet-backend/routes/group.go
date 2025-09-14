// routes/group.go
package routes

import (
	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/handlers"
	"chama-wallet-backend/middleware"
)

func GroupRoutes(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
