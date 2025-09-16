package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/utils"
)

var jwtSecret = []byte("your-secret-key-change-in-production") // Change this in production

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

