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

	fmt.Printf("🔧 Deploying contract from WASM: %s on %s\n", wasmPath, config.Config.Network)
	fmt.Printf("🔧 Using source account: %s\n", source)

	network := config.GetSorobanNetwork()

	// Deploy using source account name (should be configured in soroban keys)
	cmd := exec.Command("soroban",
		"contract", "deploy",
		"--wasm", wasmPath,
		"--source-account", source,
		"--network", network,
	)

	fmt.Printf("🚀 Running deployment command on %s...\n", network)

	// Capture stdout and stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute
	execErr := cmd.Run()
	if execErr != nil {
		fmt.Printf("❌ Deployment error: %v\n", execErr)
		fmt.Printf("❗ stderr: %s\n", stderr.String())
		fmt.Printf("❗ stdout: %s\n", out.String())

		// Try alternative method with temporary key storage
		return deployWithKeyStorage(source, secret)
	}

	output := strings.TrimSpace(out.String())
	fmt.Printf("✅ Contract deployed successfully on %s. Output: %s\n", network, output)

	// Extract contract address from output (usually the last line)
	lines := strings.Split(output, "\n")
	contractAddress := strings.TrimSpace(lines[len(lines)-1])

	