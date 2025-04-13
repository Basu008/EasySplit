package schema

type CreateGroupOpts struct {
	OwnerID uint
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type" validate:"required"`
	UserIDs []uint `json:"user_ids" validate:"required,min=1"`
}

type EditGroupInfoOpts struct {
	ID   uint
	Name string `json:"name"`
	Type string `json:"type"`
}

type RemoveGroupMemberOpts struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
}

type AddGroupMembersOpts struct {
	ID      uint   `json:"id" validate:"required"`
	UserIDs []uint `json:"user_ids" validate:"required,min=1"`
}

type GroupResponse struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	Type          string           `json:"type"`
	TotalExpense  float64          `json:"total_expense"`
	SettledAmount float64          `json:"settled_amount"`
	Members       []MemberResponse `json:"members,omitempty"`
}

type MemberResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}
