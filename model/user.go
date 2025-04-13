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
	ID            uint      `json:"id,omitempty" gorm:"primary key;autoIncrement"`
	Username      *string   `json:"username,omitempty" gorm:"unique"`
	PhoneNumber   string    `json:"phone_number,omitempty" gorm:"unique;not null"`
	CountryCode   string    `json:"country_code"`
	Email         *string   `json:"email,omitempty" gorm:"unique"`
	Plan          string    `json:"plan,omitempty" gorm:"default:FREE"`
	PhoneVerified bool      `json:"phone_verified,omitempty"  gorm:"default:false"`
	EmailVerified bool      `json:"email_verified,omitempty"  gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"-" gorm:"autoUpdateime"`
}
