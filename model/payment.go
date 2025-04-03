package model

import "time"

type Payment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	PayerID   uint      `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	PayeeID   uint      `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	ExpenseID uint      `gorm:"not null;index;constraint:OnDelete:CASCADE;"`
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	Mode      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	Payer   User    `gorm:"foreignKey:PayerID"`
	Payee   User    `gorm:"foreignKey:PayeeID"`
	Expense Expense `gorm:"foreignKey:ExpenseID"`
}
