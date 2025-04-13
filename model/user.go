package model

import "time"

// fields
const (
	Username      = "username"
	Email         = "email"
	PhoneNumber   = "phone_number"
	OTP           = "otp"
	PhoneVerified = "phone_verified"
)

type User struct {
	ID          uint      `json:"id,omitempty" gorm:"primary key;autoIncrement"`
	FullName    string    `json:"name,omitempty" gorm:"not null"`
	Username    string    `json:"username,omitempty" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"`
	CountryCode string    `json:"country_code" gorm:"not null"`
	PhoneNumber string    `json:"phone_number,omitempty" gorm:"unique;not null"`
	Email       string    `json:"email,omitempty" gorm:"unique;not null"`
	Plan        string    `json:"plan,omitempty" gorm:"default:FREE"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"-" gorm:"autoUpdateime"`
}
