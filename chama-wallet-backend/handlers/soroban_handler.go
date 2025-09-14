package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"chama-wallet-backend/services"
	"chama-wallet-backend/config"
)

// ContributeHandler handles direct Soroban contributions
func ContributeHandler(c *fiber.Ctx) error {
	type RequestBody struct {
		ContractID  string `json:"contract_id"`
		UserAddress string `json:"user_address"`
		Amount      string `json:"amount"`
		SecretKey   string `json:"secret_key,omitempty"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		fmt.Printf("❌ Failed to parse contribute request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
