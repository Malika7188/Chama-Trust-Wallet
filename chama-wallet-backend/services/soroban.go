// services/soroban.go
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

type SorobanInvokeRequest struct {
	ContractID string        `json:"contract_id"`
	Function   string        `json:"function"`
	Args       []interface{} `json:"args"`
}

// validateContractID ensures the contract ID is valid
func validateContractID(contractID string) error {
	if contractID == "" {
		return fmt.Errorf("contract ID cannot be empty")
	}
	if len(contractID) != 56 {
		return fmt.Errorf("invalid contract ID length: expected 56 characters, got %d", len(contractID))
	}
	if !strings.HasPrefix(contractID, "C") {
		return fmt.Errorf("contract ID must start with 'C'")
	}
	return nil
}
