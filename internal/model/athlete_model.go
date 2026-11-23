package model

import "time"

type PostAthleteV1Request struct {
	FirstName        string    `json:"firstName" binding:"required,min=1,max=100"`
	SecondName       string    `json:"secondName" binding:"required,min=1,max=100"`
	MiddleName       *string   `json:"middleName" binding:"omitempty,max=100"`
	BirthDate        time.Time `json:"birthDate" binding:"required"`
	PhotoURL         *string   `json:"photoUrl" binding:"omitempty,max=2048,http_url"`
	Phone            string    `json:"phone" binding:"omitempty,min=11,max=16,e164,excludes=+0"`
	ParentFirstName  *string   `json:"parentFirstName" binding:"omitempty,max=100"`
	ParentSecondName *string   `json:"parentSecondName" binding:"omitempty,max=100"`
	ParentMiddleName *string   `json:"parentMiddleName" binding:"omitempty,max=100"`
	ParentPhone      string    `json:"parentPhone" binding:"omitempty,min=11,max=16,e164,excludes=+0"`
}

type AthleteV1Response struct {
	ID               string    `json:"id"`
	FirstName        string    `json:"firstName"`
	SecondName       string    `json:"secondName"`
	MiddleName       *string   `json:"middleName"`
	BirthDate        time.Time `json:"birthDate"`
	PhotoURL         *string   `json:"photoUrl"`
	Phone            *string   `json:"phone"`
	ParentFirstName  *string   `json:"parentFirstName"`
	ParentSecondName *string   `json:"parentSecondName"`
	ParentMiddleName *string   `json:"parentMiddleName"`
	ParentPhone      *string   `json:"parentPhone"`
	IsActive         bool      `json:"isActive"`
}

type GetAthletesV1Response struct {
	Athletes []AthleteV1Response `json:"athletes"`
}
