package model

type Error struct {
	Err     error
	Message string
	Code    int
}

const (
	InvalidPhoneNo = "phone number doesn't exists"
)

func NewError(err error, message string, code int) *Error {
	return &Error{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
