package schema

type CreateExpense struct {
	CreatedBy         uint
	GroupID           uint                `json:"group_id" validate:"required"`
	Amount            float64             `json:"amount" validate:"required"`
	Description       string              `json:"description"`
	ExpenseShareType  string              `json:"expense_share_type" valdiate:"required,oneof=equal percent custom"`
	UserShare         float64             `json:"user_share"`
	MemberIDWithShare []MemberIDWithShare `json:"member_id_with_share" valdiate:"omitempty,min=1,dive"`
}

type MemberIDWithShare struct {
	ID    uint    `json:"id" valdiate:"required"`
	Share float64 `json:"share"`
}
