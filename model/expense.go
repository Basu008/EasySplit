package model

import "time"

const (
	Equal   = "equal"
	Percent = "percent"
	Custom  = "custom"
)

type Expense struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID     uint      `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	CreatedBy   uint      `json:"created_by" gorm:"index;constraint:OnDelete:SET NULL;"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Group Group `json:"-" gorm:"foreignKey:GroupID"`
	User  *User `json:"-" gorm:"foreignKey:CreatedBy"`
}

type ExpenseShare struct {
	ID        uint    `json:"-" gorm:"primaryKey;autoIncrement"`
	ExpenseID uint    `json:"-" gorm:"not null;index:idx_expense_user;constraint:OnDelete:CASCADE;"`
	UserID    uint    `json:"user_id" gorm:"not null;index:idx_expense_user;constraint:OnDelete:CASCADE;"`
	Amount    float64 `json:"amount" gorm:"type:decimal(10,2);not null"`
	IsSettled bool    `gorm:"default:false;not null"`

	// Relationships
	Expense Expense `json:"-" gorm:"foreignKey:ExpenseID"`
	User    User    `json:"-" gorm:"foreignKey:UserID"`
}

type ExpenseWithShares struct {
	ID           uint           `json:"id"`
	TotalAmount  float64        `json:"total_amount"`
	Description  string         `json:"descript"`
	CreatedBy    uint           `json:"created_by"`
	MembersShare []MembersShare `json:"members_share"`
}

type MembersShare struct {
	ID     uint    `json:"id"`
	Amount float64 `json:"amount"`
}
