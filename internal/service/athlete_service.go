package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
)

var (
	ErrAthleteNotFound = errors.New("athlete not found")
)

type AthleteService interface {
	CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error)
	GetAthletesV1(ctx context.Context) (model.GetAthletesV1Response, error)
	GetAthleteByIDV1(ctx context.Context, id int64) (model.AthleteV1Response, error)
	UpdateAthleteV1(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error)
}

type athleteService struct {
	repository repository.AthleteRepository
}

func NewAthleteService(repository repository.AthleteRepository) AthleteService {
	return &athleteService{repository: repository}
}

func (s *athleteService) CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
	params := database.CreateAthleteV1Params{
		SecondName:       strings.TrimSpace(request.SecondName),
		FirstName:        strings.TrimSpace(request.FirstName),
		MiddleName:       nullableRequestText(request.MiddleName),
		BirthDate:        request.BirthDate,
		PhotoUrl:         nullableRequestText(request.PhotoURL),
		Phone:            nullableRequestValue(request.Phone),
		ParentSecondName: nullableRequestText(request.ParentSecondName),
		ParentFirstName:  nullableRequestText(request.ParentFirstName),
		ParentMiddleName: nullableRequestText(request.ParentMiddleName),
		ParentPhone:      nullableRequestValue(request.ParentPhone),
	}
	record, err := s.repository.CreateAthleteV1(ctx, params)
	if err != nil {
		return model.AthleteV1Response{}, err
	}
	return athleteResponse(record), nil
}

func (s *athleteService) GetAthletesV1(ctx context.Context) (model.GetAthletesV1Response, error) {
	records, err := s.repository.GetAthletesV1(ctx)
	if err != nil {
		return model.GetAthletesV1Response{}, err
	}
	response := model.GetAthletesV1Response{Athletes: make([]model.AthleteV1Response, 0, len(records))}
	for _, record := range records {
		response.Athletes = append(response.Athletes, athleteResponse(record))
	}
	return response, nil
}

func (s *athleteService) GetAthleteByIDV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	record, err := s.repository.GetAthleteByIDV1(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.AthleteV1Response{}, ErrAthleteNotFound
	}
	if err != nil {
		return model.AthleteV1Response{}, err
	}
	return athleteResponse(record), nil
}

func (s *athleteService) UpdateAthleteV1(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
	params := database.UpdateAthleteV1Params{
		ID:               id,
		SecondName:       strings.TrimSpace(request.SecondName),
		FirstName:        strings.TrimSpace(request.FirstName),
		MiddleName:       nullableRequestText(request.MiddleName),
		BirthDate:        request.BirthDate,
		Phone:            nullableRequestValue(request.Phone),
		ParentSecondName: nullableRequestText(request.ParentSecondName),
		ParentFirstName:  nullableRequestText(request.ParentFirstName),
		ParentMiddleName: nullableRequestText(request.ParentMiddleName),
		ParentPhone:      nullableRequestValue(request.ParentPhone),
	}
	record, err := s.repository.UpdateAthleteV1(ctx, params)
	if errors.Is(err, sql.ErrNoRows) {
		return model.AthleteV1Response{}, ErrAthleteNotFound
	}
	if err != nil {
		return model.AthleteV1Response{}, err
	}
	return athleteResponse(record), nil
}

func nullableRequestText(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

func nullableRequestValue(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func athleteResponse(record database.RootAthlete) model.AthleteV1Response {
	return model.AthleteV1Response{
		ID: strconv.FormatInt(record.ID, 10), SecondName: record.SecondName, FirstName: record.FirstName,
		MiddleName:       nullableString(record.MiddleName),
		BirthDate:        record.BirthDate,
		PhotoURL:         nullableString(record.PhotoUrl),
		Phone:            nullableString(record.Phone),
		ParentFirstName:  nullableString(record.ParentFirstName),
		ParentSecondName: nullableString(record.ParentSecondName),
		ParentMiddleName: nullableString(record.ParentMiddleName),
		ParentPhone:      nullableString(record.ParentPhone),
		IsActive:         record.IsActive,
		UpdatedAt:        record.UpdatedAt,
	}
}

func nullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
