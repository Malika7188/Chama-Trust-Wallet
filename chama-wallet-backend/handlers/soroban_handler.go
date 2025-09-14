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

	// Validate required fields
	if body.ContractID == "" || body.UserAddress == "" || body.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: contract_id, user_address, and amount are required",
		})
	}

	// Validate amount
	amount, err := strconv.ParseFloat(body.Amount, 64)
	if err != nil || amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be a positive number",
		})
	}

	// Validate amount limits for mainnet
	if config.Config.IsMainnet {
		minAmount := 0.0000001 // Minimum XLM amount
		if amount < minAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount below minimum of %f XLM for mainnet", minAmount),
			})
		}
	}

	