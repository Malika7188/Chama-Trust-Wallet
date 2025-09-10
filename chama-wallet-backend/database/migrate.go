package database

import (
	"log"

	"chama-wallet-backend/models"
)

func RunMigrations() {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.Member{},
		&models.Contribution{},
		&models.GroupInvitation{},
		&models.AdminNomination{},
		&models.PayoutRequest{},
		&models.PayoutApproval{},
		&models.PayoutSchedule{},
		&models.Notification{},
		&models.RoundContribution{},
		&models.RoundStatus{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Database migrations completed successfully")
}
