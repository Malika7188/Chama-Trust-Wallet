package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/utils"
)

var groups = make(map[string]models.Group)

func CreateGroup(name, description, creatorID string) (models.Group, error) {
	wallet, err := utils.GenerateStellarWallet()
	if err != nil {
		return models.Group{}, err
	}

	group := models.Group{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Wallet:      wallet.PublicKey,
		CreatorID:   creatorID,
		Status:      "pending",
	}

	