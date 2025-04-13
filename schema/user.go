package schema

type GetUserOpts struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
}

type SignupOpts struct {
	FullName    string `json:"full_name" validate:"required,alpha_space"`
	Username    string `json:"username" validate:"required,username"`
	Password    string `json:"password" validate:"required,password"`
	PhoneNumber string `json:"phone_number" validate:"required,len=10,numeric"`
	CountryCode string `json:"country_code" validate:"required,eq=+91"`
	Email       string `json:"email" validate:"required,email"`
}

type LoginOpts struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserOpts struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
