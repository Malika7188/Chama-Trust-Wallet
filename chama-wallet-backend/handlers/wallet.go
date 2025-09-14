package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"

	"chama-wallet-backend/config"
	"chama-wallet-backend/services"
)

// wallet handlers
// Creates and returns a wallet
func CreateWallet(c *fiber.Ctx) error {
	address, seed := services.CreateWallet()
	return c.JSON(fiber.Map{
		"address": address,
		"seed":    seed,
		"network": config.Config.Network,
	})
}

func GetBalance(c *fiber.Ctx) error {
	address := c.Params("address")

	client := config.GetHorizonClient()
	accountRequest := horizonclient.AccountRequest{AccountID: address}
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		// For mainnet, provide more helpful error messages
		if config.Config.IsMainnet {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Account not found on mainnet. Please ensure the account is funded with real XLM first.",
				"network": config.Config.Network,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var balances []string
	for _, b := range account.Balances {
		assetInfo := "XLM"
		if b.Asset.Type != "native" {
			assetInfo = fmt.Sprintf("%s:%s", b.Asset.Code, b.Asset.Issuer)
		}
		balances = append(balances, fmt.Sprintf("%s: %s", assetInfo, b.Balance))
	}

	