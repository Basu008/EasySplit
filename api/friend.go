package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) sendFriendRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.FriendRequestOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	s.SenderUserID = requestCTX.UserClaim.ID
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.SenderUserID == s.ReceiverUserID {
		requestCTX.SetErr(nil, model.SelfFriendRequest, http.StatusBadRequest)
		return
	}
	s.RequestStatus = model.Requested
	err := a.App.Friend.SendFriendRequest(&s)
	if err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) updateFriendRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateFriendRequestOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	s.ReceiverUserID = requestCTX.UserClaim.ID
	err := a.App.Friend.UpdateFriendRequest(&s)
	if err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) allFriends(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	userID := requestCTX.UserClaim.ID
	page := a.getPageValue(r)
	friends := a.App.Friend.GetAllFriends(userID, page)
	requestCTX.SetAppResponse(friends, http.StatusOK)
}

func (a *API) getFriendStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	userID := requestCTX.UserClaim.ID
	friendID := a.getIDfromPath(r, "id")
	friend, err := a.App.Friend.GetFriendStatus(userID, friendID)
	if err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
		return
	}
	requestCTX.SetAppResponse(friend, http.StatusOK)
}
