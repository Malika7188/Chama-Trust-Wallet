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
