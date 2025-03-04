package schema

type GetUserOpts struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
}

type PhoneNoLogin struct {
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number" validate:"len=10"`
}

type ConfirmOTPOpts struct {
	PhoneNumber string `json:"phone_number" validate:"len=10"`
	OTP         string `json:"otp" validate:"len=4"`
}

type UpdateUserOpts struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
