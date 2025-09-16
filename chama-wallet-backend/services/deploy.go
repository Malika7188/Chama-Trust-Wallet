package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"chama-wallet-backend/config"
)


func DeployChamaContract() (string, error) {
	if config.Config.IsMainnet {
		return "", fmt.Errorf("contract deployment should be done manually on mainnet for security. Use the configured SOROBAN_CONTRACT_ID instead")
	}

	// Load keys from environment
	source := os.Getenv("SOROBAN_PUBLIC_KEY")
	secret := os.Getenv("SOROBAN_SECRET_KEY")

	if source == "" || secret == "" {
		// Fallback to default test account
		source = "malika"
		secret = os.Getenv("SOROBAN_SECRET_KEY")
		if secret == "" {
			return "", fmt.Errorf("missing SOROBAN_SECRET_KEY in environment")
		}
	}

	