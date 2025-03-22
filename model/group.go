package model

import "time"

type Group struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"not null"`
	OwnerID       uint    `gorm:"not null;index"`
	Type          string  `gorm:"not null"`
	TotalExpense  float64 `gorm:"check:total_expense>=0;default:0"`
	SettledAmount float64 `gorm:"check:settled_amount >= 0;default:0;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Owner   User          `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Members []GroupMember `gorm:"foreignKey:GroupID"`
}

type GroupMember struct {
	GroupID  uint      `gorm:"not null;index"`
	MemberID uint      `gorm:"not null;index"`
	AddedAt  time.Time `gorm:"autoCreateTime"`

	Group  Group `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Member User  `gorm:"foreignKey:MemberID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
