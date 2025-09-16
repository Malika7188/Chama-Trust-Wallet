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
		