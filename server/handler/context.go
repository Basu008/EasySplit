package handler

import (
	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/server/auth"
)

type RequestContext struct {
	RequestID    string
	Path         string
	Response     Response
	Err          *AppErr
	ErrMsg       string
	ResponseType ResponseType
	ResponseCode int
	UserClaim    auth.UserClaim
}

func (requestCTX *RequestContext) SetErr(err error, errMsg string, statusCode int) {
	appErr := requestCTX.Err
	requestCTX.ResponseType = ErrorResp
	requestCTX.ResponseCode = statusCode
	if appErr == nil {
		appErr = &AppErr{}
	}
	appErr.Error = append(appErr.Error, *model.NewError(err, errMsg, statusCode))
	requestCTX.Err = appErr
}

func (requestCTX *RequestContext) SetAppResponse(message interface{}, statusCode int) {
	requestCTX.ResponseType = JSONResp
	requestCTX.ResponseCode = statusCode
	requestCTX.Response = &AppResponse{
		Payload: message,
	}
}

func (requestCTX *RequestContext) SetErrs(errs []error, statusCode int) {
	for _, e := range errs {
		requestCTX.SetErr(e, e.Error(), statusCode)
	}
}
