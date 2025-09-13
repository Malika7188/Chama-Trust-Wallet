package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
)

func CreatePayoutRequest(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		RecipientID string  `json:"recipient_id"`
		Amount      float64 `json:"amount"`
		Round       int     `json:"round"`
	}

	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("❌ Failed to parse payout request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Validate required fields
	if payload.RecipientID == "" || payload.Amount <= 0 || payload.Round <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid required fields",
		})
	}

	// Get group details
	var group models.Group
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	// Validate group is active
	if group.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group must be active to create payout requests",
		})
	}

	// Check if user is admin/creator
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?",
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can create payout requests"})
	}

	// Verify recipient is group member
	var recipient models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, payload.RecipientID, "approved").First(&recipient).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Recipient is not a group member"})
	}

	// Check if payout request already exists for this round
	var existingRequest models.PayoutRequest
	if err := database.DB.Where("group_id = ? AND round = ? AND status IN ?",
		groupID, payload.Round, []string{"pending", "approved"}).First(&existingRequest).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Payout request already exists for this round",
		})
	}

	// Validate payout amount against group balance
	groupBalance, err := services.CheckBalance(group.Wallet)
	if err != nil {
		fmt.Printf("⚠️ Warning: Could not check group balance: %v\n", err)
	} else {
		if balance, parseErr := strconv.ParseFloat(groupBalance, 64); parseErr == nil {
			if payload.Amount > balance {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Insufficient group balance. Available: %.2f XLM, Requested: %.2f XLM", balance, payload.Amount),
				})
			}
		}
	}

	payoutRequest := models.PayoutRequest{
		ID:          uuid.NewString(),
		GroupID:     groupID,
		RecipientID: payload.RecipientID,
		Amount:      payload.Amount,
		Round:       payload.Round,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&payoutRequest).Error; err != nil {
		fmt.Printf("❌ Failed to create payout request: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("✅ Payout request created: %s\n", payoutRequest.ID)

	// Notify all admins about the payout request (excluding the creator)
	var admins []models.Member
	database.DB.Where("group_id = ? AND role IN ? AND status = ?",
		groupID, []string{"creator", "admin"}, "approved").Find(&admins)

	for _, admin := range admins {
		if admin.UserID != user.ID {
			services.CreateNotification(
				admin.UserID,
				groupID,
				"payout_request",
				"Payout Request Created",
				fmt.Sprintf("New payout request for %.2f XLM to %s requires approval", payload.Amount, recipient.User.Name),
			)
		}
	}

	return c.JSON(fiber.Map{
		"message": "Payout request created successfully",
		"request": payoutRequest,
	})
}

func ApprovePayoutRequest(c *fiber.Ctx) error {
	payoutID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Approved bool `json:"approved"`
	}

	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("❌ Failed to parse approval request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	