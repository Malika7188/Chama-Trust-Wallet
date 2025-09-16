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

	// Check if WASM file exists
	wasmPath := "./chama_savings/target/wasm32-unknown-unknown/release/chama_savings.wasm"
	if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
		// Try alternative path
		wasmPath = "./chama_savings.wasm"
		if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
			return "", fmt.Errorf("WASM file not found. Please build the contract first with: cd chama_savings && stellar contract build")
		}
	}

	