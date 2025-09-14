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

	// Network info endpoint
	app.Get("/network", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"network":            config.Config.Network,
			"horizon_url":        config.Config.HorizonURL,
			"soroban_rpc_url":    config.Config.SorobanRPCURL,
			"network_passphrase": config.Config.NetworkPassphrase,
			"contract_id":        config.Config.ContractID,
			"is_mainnet":         config.Config.IsMainnet,
			"supported_assets":   config.GetAssetInfo(),
		})
	})
	