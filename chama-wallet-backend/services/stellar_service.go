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
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return fmt.Errorf("could not load source account: %w", err)
	}

	op := txnbuild.Payment{
		Destination: toAddress,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	// Add memo for mainnet compliance
	var memo txnbuild.Memo
	