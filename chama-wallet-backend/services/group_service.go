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

func AddMemberToGroup(groupID, userID, walletAddress string) (models.Group, error) {
	var group models.Group
	if err := database.DB.Preload("Members").First(&group, "id = ?", groupID).Error; err != nil {
		return group, err
	}

	// Check if member already exists
	for _, member := range group.Members {
		if member.UserID == userID {
			return group, nil // Member already exists
		}
	}

	member := models.Member{
		ID:       uuid.NewString(),
		GroupID:  groupID,
		UserID:   userID,
		Wallet:   walletAddress,
		Role:     "member",
		Status:   "approved",
		JoinedAt: time.Now(),
	}
	if err := database.DB.Create(&member).Error; err != nil {
		return group, err
	}

	