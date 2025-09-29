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
	