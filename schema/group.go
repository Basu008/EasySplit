package schema

type CreateGroupOpts struct {
	OwnerID   uint
	Name      string `json:"name" validate:"required"`
	Type      string `json:"type" validate:"required"`
	MemberIDs []uint `json:"member_ids" validate:"required,min=1"`
}

type EditGroupInfoOpts struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UpdateGroupMembers struct {
	ID        uint
	MemberID  uint   `json:"member_id"`
	Operation string `json:"operation" validate:"oneof=add remove"`
}

type GroupResponse struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	Type          string           `json:"type"`
	TotalExpense  float64          `json:"total_expense"`
	SettledAmount float64          `json:"settled_amount"`
	Members       []MemberResponse `json:"members"`
}

type MemberResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
