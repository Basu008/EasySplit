package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) getUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	userService := a.App.User
	id := a.getIDfromPath(r, "id")
	if id != 0 {
		user, err := userService.GetUserByID(uint(id))
		if err != nil {
			requestCTX.SetErr(err.Err, err.Message, err.Code)
			return
		}
		requestCTX.SetAppResponse(user, http.StatusOK)
		return
	}
	username := r.URL.Query().Get("username")
	if username != "" {
		user, err := userService.GetUserByUsername(username)
		if err != nil {
			requestCTX.SetErr(err.Err, err.Message, err.Code)
			return
		}
		requestCTX.SetAppResponse(user, http.StatusOK)
		return
	}
	phoneNumber := r.URL.Query().Get("phone_number")
	if phoneNumber != "" {
		user, err := userService.GetUserByPhoneNo(phoneNumber)
		if err != nil {
			requestCTX.SetErr(err.Err, err.Message, err.Code)
			return
		}
		requestCTX.SetAppResponse(user, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(nil, http.StatusOK)
}

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
	var s schema.ConfirmOTPOpts
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

func (a *API) updateUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateUserOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if s.Username != "" && !a.IsUsernameValid(s.Username) {
		requestCTX.SetErr(nil, "weak username", http.StatusInternalServerError)
		return
	}
	s.ID = requestCTX.UserClaim.ID
	if custErr := a.App.User.UpdateUser(&s); custErr != nil {
		requestCTX.SetErr(custErr.Err, custErr.Message, custErr.Code)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusAccepted)
}
