package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Basu008/EasySplit.git/server/auth"
)

type Request struct {
	HandlerFunc func(*RequestContext, http.ResponseWriter, *http.Request)
	AuthFunc    auth.TokenAuthentication
	IsLoggedIn  bool
	IsSudoUser  bool
}

func (rh *Request) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestCTX := &RequestContext{}
	requestCTX.Path = r.URL.Path
	auth := r.Header.Get("Authorization")
	if auth != "" {
		d := strings.Split(auth, "Bearer ")
		if len(d) < 2 {
			requestCTX.SetErr(nil, fmt.Sprintf("%s: authorization is of type bearer", BadRequest), http.StatusUnauthorized)
		} else {
			claim, err := rh.AuthFunc.VerifyToken(d[1])
			if err != nil {
				requestCTX.SetErr(fmt.Errorf("%s: failed to verify token", Unauthorized), "", http.StatusUnauthorized)
			} else {
				requestCTX.UserClaim = *claim
			}
		}
	}
	if requestCTX.Err == nil {
		rh.HandlerFunc(requestCTX, w, r)
	}
	switch t := requestCTX.ResponseType; t {
	case JSONResp:
		w.Header().Set("Content-Type", "application/json")
		if requestCTX.ResponseCode != 0 {
			w.WriteHeader(requestCTX.ResponseCode)
		}
		json.NewEncoder(w).Encode(requestCTX.Response)
	case ErrorResp:
		w.Header().Set("Content-Type", "application/json")
		if requestCTX.ResponseCode != 0 {
			w.WriteHeader(requestCTX.ResponseCode)
		}
		if requestCTX.ErrMsg != "" {
			json.NewEncoder(w).Encode(&requestCTX.ErrMsg)
		} else {
			json.NewEncoder(w).Encode(&requestCTX.Err)
		}
	}

}
