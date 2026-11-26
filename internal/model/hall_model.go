package model

import "time"

type PostHallV1Request struct {
	BranchID int64   `json:"branch_id" binding:"required,gt=0"`
	Name     string  `json:"name" binding:"required,min=1,hall_name_max"`
	Subtitle *string `json:"subtitle,omitempty" binding:"omitempty,min=1,hall_subtitle_max"`
	Capacity int     `json:"capacity" binding:"required,hall_capacity_range"`
	Note     *string `json:"note,omitempty" binding:"omitempty,hall_note_max"`
}

type PutHallV1Request struct {
	Name           string  `json:"name" binding:"required,min=1,hall_name_max"`
	Subtitle       *string `json:"subtitle,omitempty" binding:"omitempty,min=1,hall_subtitle_max"`
	Capacity       int     `json:"capacity" binding:"required,hall_capacity_range"`
	Note           *string `json:"note,omitempty" binding:"omitempty,hall_note_max"`
	ActivityStatus string  `json:"activity_status" binding:"required,oneof=active inactive"`
}

func (r PutHallV1Request) AsCreateShape() PostHallV1Request {
	return PostHallV1Request{
		Name:     r.Name,
		Subtitle: r.Subtitle,
		Capacity: r.Capacity,
		Note:     r.Note,
	}
}

type HallV1Response struct {
	ID                   int64      `json:"id"`
	BranchID             int64      `json:"branch_id"`
	Name                 string     `json:"name"`
	Subtitle             *string    `json:"subtitle,omitempty"`
	Capacity             int        `json:"capacity"`
	Note                 *string    `json:"note,omitempty"`
	ActivityStatus       string     `json:"activity_status"`
	BranchName           *string    `json:"branch_name,omitempty"`
	BranchAddress        *string    `json:"branch_address,omitempty"`
	BranchActivityStatus *string    `json:"branch_activity_status,omitempty"`
	ArchivedAt           *time.Time `json:"archived_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type GetBranchHallsV1Query struct {
	Q      string `form:"q" binding:"omitempty,halls_search_q_max"`
	Limit  int    `form:"limit" binding:"omitempty,halls_page_limit"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
}

type GetBranchHallsV1Response struct {
	Halls  []HallV1Response `json:"halls"`
	Total  int64            `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}
