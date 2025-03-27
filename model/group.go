package model

import "time"

const (
	GroupName = "name"
	GroupType = "type"
)

type Group struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Type      string    `gorm:"type:varchar(255);not null"`
	CreatedBy uint      `gorm:"index;constraint:OnDelete:SET NULL;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	User      *User     `gorm:"foreignKey:CreatedBy"`
}

type GroupMember struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	GroupID  uint      `gorm:"not null;index:idx_group_user;constraint:OnDelete:CASCADE;"`
	UserID   uint      `gorm:"not null;index:idx_group_user;constraint:OnDelete:CASCADE;"`
	JoinedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	Group Group `gorm:"foreignKey:GroupID"`
	User  User  `gorm:"foreignKey:UserID"`
}
