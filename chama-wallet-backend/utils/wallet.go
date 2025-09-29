package utils

import (
	"fmt"

	"github.com/stellar/go/keypair"
)

type StellarWallet struct {
	PublicKey string
	SecretKey string
}

// GenerateStellarWallet creates a new Stellar keypair
func GenerateStellarWallet() (*StellarWallet, error) {
	// Generate a new keypair
	pair, err := keypair.Random()
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	// Validate the generated keys
	if pair.Address() == "" {
	