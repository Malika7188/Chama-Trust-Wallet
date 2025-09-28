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
