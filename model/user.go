package model

import "time"

// fields
const (
	PhoneNumber = "phone_number"
	OTP         = "otp"
)

type User struct {
	ID            uint      `json:"id" gorm:"primary key;autoIncrement"`
	Username      *string   `json:"username" gorm:"unique"`
	PhoneNumber   string    `json:"phone_number" gorm:"unique;not null"`
	CountryCode   string    `json:"country_code"`
	OTP           string    `json:"otp"`
	Email         string    `json:"email"`
	Plan          string    `json:"plan" gorm:"default:FREE"`
	PhoneVerified bool      `json:"phone_verified"  gorm:"default:false"`
	EmailVerified bool      `json:"email_verified"  gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateime"`
}
