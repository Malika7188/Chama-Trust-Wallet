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
	if config.Config.IsMainnet {
		return fmt.Errorf("funding not available on mainnet - use real XLM deposits")
	}

	url := fmt.Sprintf("https://friendbot.stellar.org/?addr=%s", address)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Friendbot response for %s: %s\n", address, string(body))
	return nil
}

// SendXLM transfers XLM from sender to receiver
func SendXLM(seed, destination, amount string) (horizon.Transaction, error) {
	client := GetHorizonClient()

	// Load source account
	kp, err := keypair.ParseFull(seed)
	if err != nil {
		return horizon.Transaction{}, err
	}

	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return horizon.Transaction{}, err
	}

	