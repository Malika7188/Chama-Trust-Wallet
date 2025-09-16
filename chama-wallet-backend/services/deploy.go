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
