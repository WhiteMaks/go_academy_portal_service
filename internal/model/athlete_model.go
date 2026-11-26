package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DateOnly is a calendar date (YYYY-MM-DD) in JSON.
type DateOnly struct {
	time.Time
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Time = time.Time{}
		return nil
	}
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	raw = strings.TrimSpace(raw)
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return fmt.Errorf("invalid date %q: %w", raw, err)
	}
	d.Time = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
	return nil
}

type PostAthleteV1Request struct {
	LastName         string    `json:"last_name" binding:"required,min=1,athlete_name_max"`
	FirstName        string    `json:"first_name" binding:"required,min=1,athlete_name_max"`
	Patronymic       *string   `json:"patronymic,omitempty" binding:"omitempty,min=1,athlete_name_max"`
	BirthDate        DateOnly  `json:"birth_date" binding:"required"`
	Phone            string    `json:"phone" binding:"required,e164"`
	ParentLastName   *string   `json:"parent_last_name,omitempty" binding:"omitempty,min=1,athlete_name_max"`
	ParentFirstName  *string   `json:"parent_first_name,omitempty" binding:"omitempty,min=1,athlete_name_max"`
	ParentPatronymic *string   `json:"parent_patronymic,omitempty" binding:"omitempty,min=1,athlete_name_max"`
	ParentBirthDate  *DateOnly `json:"parent_birth_date,omitempty"`
	ParentPhone      *string   `json:"parent_phone,omitempty" binding:"omitempty,e164"`
	Note             *string   `json:"note,omitempty" binding:"omitempty,athlete_note_max"`
}

type AthleteV1Response struct {
	ID               int64      `json:"id"`
	LastName         string     `json:"last_name"`
	FirstName        string     `json:"first_name"`
	Patronymic       *string    `json:"patronymic,omitempty"`
	BirthDate        DateOnly   `json:"birth_date"`
	Phone            string     `json:"phone"`
	ParentLastName   *string    `json:"parent_last_name,omitempty"`
	ParentFirstName  *string    `json:"parent_first_name,omitempty"`
	ParentPatronymic *string    `json:"parent_patronymic,omitempty"`
	ParentBirthDate  *DateOnly  `json:"parent_birth_date,omitempty"`
	ParentPhone      *string    `json:"parent_phone,omitempty"`
	Note             *string    `json:"note,omitempty"`
	ActivityStatus   string     `json:"activity_status"`
	ArchivedAt       *time.Time `json:"archived_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type GetAthletesV1Query struct {
	Q      string `form:"q" binding:"omitempty,athletes_search_q_max"`
	Limit  int    `form:"limit" binding:"omitempty,athletes_page_limit"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
}

type GetAthletesV1Response struct {
	Athletes []AthleteV1Response `json:"athletes"`
	Total    int64               `json:"total"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}
