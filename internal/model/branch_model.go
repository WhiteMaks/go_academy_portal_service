package model

import "time"

type PostBranchV1Request struct {
	Name    string  `json:"name" binding:"required,min=1,branch_name_max"`
	Address string  `json:"address" binding:"required,min=1,branch_address_max"`
	Note    *string `json:"note,omitempty" binding:"omitempty,branch_note_max"`
}

type PutBranchV1Request struct {
	Name           string  `json:"name" binding:"required,min=1,branch_name_max"`
	Address        string  `json:"address" binding:"required,min=1,branch_address_max"`
	Note           *string `json:"note,omitempty" binding:"omitempty,branch_note_max"`
	ActivityStatus string  `json:"activity_status" binding:"required,oneof=active inactive"`
}

func (r PutBranchV1Request) AsCreateShape() PostBranchV1Request {
	return PostBranchV1Request{
		Name:    r.Name,
		Address: r.Address,
		Note:    r.Note,
	}
}

type BranchV1Response struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	Address        string     `json:"address"`
	Note           *string    `json:"note,omitempty"`
	ActivityStatus string     `json:"activity_status"`
	ArchivedAt     *time.Time `json:"archived_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type GetBranchesV1Query struct {
	Q      string `form:"q" binding:"omitempty,branches_search_q_max"`
	Limit  int    `form:"limit" binding:"omitempty,branches_page_limit"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
}

type GetBranchesV1Response struct {
	Branches []BranchV1Response `json:"branches"`
	Total    int64              `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}
