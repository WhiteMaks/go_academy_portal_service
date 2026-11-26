package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/lib/pq"
)

var (
	ErrCoachValidation    = errors.New("coach validation failed")
	ErrCoachNotFound      = errors.New("coach not found")
	ErrCoachPhoneConflict = errors.New("coach with the same phone already exists")
	ErrCoachListQuery     = errors.New("invalid coaches list query")
)

type CoachService interface {
	CreateCoachV1(ctx context.Context, request model.PostCoachV1Request) (model.CoachV1Response, error)
	GetCoachesV1(ctx context.Context, query model.GetCoachesV1Query) (model.GetCoachesV1Response, error)
	GetCoachV1(ctx context.Context, id int64) (model.CoachV1Response, error)
	UpdateCoachV1(ctx context.Context, id int64, request model.PutCoachV1Request) (model.CoachV1Response, error)
}

type coachService struct {
	coachRepository repository.CoachRepository
}

func NewCoachService(coachRepository repository.CoachRepository) CoachService {
	return &coachService{coachRepository: coachRepository}
}

func (s *coachService) CreateCoachV1(ctx context.Context, request model.PostCoachV1Request) (model.CoachV1Response, error) {
	params, err := prepareCreateCoachV1Params(request)
	if err != nil {
		return model.CoachV1Response{}, err
	}

	record, err := s.coachRepository.CreateCoachV1(ctx, params)
	if err != nil {
		return model.CoachV1Response{}, classifyCoachError(err)
	}

	return mapCoachRecordToResponse(record), nil
}

func (s *coachService) GetCoachesV1(ctx context.Context, query model.GetCoachesV1Query) (model.GetCoachesV1Response, error) {
	if err := model.ValidateGetCoachesV1Query(query); err != nil {
		return model.GetCoachesV1Response{}, ErrCoachListQuery
	}

	limit, offset := normalizeCoachesPagination(query.Limit, query.Offset)
	searchQuery := strings.TrimSpace(query.Q)

	listParams := database.GetCoachesV1Params{
		PLimit:  sqlNullInt32(limit),
		POffset: sqlNullInt32(offset),
	}
	if searchQuery != "" {
		listParams.PQ = sql.NullString{String: searchQuery, Valid: true}
	}

	records, err := s.coachRepository.GetCoachesV1(ctx, listParams)
	if err != nil {
		return model.GetCoachesV1Response{}, err
	}

	countQ := sql.NullString{}
	if searchQuery != "" {
		countQ = sql.NullString{String: searchQuery, Valid: true}
	}
	total, err := s.coachRepository.CountCoachesV1(ctx, countQ)
	if err != nil {
		return model.GetCoachesV1Response{}, err
	}

	coaches := make([]model.CoachV1Response, 0, len(records))
	for _, record := range records {
		coaches = append(coaches, mapCoachRecordToResponse(record))
	}

	return model.GetCoachesV1Response{
		Coaches: coaches,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}, nil
}

func (s *coachService) GetCoachV1(ctx context.Context, id int64) (model.CoachV1Response, error) {
	record, err := s.coachRepository.GetCoachV1(ctx, id)
	if err != nil {
		return model.CoachV1Response{}, classifyCoachError(err)
	}
	return mapCoachRecordToResponse(record), nil
}

func (s *coachService) UpdateCoachV1(ctx context.Context, id int64, request model.PutCoachV1Request) (model.CoachV1Response, error) {
	params, err := prepareUpdateCoachV1Params(id, request)
	if err != nil {
		return model.CoachV1Response{}, err
	}

	record, err := s.coachRepository.UpdateCoachV1(ctx, params)
	if err != nil {
		return model.CoachV1Response{}, classifyCoachError(err)
	}

	return mapCoachRecordToResponse(record), nil
}

func prepareUpdateCoachV1Params(id int64, request model.PutCoachV1Request) (database.UpdateCoachV1Params, error) {
	createParams, err := prepareCreateCoachV1Params(request.AsCreateShape())
	if err != nil {
		return database.UpdateCoachV1Params{}, err
	}

	return database.UpdateCoachV1Params{
		PID:             id,
		PLastName:       createParams.PLastName,
		PFirstName:      createParams.PFirstName,
		PPatronymic:     createParams.PPatronymic,
		PPhone:          createParams.PPhone,
		PSpecialization: createParams.PSpecialization,
		PNote:           createParams.PNote,
		PActivityStatus: request.ActivityStatus,
		PBirthDate:      createParams.PBirthDate,
	}, nil
}

func prepareCreateCoachV1Params(request model.PostCoachV1Request) (database.CreateCoachV1Params, error) {
	if err := model.ValidateCoachCreateFieldLengths(request); err != nil {
		return database.CreateCoachV1Params{}, ErrCoachValidation
	}

	lastName := strings.TrimSpace(request.LastName)
	firstName := strings.TrimSpace(request.FirstName)
	specialization := strings.TrimSpace(request.Specialization)
	if lastName == "" || firstName == "" || specialization == "" {
		return database.CreateCoachV1Params{}, ErrCoachValidation
	}

	phone, err := normalizeCoachPhone(request.Phone)
	if err != nil {
		return database.CreateCoachV1Params{}, err
	}

	if request.BirthDate.Time.IsZero() {
		return database.CreateCoachV1Params{}, ErrCoachValidation
	}
	if err := model.ValidateCoachCreateBirthDates(request.BirthDate); err != nil {
		return database.CreateCoachV1Params{}, ErrCoachValidation
	}

	return database.CreateCoachV1Params{
		PLastName:       lastName,
		PFirstName:      firstName,
		PPatronymic:     optionalTrimmedNullString(request.Patronymic),
		PBirthDate:      dateOnlyToTime(request.BirthDate),
		PPhone:          phone,
		PSpecialization: specialization,
		PNote:           optionalTrimmedNullString(request.Note),
	}, nil
}

func mapCoachRecordToResponse(record database.RootCoach) model.CoachV1Response {
	response := model.CoachV1Response{
		ID:             record.ID,
		LastName:       record.LastName,
		FirstName:      record.FirstName,
		BirthDate:      model.DateOnly{Time: record.BirthDate},
		Phone:          record.Phone,
		Specialization: record.Specialization,
		ActivityStatus: string(record.ActivityStatus),
		CreatedAt:      record.CreatedAt,
		UpdatedAt:      record.UpdatedAt,
	}

	if record.Patronymic.Valid {
		value := record.Patronymic.String
		response.Patronymic = &value
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

func classifyCoachError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case database.PgErrNoDataFound:
			return ErrCoachNotFound
		case database.PgErrUniqueViolation:
			return ErrCoachPhoneConflict
		case database.PgErrInvalidParameterValue:
			return mapCoachValidationMessage(pqErr.Message)
		}
	}
	return err
}

func mapCoachValidationMessage(message string) error {
	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "phone"):
		return ErrCoachValidation
	default:
		return ErrCoachValidation
	}
}

func normalizeCoachPhone(raw string) (string, error) {
	normalized, err := util.NormalizePhoneE164(raw)
	if err != nil {
		return "", ErrCoachValidation
	}
	if len(normalized) > model.CoachPhoneMaxLen {
		return "", ErrCoachValidation
	}
	return normalized, nil
}

func normalizeCoachesPagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = model.CoachesListDefaultLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
