package api

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/handler"
)

func (a *API) createPayment(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreatePaymentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, "", http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	s.PayerID = a.TokenAuth.UserClaim.ID
	expenseShare, err := a.App.Expense.GetExpenseShare(s.ExpenseID, s.PayerID)
	if err != nil {
		requestCTX.SetErr(nil, "expense doesn't exists", http.StatusBadRequest)
		return
	}
	if s.Amount < expenseShare.Amount {
		requestCTX.SetErr(nil, "partial payments are not allowed as of now", http.StatusBadRequest)
		return
	}
	s.PayeeID = expenseShare.Expense.CreatedBy
	if errResp := a.App.Payment.CreatePayment(&s); errResp != nil {
		requestCTX.SetErr(errResp.Err, errResp.Message, errResp.Code)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
}
