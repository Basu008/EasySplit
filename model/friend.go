package model

import "time"

// keys
const (
	SenderUserID   = "sender_user_id"
	ReceiverUserID = "receiver_user_id"
	RequestStatus  = "request_status"
)

// friend request status
const (
	Requested = "requested"
	Rejected  = "rejected"
	Accepted  = "accepted"
)

type Friend struct {
	SenderUserID   uint      `json:"sender_user_id" gorm:"not null;index:idx_sender_receiver"`
	ReceiverUserID uint      `json:"receiver_user_id" gorm:"not null;index:idx_sender_receiver"`
	RequestStatus  string    `json:"request_status"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateime"`

	// foreign keys
	SenderUser   User `json:"-" gorm:"foreignKey:SenderUserID;constraint:OnDelete:CASCADE;"`
	ReceiverUser User `json:"-" gorm:"foreignKey:ReceiverUserID;constraint:OnDelete:CASCADE;"`
}
