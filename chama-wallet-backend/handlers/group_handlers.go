package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
	"chama-wallet-backend/config"
)

func ContributeToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	