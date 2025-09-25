package services

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"chama-wallet-backend/config"
)

// FundTestAccount uses the Stellar Friendbot to send test XLM to a new account (testnet only)
func FundTestAccount(address string) error {
	if config.Config.IsMainnet {
		return fmt.Errorf("friendbot funding not available on mainnet - use real XLM deposits")
	}

	url := fmt.Sprintf("https://friendbot.stellar.org/?addr=%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to call friendbot: %w", err)
	}
	defer resp.Body.Close()

	
