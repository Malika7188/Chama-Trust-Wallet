package services

import (
	"fmt"
	"time"

	"github.com/stellar/go/clients/horizonclient"

	"chama-wallet-backend/config"
)

func CheckBalance(address string) (string, error) {
	client := config.GetHorizonClient()

	