package handler

import (
	"encoding/json"

	"github.com/Basu008/EasySplit.git/model"
)

type ResponseType string

const (
	HTMLResp     ResponseType = "html"
	JSONResp     ResponseType = "json"
	RedirectResp ResponseType = "redirect"
	FileResp     ResponseType = "file"
	ErrorResp    ResponseType = "error"
)

type AppErr struct {
	Error []model.Error
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

func (err *AppErr) MarshalJSON() ([]byte, error) {
	errs := []string{}
	for _, e := range err.Error {
		if e.Message != "" {
			errs = append(errs, e.Message)
		} else {
			errs = append(errs, e.Err.Error())
		}
	}
	return json.Marshal(&struct {
		Error   []string `json:"errors"`
		Success bool     `json:"success"`
	}{
		Error:   errs,
		Success: false,
	})
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
