package model

import "time"

type PostCoachV1Request struct {
	LastName       string   `json:"last_name" binding:"required,min=1,coach_name_max"`
	FirstName      string   `json:"first_name" binding:"required,min=1,coach_name_max"`
	Patronymic     *string  `json:"patronymic,omitempty" binding:"omitempty,min=1,coach_name_max"`
	BirthDate      DateOnly `json:"birth_date" binding:"required"`
	Phone          string   `json:"phone" binding:"required,e164"`
	Specialization string   `json:"specialization" binding:"required,min=1,coach_specialization_max"`
	Note           *string  `json:"note,omitempty" binding:"omitempty,coach_note_max"`
}

type PutCoachV1Request struct {
	LastName       string   `json:"last_name" binding:"required,min=1,coach_name_max"`
	FirstName      string   `json:"first_name" binding:"required,min=1,coach_name_max"`
	Patronymic     *string  `json:"patronymic,omitempty" binding:"omitempty,min=1,coach_name_max"`
	BirthDate      DateOnly `json:"birth_date" binding:"required"`
	Phone          string   `json:"phone" binding:"required,e164"`
	Specialization string   `json:"specialization" binding:"required,min=1,coach_specialization_max"`
	Note           *string  `json:"note,omitempty" binding:"omitempty,coach_note_max"`
	ActivityStatus string   `json:"activity_status" binding:"required,oneof=active inactive"`
}

// AsCreateShape returns the shared card fields for create/update validation helpers.
func (r PutCoachV1Request) AsCreateShape() PostCoachV1Request {
	return PostCoachV1Request{
		LastName:       r.LastName,
		FirstName:      r.FirstName,
		Patronymic:     r.Patronymic,
		BirthDate:      r.BirthDate,
		Phone:          r.Phone,
		Specialization: r.Specialization,
		Note:           r.Note,
	}
}

type CoachV1Response struct {
	ID             int64      `json:"id"`
	LastName       string     `json:"last_name"`
	FirstName      string     `json:"first_name"`
	Patronymic     *string    `json:"patronymic,omitempty"`
	BirthDate      DateOnly   `json:"birth_date"`
	Phone          string     `json:"phone"`
	Specialization string     `json:"specialization"`
	Note           *string    `json:"note,omitempty"`
	ActivityStatus string     `json:"activity_status"`
	ArchivedAt     *time.Time `json:"archived_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type GetCoachesV1Query struct {
	Q      string `form:"q" binding:"omitempty,coaches_search_q_max"`
	Limit  int    `form:"limit" binding:"omitempty,coaches_page_limit"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
}

type GetCoachesV1Response struct {
	Coaches []CoachV1Response `json:"coaches"`
	Total   int64             `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
}
