package model

type Error struct {
	Err     error
	Message string
	Code    int
}

const (
	InvalidPhoneNo  = "phone number doesn't exists"
	InvalidUsername = "no user found with this username"
	InvalidEmail    = "no user found with this email"
)

const (
	RequestAlreadyExists = "request already exists"
	RequestDoesntExist   = "request doesn't exists"
	RequestProcessUnable = "unable to process request"
)

const (
	SelfFriendRequest = "you can't send request to yourself"
)

func NewError(err error, message string, code int) *Error {
	return &Error{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
