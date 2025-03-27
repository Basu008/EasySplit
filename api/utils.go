package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Basu008/EasySplit.git/server/handler"
	"github.com/gorilla/mux"
)

func (a *API) healthCheck(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) DecodeJSONBody(r *http.Request, res interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		if r.Header.Get("Content-Type") != "application/json" {
			err := errors.New("unsupported content-type request: Content-Type header is not application/json")
			return err
		}
	}
	if r.ContentLength == 0 {
		return errors.New("request body must not be empty")
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(&res)
}

func (a *API) IsUsernameValid(username string) bool {
	regex := `^[a-zA-Z0-9@._]{5,30}$`
	matched, err := regexp.MatchString(regex, username)
	if err != nil {
		return false
	}
	return matched
}

func (a *API) getPageValue(r *http.Request) int {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	return page
}

func (a *API) getIDfromPath(r *http.Request, key string) uint {
	idString := mux.Vars(r)[key]
	var id uint64
	if idString != "" {
		var err error
		id, err = strconv.ParseUint(idString, 10, 32)
		if err != nil {
			return 0
		}
		return uint(id)
	}
	return uint(id)
}
