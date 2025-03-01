package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) healthCheck(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	requestCTX.SetAppResponse(true, http.StatusOK)
}
