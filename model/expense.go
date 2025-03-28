package model

import "time"

const (
	Equal   = "equal"
	Percent = "percent"
	Custom  = "custom"
)

type Expense struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	GroupID     uint      `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	CreatedBy   uint      `gorm:"index;constraint:OnDelete:SET NULL;"`
	Amount      float64   `gorm:"type:decimal(10,2);not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	// Relationships
	Group Group `gorm:"foreignKey:GroupID"`
	User  *User `gorm:"foreignKey:CreatedBy"`
}

type ExpenseShare struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	ExpenseID uint    `gorm:"not null;index:idx_expense_user;constraint:OnDelete:CASCADE;"`
	UserID    uint    `gorm:"not null;index:idx_expense_user;constraint:OnDelete:CASCADE;"`
	Amount    float64 `gorm:"type:decimal(10,2);not null"`

	// Relationships
	Expense Expense `gorm:"foreignKey:ExpenseID"`
	User    User    `gorm:"foreignKey:UserID"`
}
