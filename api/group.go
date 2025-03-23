package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) createGroup(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateGroupOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	s.OwnerID = requestCTX.UserClaim.ID
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if errResp := a.App.Group.CreateGroup(&s); errResp != nil {
		requestCTX.SetErr(errResp.Err, errResp.Message, errResp.Code)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
}

func (a *API) getGroupByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	groupID := a.getIDfromPath(r, "ID")
	group, errResp := a.App.Group.GetGroupByID(groupID)
	if errResp != nil {
		requestCTX.SetErr(errResp.Err, errResp.Message, errResp.Code)
		return
	}
	requestCTX.SetAppResponse(group, http.StatusOK)
}

func (a *API) getGroups(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	page := a.getPageValue(r)
	ownerID := requestCTX.UserClaim.ID
	groups, errResp := a.App.Group.GetGroups(ownerID, page)
	if errResp != nil {
		requestCTX.SetErr(errResp.Err, errResp.Message, errResp.Code)
		return
	}
	requestCTX.SetAppResponse(groups, http.StatusOK)
}

func (a *API) editGroup(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditGroupInfoOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	groupID := a.getIDfromPath(r, "ID")
	s.ID = groupID
	errResp := a.App.Group.EditGroup(&s)
	if errResp != nil {
		requestCTX.SetErr(errResp.Err, "", errResp.Code)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
