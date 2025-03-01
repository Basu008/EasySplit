package handler

import "encoding/json"

type ResponseType string

const (
	HTMLResp     ResponseType = "html"
	JSONResp     ResponseType = "json"
	RedirectResp ResponseType = "redirect"
	FileResp     ResponseType = "file"
	ErrorResp    ResponseType = "error"
)

type AppErr struct {
	Error     []error
	RequestID *string
}

type Response interface {
	MarshalJSON() ([]byte, error)
	GetRaw() interface{}
}

type AppResponse struct {
	Payload interface{}
}

func (r *AppResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Success bool        `json:"success"`
		Payload interface{} `json:"payload"`
	}{
		Success: true,
		Payload: &r.Payload,
	})
}

func (r *AppResponse) GetRaw() interface{} {
	return r.Payload
}

var (
	NoType             = "NoType"
	BadRequest         = "BadRequest"
	NotFound           = "NotFound"
	DBError            = "DBError"
	Unauthorized       = "Unauthorized"
	PermissionDenied   = "PermissionDenied"
	SomethingWentWrong = "SomethingWentWrong"
)
