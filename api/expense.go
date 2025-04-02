package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) createExpense(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateExpense
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	s.CreatedBy = requestCTX.UserClaim.ID
	group, err := a.App.Group.GetGroupByID(s.GroupID)
	if err != nil {
		requestCTX.SetErr(err.Err, err.Message, err.Code)
		return
	}
	if s.ExpenseShareType != model.Equal {
		if len(s.MemberIDWithShare) < 1 {
			requestCTX.SetErr(nil, "member_id_with_share is a required field for expense_share_type = percent & custom", http.StatusBadRequest)
			return
		}
		if !validateMembers(group, &s) {
			requestCTX.SetErr(nil, "invalid members", http.StatusBadRequest)
			return
		}
		if !validateShare(&s) {
			requestCTX.SetErr(nil, "incorrect share distribution", http.StatusBadRequest)
			return
		}
	}
	if len(s.MemberIDWithShare) == 0 {
		for _, member := range group.Members {
			if member.ID == s.CreatedBy {
				continue
			}
			s.MemberIDWithShare = append(s.MemberIDWithShare, schema.MemberIDWithShare{
				ID: member.ID,
			})
		}
	}
	if errResp := a.App.Expense.CreateExpense(&s); errResp != nil {
		requestCTX.SetErr(errResp.Err, errResp.Message, errResp.Code)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
}

func (a *API) getExpense(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	expenseID := a.getIDfromPath(r, "expense_id")
	expense, errResp := a.App.Expense.GetExpense(expenseID)
	if errResp != nil {
		requestCTX.SetErr(errResp.Err, errResp.Message, errResp.Code)
		return
	}
	requestCTX.SetAppResponse(expense, http.StatusOK)
}

func (a *API) getExpenses(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	groupID := a.getIDfromPath(r, "group_id")
	expenses, errResp := a.App.Expense.GetExpenses(groupID)
	if errResp != nil {
		requestCTX.SetErr(errResp.Err, errResp.Message, errResp.Code)
		return
	}
	requestCTX.SetAppResponse(expenses, http.StatusOK)
}

func (a *API) deleteExpense(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	expenseID := a.getIDfromPath(r, "expense_id")
	requestCTX.SetAppResponse(a.App.Expense.DeleteExpense(expenseID), http.StatusOK)
}

func validateMembers(group *schema.GroupResponse, s *schema.CreateExpense) bool {
	for _, memeberOpt := range s.MemberIDWithShare {
		var found bool
		for _, member := range group.Members {
			if member.ID == memeberOpt.ID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func validateShare(s *schema.CreateExpense) bool {
	amount := s.Amount
	var share float64
	for _, opt := range s.MemberIDWithShare {
		share += opt.Share
	}
	if s.UserShare != 0 {
		share += s.UserShare
	}
	if s.ExpenseShareType == model.Percent {
		return share == 100
	}
	return share == amount
}
