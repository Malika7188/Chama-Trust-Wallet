package services

import (
	"fmt"
	"time"

	"github.com/stellar/go/clients/horizonclient"

	"chama-wallet-backend/config"
)

func CheckBalance(address string) (string, error) {
	client := config.GetHorizonClient()

	// First try to get account details
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: address})
	if err != nil {
		// Check if it's a "Resource Missing" error (account doesn't exist)
		if horizonError, ok := err.(*horizonclient.Error); ok {
			if horizonError.Problem.Status == 404 {
				if config.Config.IsMainnet {
					return "0", fmt.Errorf("account not found on mainnet - account needs to be funded with real XLM first")
				}

				fmt.Printf("⚠️ Account %s not found on testnet. Attempting to fund...\n", address)

				// Try to fund the account
				if fundErr := FundTestAccount(address); fundErr != nil {
					return "0", fmt.Errorf("account not found and funding failed: %w", fundErr)
				}
