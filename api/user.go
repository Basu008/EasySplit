package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) signupUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SignupOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	user, err := a.App.User.SignupUser(&s)
	if err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
		return
	}
	requestCTX.SetAppResponse(user, http.StatusCreated)
}

func (a *API) loginUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.LoginOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	claim, err := a.App.User.LoginUser(&s)
	if err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
		return
	}
	a.TokenAuth.UserClaim = claim
	token, tokenErr := a.TokenAuth.SignToken()
	if tokenErr != nil {
		requestCTX.SetErr(tokenErr, "", http.StatusInternalServerError)
	}
	requestCTX.SetAppResponse(token, http.StatusOK)
}

func (a *API) getUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	userService := a.App.User
	id := a.getIDfromQuery(r, "id")
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
