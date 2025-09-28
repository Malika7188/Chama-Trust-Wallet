package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"

	"chama-wallet-backend/config"
)

// GetHorizonClient returns the appropriate Horizon client based on network configuration
func GetHorizonClient() *horizonclient.Client {
	return config.GetHorizonClient()
}

// CreateWallet generates a new Stellar keypair
func CreateWallet() (string, string) {
	kp, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	return kp.Address(), kp.Seed()
}

// FundWallet funds a wallet using Friendbot (testnet only)
func FundWallet(address string) error {
	