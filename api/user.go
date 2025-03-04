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
	if err := a.App.User.LoginUser(&s); err != nil {
		requestCTX.SetErr(err.Err, err.Message, http.StatusInternalServerError)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
