package schema

type PhoneNoLogin struct {
	CountryCode string `json:"country_code"`
	PhoneNumber string `json:"phone_number" validate:"len=10"`
}
