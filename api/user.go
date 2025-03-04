package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) loginUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.PhoneNoLogin
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if err := a.App.User.LoginUser(&s); err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) confirmOTP(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ConfirmOTPRequest
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if userClaim, err := a.App.User.ConfirmOTP(&s); err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
		return
	} else {
		a.TokenAuth.UserClaim = userClaim
		token, err := a.TokenAuth.SignToken()
		if err != nil {
			requestCTX.SetErr(nil, "unable to sign in", http.StatusInternalServerError)
		}
		requestCTX.SetAppResponse(token, http.StatusOK)
	}
}
