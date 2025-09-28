package services

import (
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"

	"chama-wallet-backend/config"
)

func SendPayment(fromSecret, toAddress, amount string) error {
	senderKP, err := keypair.ParseFull(fromSecret)
	if err != nil {
		return fmt.Errorf("invalid secret key: %w", err)
	}

	client := config.GetHorizonClient()
	ar := horizonclient.AccountRequest{AccountID: senderKP.Address()}
