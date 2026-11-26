package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/lib/pq"
)

var (
	ErrAthleteValidation    = errors.New("athlete validation failed")
	ErrAthleteNotFound      = errors.New("athlete not found")
	ErrAthletePhoneConflict = errors.New("athlete with the same phone already exists")
	ErrAthleteParentPartial = errors.New("fill all parent fields or leave the parent block empty")
	ErrAthleteListQuery     = errors.New("invalid athletes list query")
)

type AthleteService interface {
	CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error)
	GetAthletesV1(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error)
	GetAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error)
}

type athleteService struct {
	athleteRepository repository.AthleteRepository
}

func NewAthleteService(athleteRepository repository.AthleteRepository) AthleteService {
	return &athleteService{athleteRepository: athleteRepository}
}

func (s *athleteService) CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
	params, err := prepareCreateAthleteV1Params(request)
	if err != nil {
		return model.AthleteV1Response{}, err
	}

	record, err := s.athleteRepository.CreateAthleteV1(ctx, params)
	if err != nil {
		return model.AthleteV1Response{}, classifyAthleteError(err)
	}

	return mapAthleteRecordToResponse(record), nil
}

func (s *athleteService) GetAthletesV1(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
	if err := model.ValidateGetAthletesV1Query(query); err != nil {
		return model.GetAthletesV1Response{}, ErrAthleteListQuery
	}

	limit, offset := normalizeAthletesPagination(query.Limit, query.Offset)
	searchQuery := strings.TrimSpace(query.Q)

	listParams := database.GetAthletesV1Params{
		PLimit:  sqlNullInt32(limit),
		POffset: sqlNullInt32(offset),
	}
	if searchQuery != "" {
		listParams.PQ = sql.NullString{String: searchQuery, Valid: true}
	}

	records, err := s.athleteRepository.GetAthletesV1(ctx, listParams)
	if err != nil {
		return model.GetAthletesV1Response{}, err
	}

	countQ := sql.NullString{}
	if searchQuery != "" {
		countQ = sql.NullString{String: searchQuery, Valid: true}
	}
	total, err := s.athleteRepository.CountAthletesV1(ctx, countQ)
	if err != nil {
		return model.GetAthletesV1Response{}, err
	}

	athletes := make([]model.AthleteV1Response, 0, len(records))
	for _, record := range records {
		athletes = append(athletes, mapAthleteRecordToResponse(record))
	}

	return model.GetAthletesV1Response{
		Athletes: athletes,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (s *athleteService) GetAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	record, err := s.athleteRepository.GetAthleteV1(ctx, id)
	if err != nil {
		return model.AthleteV1Response{}, classifyAthleteError(err)
	}
	return mapAthleteRecordToResponse(record), nil
}

func prepareCreateAthleteV1Params(request model.PostAthleteV1Request) (database.CreateAthleteV1Params, error) {
	if err := model.ValidateAthleteCreateFieldLengths(request); err != nil {
		return database.CreateAthleteV1Params{}, ErrAthleteValidation
	}

	lastName := strings.TrimSpace(request.LastName)
	firstName := strings.TrimSpace(request.FirstName)
	if lastName == "" || firstName == "" {
		return database.CreateAthleteV1Params{}, ErrAthleteValidation
	}

	phone, err := normalizeAthletePhone(request.Phone)
	if err != nil {
		return database.CreateAthleteV1Params{}, err
	}

	if request.BirthDate.Time.IsZero() {
		return database.CreateAthleteV1Params{}, ErrAthleteValidation
	}
	if err := model.ValidateAthleteCreateBirthDates(request); err != nil {
		return database.CreateAthleteV1Params{}, ErrAthleteValidation
	}

	parentPartial := isParentFieldSet(request.ParentLastName) ||
		isParentFieldSet(request.ParentFirstName) ||
		isParentFieldSet(request.ParentPatronymic) ||
		request.ParentBirthDate != nil && !request.ParentBirthDate.Time.IsZero() ||
		isParentFieldSet(request.ParentPhone)
	parentComplete := isParentFieldSet(request.ParentLastName) &&
		isParentFieldSet(request.ParentFirstName) &&
		request.ParentBirthDate != nil && !request.ParentBirthDate.Time.IsZero() &&
		isParentFieldSet(request.ParentPhone)

	if parentPartial && !parentComplete {
		return database.CreateAthleteV1Params{}, ErrAthleteParentPartial
	}

	params := database.CreateAthleteV1Params{
		PLastName:   lastName,
		PFirstName:  firstName,
		PBirthDate:  dateOnlyToTime(request.BirthDate),
		PPhone:      phone,
		PPatronymic: optionalTrimmedNullString(request.Patronymic),
		PNote:       optionalTrimmedNullString(request.Note),
	}

	if parentComplete {
		params.PParentLastName = optionalTrimmedNullString(request.ParentLastName)
		params.PParentFirstName = optionalTrimmedNullString(request.ParentFirstName)
		params.PParentPatronymic = optionalTrimmedNullString(request.ParentPatronymic)
		parentBirth := dateOnlyToTime(*request.ParentBirthDate)
		params.PParentBirthDate = sql.NullTime{Time: parentBirth, Valid: true}
		parentPhone, phoneErr := normalizeAthletePhone(ptrToString(request.ParentPhone))
		if phoneErr != nil {
			return database.CreateAthleteV1Params{}, phoneErr
		}
		params.PParentPhone = sql.NullString{String: parentPhone, Valid: true}
	}

	return params, nil
}

func mapAthleteRecordToResponse(record database.RootAthlete) model.AthleteV1Response {
	response := model.AthleteV1Response{
		ID:             record.ID,
		LastName:       record.LastName,
		FirstName:      record.FirstName,
		BirthDate:      model.DateOnly{Time: record.BirthDate},
		Phone:          record.Phone,
		ActivityStatus: string(record.ActivityStatus),
		CreatedAt:      record.CreatedAt,
		UpdatedAt:      record.UpdatedAt,
	}

	if record.Patronymic.Valid {
		value := record.Patronymic.String
		response.Patronymic = &value
	}
	if record.ParentLastName.Valid {
		value := record.ParentLastName.String
		response.ParentLastName = &value
	}
	if record.ParentFirstName.Valid {
		value := record.ParentFirstName.String
		response.ParentFirstName = &value
	}
	if record.ParentPatronymic.Valid {
		value := record.ParentPatronymic.String
		response.ParentPatronymic = &value
	}
	if record.ParentBirthDate.Valid {
		response.ParentBirthDate = &model.DateOnly{Time: record.ParentBirthDate.Time}
	}
	if record.ParentPhone.Valid {
		value := record.ParentPhone.String
		response.ParentPhone = &value
	}
	if record.Note.Valid {
		value := record.Note.String
		response.Note = &value
	}
	if record.ArchivedAt.Valid {
		value := record.ArchivedAt.Time
		response.ArchivedAt = &value
	}

	return response
}

func classifyAthleteError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case database.PgErrNoDataFound:
			return ErrAthleteNotFound
		case database.PgErrUniqueViolation:
			return ErrAthletePhoneConflict
		case database.PgErrInvalidParameterValue:
			return mapAthleteValidationMessage(pqErr.Message)
		}
	}
	return err
}

func mapAthleteValidationMessage(message string) error {
	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "parent"):
		return ErrAthleteParentPartial
	case strings.Contains(lower, "phone"):
		return ErrAthleteValidation
	default:
		return ErrAthleteValidation
	}
}

func normalizeAthletePhone(raw string) (string, error) {
	normalized, err := util.NormalizePhoneE164(raw)
	if err != nil {
		return "", ErrAthleteValidation
	}
	if len(normalized) > model.AthletePhoneMaxLen {
		return "", ErrAthleteValidation
	}
	return normalized, nil
}

func normalizeAthletesPagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = model.AthletesListDefaultLimit
	} else if limit > model.AthletesListMaxLimit {
		limit = model.AthletesListMaxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func optionalTrimmedNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

func sqlNullInt32(value int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(value), Valid: true}
}

func isParentFieldSet(value *string) bool {
	return value != nil && strings.TrimSpace(*value) != ""
}

func ptrToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func dateOnlyToTime(date model.DateOnly) time.Time {
	t := date.Time
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
