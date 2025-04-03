package schema

type CreatePaymentOpts struct {
	PayerID   uint
	PayeeID   uint
	ExpenseID uint    `json:"expense_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=1"`
	Mode      string  `json:"mode" validate:"required,oneof=card upi cash net-banking"`
}
