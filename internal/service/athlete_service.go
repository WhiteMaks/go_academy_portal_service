package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/lib/pq"
)

var (
	ErrAthleteNotFound        = errors.New("athlete not found")
	ErrAthleteAlreadyArchived = errors.New("athlete already archived")
	ErrAthleteAlreadyActive   = errors.New("athlete is already active")
)

type AthleteService interface {
	CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error)
	GetAthletesV1(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error)
	GetAthleteByIDV1(ctx context.Context, id int64) (model.AthleteV1Response, error)
	UpdateAthleteV1(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error)
	ArchiveAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error)
	RestoreAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error)
	SetAthleteCoachV1(ctx context.Context, id int64, isCoach bool) (model.AthleteV1Response, error)
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

func (s *athleteService) GetAthletesV1(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
	listParams := database.GetAthletesV1Params{
		IsActive:         query.IsActive,
		SearchQuery:      query.SearchQuery,
		IsBirthDateQuery: query.IsBirthDateQuery,
		LimitCount:       query.Limit,
		OffsetCount:      query.Offset,
		IsCoach:          optionalBool(query.IsCoach),
	}
	countParams := database.CountAthletesV1Params{
		IsActive:         query.IsActive,
		SearchQuery:      query.SearchQuery,
		IsBirthDateQuery: query.IsBirthDateQuery,
		IsCoach:          optionalBool(query.IsCoach),
	}
	records, err := s.repository.GetAthletesV1(ctx, listParams)
	if err != nil {
		return model.GetAthletesV1Response{}, err
	}
	total, err := s.repository.CountAthletesV1(ctx, countParams)
	if err != nil {
		return model.GetAthletesV1Response{}, err
	}
	response := model.GetAthletesV1Response{
		Athletes: make([]model.AthleteV1Response, 0, len(records)),
		Total:    total,
		Limit:    query.Limit,
		Offset:   query.Offset,
	}
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

func (s *athleteService) ArchiveAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	record, err := s.repository.ArchiveAthleteV1(ctx, id)
	if err != nil {
		return model.AthleteV1Response{}, classifyAthleteMutationError(err, ErrAthleteAlreadyArchived)
	}
	return athleteResponse(record), nil
}

func (s *athleteService) RestoreAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	record, err := s.repository.RestoreAthleteV1(ctx, id)
	if err != nil {
		return model.AthleteV1Response{}, classifyAthleteMutationError(err, ErrAthleteAlreadyActive)
	}
	return athleteResponse(record), nil
}

func (s *athleteService) SetAthleteCoachV1(ctx context.Context, id int64, isCoach bool) (model.AthleteV1Response, error) {
	record, err := s.repository.SetAthleteCoachV1(ctx, database.SetAthleteCoachV1Params{ID: id, IsCoach: isCoach})
	if err != nil {
		return model.AthleteV1Response{}, classifyAthleteNotFoundError(err)
	}
	return athleteResponse(record), nil
}

func classifyAthleteNotFoundError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrAthleteNotFound
	}
	var pqError *pq.Error
	if errors.As(err, &pqError) && string(pqError.Code) == "P0002" {
		return ErrAthleteNotFound
	}
	return err
}

func classifyAthleteMutationError(err error, conflict error) error {
	var pqError *pq.Error
	if errors.As(err, &pqError) {
		switch string(pqError.Code) {
		case "P0002":
			return ErrAthleteNotFound
		case "P4091", "P4092":
			return conflict
		}
	}
	return err
}

func NormalizeAthleteSearchQuery(value string) (string, bool) {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return "", false
	}
	isDateQuery := true
	hasDigit := false
	var dateBuilder strings.Builder
	for _, r := range value {
		switch {
		case r >= '0' && r <= '9':
			hasDigit = true
			dateBuilder.WriteRune(r)
		case r == '.' || unicode.IsSpace(r):
		default:
			isDateQuery = false
		}
	}
	if isDateQuery && hasDigit {
		return dateBuilder.String(), true
	}
	for i, field := range fields {
		fields[i] = strings.Map(func(r rune) rune {
			if r == 'ё' || r == 'Ё' {
				return 'е'
			}
			return unicode.ToLower(r)
		}, field)
	}
	return strings.Join(fields, " "), false
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

func optionalBool(value *bool) sql.NullBool {
	if value == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *value, Valid: true}
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
		IsCoach:          record.IsCoach,
		UpdatedAt:        record.UpdatedAt,
	}
}

func nullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
