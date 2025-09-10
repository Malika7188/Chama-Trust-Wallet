package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"chama-wallet-backend/database"
	"chama-wallet-backend/models"
	"chama-wallet-backend/services"
	"chama-wallet-backend/config"
)

func ContributeToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		From   string `json:"from"`
		Secret string `json:"secret"`
		Amount string `json:"amount"`
	}
	
	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("❌ Failed to parse request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if payload.From == "" || payload.Secret == "" || payload.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: from, secret, and amount are required",
		})
	}

	// Validate amount limits for mainnet
	if config.Config.IsMainnet {
		amount, err := strconv.ParseFloat(payload.Amount, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid amount format",
			})
		}
		
		minAmount := 0.0000001 // Minimum XLM amount
		if amount < minAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount below minimum of %f XLM", minAmount),
			})
		}
	}
	// Verify user is a member of the group
	var member models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		groupID, user.ID, "approved").First(&member).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not an approved member of this group",
		})
	}

	// Get group details
	group, err := services.GetGroupByID(groupID)
	if err != nil {
		fmt.Printf("❌ Group not found: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Group not found"})
	}

	// Validate group status
	if group.Status != "active" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Group is not active (current status: %s)", group.Status),
		})
	}

	// Validate contract ID
	if group.ContractID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Group contract ID is not set",
		})
	}

	// Validate user's wallet address matches the 'from' field
	if payload.From != user.Wallet {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "From address must match your wallet address",
		})
	}

	fmt.Printf("🔄 Processing contribution: %s XLM from %s to group %s (contract: %s) on %s\n", 
		payload.Amount, payload.From, group.Name, group.ContractID, config.Config.Network)

	// Make authenticated Soroban contract call
	output, err := services.ContributeWithAuth(group.ContractID, payload.From, payload.Amount, payload.Secret)
	if err != nil {
		fmt.Printf("❌ Soroban contribution failed: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Blockchain transaction failed: %v", err),
		})
	}

	// Record the contribution in the database
	contribution := models.Contribution{
		ID:        uuid.NewString(),
		GroupID:   groupID,
		UserID:    user.ID,
		Amount:    parseFloat64(payload.Amount),
		Status:    "confirmed",
		TxHash:    output,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&contribution).Error; err != nil {
		fmt.Printf("⚠️ Warning: Failed to record contribution in database: %v\n", err)
		// Don't fail the request since blockchain transaction succeeded
	}

	fmt.Printf("✅ Contribution successful on %s: %s\n", config.Config.Network, output)

	return c.JSON(fiber.Map{
		"message":      "Contribution successful",
		"group_id":     groupID,
		"group_name":   group.Name,
		"from":         payload.From,
		"to":           group.Wallet,
		"amount":       payload.Amount,
		"tx_hash":      output,
		"network":      config.Config.Network,
		"contribution": contribution,
	})
}

func GetUserGroups(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	groups, err := services.GetUserGroups(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(groups)
}

func GetNonGroupMembers(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is admin/creator of the group
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?", 
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	// Get users who are not members of this group
	var users []models.User
	err := database.DB.
		Where("id NOT IN (SELECT user_id FROM members WHERE group_id = ?)", groupID).
		Find(&users).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}

func InviteToGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	var payload struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Check if user is admin/creator
	var admin models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ? AND role IN ?", 
		groupID, user.ID, []string{"creator", "admin"}).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can invite users"})
	}

	// Check if user exists
	var invitedUser models.User
	if err := database.DB.Where("email = ?", payload.Email).First(&invitedUser).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User with this email not found"})
	}

	// Check if user is already a member
	var existingMember models.Member
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, invitedUser.ID).First(&existingMember).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already a member of this group"})
	}

	// Create invitation
	invitation := models.GroupInvitation{
		ID:        uuid.NewString(),
		GroupID:   groupID,
		InviterID: user.ID,
		Email:     payload.Email,
		UserID:    invitedUser.ID,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := database.DB.Create(&invitation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Create notification for invited user
	var group models.Group
	database.DB.First(&group, "id = ?", groupID)

	notification := models.Notification{
		ID:        uuid.NewString(),
		UserID:    invitedUser.ID,
		GroupID:   groupID,
		Type:      "group_invitation",
		Title:     "Group Invitation",
		Message:   fmt.Sprintf("You have been invited to join %s by %s", group.Name, user.Name),
		Status:    "unread",
		Read:      false,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		fmt.Printf("⚠️ Failed to create notification: %v\n", err)
	} else {
		fmt.Printf("✅ Notification created for user %s\n", invitedUser.ID)
	}

	return c.JSON(fiber.Map{"message": "Invitation sent successfully"})
}

func ApproveGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	user := c.Locals("user").(models.User)

	// Check if user is the creator
	var group models.Group
	if err := database.DB.Where("id = ? AND creator_id = ?", groupID, user.ID).First(&group).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only group creator can approve the group"})
	}

	// Check if group has minimum members
	var memberCount int64
	database.DB.Model(&models.Member{}).Where("group_id = ? AND status = ?", groupID, "approved").Count(&memberCount)
	if memberCount < int64(group.MinMembers) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Group needs at least %d members before approval (currently has %d)", group.MinMembers, memberCount),
		})
	}

	