package models

import "time"

type Group struct {
	ID                 string `gorm:"primaryKey"`
	Name               string
	Description        string
	Wallet             string
	SecretKey          string `gorm:"column:secret_key"`
	CreatorID          string `gorm:"column:creator_id"`
	Creator            User   `gorm:"foreignKey:CreatorID"`
	Members            []Member       `gorm:"foreignKey:GroupID"`
	Contributions      []Contribution `gorm:"foreignKey:GroupID"`
	ContractID         string         `gorm:"column:contract_id"`
	Status             string         `gorm:"default:pending"` // pending, active, completed
	ContributionAmount float64        `gorm:"column:contribution_amount"`
	ContributionPeriod int            `gorm:"column:contribution_period"` // days
	PayoutOrder        string         `gorm:"column:payout_order"` // JSON array of member IDs
	CurrentRound       int            `gorm:"column:current_round;default:0"`
	MaxMembers         int            `gorm:"column:max_members;default:20"`
	MinMembers         int            `gorm:"column:min_members;default:3"`
	NextContributionDate time.Time `gorm:"column:next_contribution_date"`
	IsApproved         bool          `gorm:"column:is_approved;default:false"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Member struct {
	ID       string `gorm:"primaryKey"`
	GroupID  string
	UserID   string
	User     User   `gorm:"foreignKey:UserID"`
	Wallet   string
	Role     string `gorm:"default:member"` // member, admin, creator
	JoinedAt time.Time
	Status   string `gorm:"default:pending"` // pending, approved, rejected
}

type GroupInvitation struct {
	ID        string `gorm:"primaryKey"`
	GroupID   string
	Group     Group  `gorm:"foreignKey:GroupID"`
	InviterID string
	Inviter   User   `gorm:"foreignKey:InviterID"`
	Email     string
	UserID    string // if user exists
	User      User   `gorm:"foreignKey:UserID"`
	Status    string `gorm:"default:pending"` // pending, accepted, rejected
	CreatedAt time.Time
	ExpiresAt time.Time
}
