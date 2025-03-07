package schema

type FriendRequestOpts struct {
	SenderUserID   uint
	ReceiverUserID uint   `json:"receiver_user_id"`
	RequestStatus  string `json:"request_status" validate:"omitempty,oneof=accepted rejected"`
}

type GetAllFriendsResponse struct {
	ID          uint   `json:"id"`
	Username    string `json:"username,omitempty"`
	PhoneNumber string `json:"phone_number"`
}
