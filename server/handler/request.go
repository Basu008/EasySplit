package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Basu008/EasySplit.git/server/auth"
)

type Request struct {
	HandlerFunc func(*RequestContext, http.ResponseWriter, *http.Request)
	AuthFunc    auth.TokenAuth
	IsLoggedIn  bool
	IsSudoUser  bool
}

func (rh *Request) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestCTX := &RequestContext{}
	requestCTX.Path = r.URL.Path
	authToken := r.Header.Get("Authorization")
	if authToken != "" {
		claim, err := rh.AuthFunc.VerifyToken(authToken)
		if err != nil {
			requestCTX.SetErr(fmt.Errorf("%s: failed to verify token", Unauthorized), http.StatusUnauthorized)
			goto SKIP_REQUEST
		} else {
			requestCTX.UserClaim = *claim
		}
	}

SKIP_REQUEST:
	if requestCTX.Err == nil {
		rh.HandlerFunc(requestCTX, w, r)
	}

	if requestCTX.ResponseCode != 0 && requestCTX.ResponseType != RedirectResp {
		w.WriteHeader(requestCTX.ResponseCode)
	}
	switch t := requestCTX.ResponseType; t {
	case JSONResp:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestCTX.Response)
	case ErrorResp:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&requestCTX.Err)
	}
}
