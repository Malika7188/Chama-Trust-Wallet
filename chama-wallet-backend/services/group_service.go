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

	if err := database.DB.Create(&group).Error; err != nil {
		return models.Group{}, err
	}

	// Automatically add creator as admin
	creator := models.Member{
		ID:       uuid.NewString(),
		GroupID:  group.ID,
		UserID:   creatorID,
		Role:     "creator",
		Status:   "approved",
		JoinedAt: time.Now(),
	}
	database.DB.Create(&creator)

	return group, nil
}

func GetGroupByID(groupID string) (models.Group, error) {
	var group models.Group
	err := database.DB.Preload("Members.User").Preload("Creator").First(&group, "id = ?", groupID).Error
	return group, err
}
