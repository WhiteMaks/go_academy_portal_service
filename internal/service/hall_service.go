package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/lib/pq"
)

var (
	ErrHallValidation                   = errors.New("hall validation failed")
	ErrHallNotFound                     = errors.New("hall not found")
	ErrHallNameConflict                 = errors.New("hall with the same name already exists in this branch")
	ErrHallListQuery                    = errors.New("invalid halls list query")
	ErrHallCannotActivateInactiveBranch = errors.New("cannot activate hall while branch is inactive")
)

type HallService interface {
	CreateHallV1(ctx context.Context, request model.PostHallV1Request) (model.HallV1Response, error)
	GetBranchHallsV1(ctx context.Context, branchID int64, query model.GetBranchHallsV1Query) (model.GetBranchHallsV1Response, error)
	GetHallV1(ctx context.Context, id int64) (model.HallV1Response, error)
	UpdateHallV1(ctx context.Context, id int64, request model.PutHallV1Request) (model.HallV1Response, error)
}

type hallService struct {
	hallRepository   repository.HallRepository
	branchRepository repository.BranchRepository
}

func NewHallService(hallRepository repository.HallRepository, branchRepository repository.BranchRepository) HallService {
	return &hallService{
		hallRepository:   hallRepository,
		branchRepository: branchRepository,
	}
}

func (s *hallService) CreateHallV1(ctx context.Context, request model.PostHallV1Request) (model.HallV1Response, error) {
	params, err := prepareCreateHallV1Params(request)
	if err != nil {
		return model.HallV1Response{}, err
	}

	record, err := s.hallRepository.CreateHallV1(ctx, params)
	if err != nil {
		return model.HallV1Response{}, classifyHallError(err)
	}

	return s.mapHallRecordToResponse(ctx, record)
}

func (s *hallService) GetBranchHallsV1(ctx context.Context, branchID int64, query model.GetBranchHallsV1Query) (model.GetBranchHallsV1Response, error) {
	if err := model.ValidateGetBranchHallsV1Query(query); err != nil {
		return model.GetBranchHallsV1Response{}, ErrHallListQuery
	}

	limit, offset := normalizeHallsPagination(query.Limit, query.Offset)
	searchQuery := strings.TrimSpace(query.Q)

	listParams := database.GetBranchHallsV1Params{
		PBranchID: branchID,
		PLimit:    sqlNullInt32(limit),
		POffset:   sqlNullInt32(offset),
	}
	countParams := database.CountBranchHallsV1Params{PBranchID: branchID}
	if searchQuery != "" {
		q := sql.NullString{String: searchQuery, Valid: true}
		listParams.PQ = q
		countParams.PQ = q
	}

	records, err := s.hallRepository.GetBranchHallsV1(ctx, listParams)
	if err != nil {
		return model.GetBranchHallsV1Response{}, classifyHallError(err)
	}

	total, err := s.hallRepository.CountBranchHallsV1(ctx, countParams)
	if err != nil {
		return model.GetBranchHallsV1Response{}, classifyHallError(err)
	}

	halls := make([]model.HallV1Response, 0, len(records))
	for _, record := range records {
		response, err := s.mapHallRecordToResponse(ctx, record)
		if err != nil {
			return model.GetBranchHallsV1Response{}, err
		}
		halls = append(halls, response)
	}

	return model.GetBranchHallsV1Response{
		Halls:  halls,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *hallService) GetHallV1(ctx context.Context, id int64) (model.HallV1Response, error) {
	record, err := s.hallRepository.GetHallV1(ctx, id)
	if err != nil {
		return model.HallV1Response{}, classifyHallError(err)
	}
	return s.mapHallRecordToResponse(ctx, record)
}

func (s *hallService) UpdateHallV1(ctx context.Context, id int64, request model.PutHallV1Request) (model.HallV1Response, error) {
	params, err := prepareUpdateHallV1Params(id, request)
	if err != nil {
		return model.HallV1Response{}, err
	}

	record, err := s.hallRepository.UpdateHallV1(ctx, params)
	if err != nil {
		return model.HallV1Response{}, classifyHallError(err)
	}

	return s.mapHallRecordToResponse(ctx, record)
}

func prepareUpdateHallV1Params(id int64, request model.PutHallV1Request) (database.UpdateHallV1Params, error) {
	body, err := prepareHallBodyParams(request.AsCreateShape())
	if err != nil {
		return database.UpdateHallV1Params{}, err
	}

	return database.UpdateHallV1Params{
		PID:             id,
		PName:           body.name,
		PSubtitle:       body.subtitle,
		PCapacity:       body.capacity,
		PNote:           body.note,
		PActivityStatus: request.ActivityStatus,
	}, nil
}

type hallBodyParams struct {
	name     string
	subtitle sql.NullString
	capacity int32
	note     sql.NullString
}

// prepareHallBodyParams — общая валидация полей зала (create и update); branch_id только для create.
func prepareHallBodyParams(request model.PostHallV1Request) (hallBodyParams, error) {
	if err := model.ValidateHallCreateFieldLengths(request); err != nil {
		return hallBodyParams{}, ErrHallValidation
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		return hallBodyParams{}, ErrHallValidation
	}

	return hallBodyParams{
		name:     name,
		subtitle: optionalTrimmedNullString(request.Subtitle),
		capacity: int32(request.Capacity),
		note:     optionalTrimmedNullString(request.Note),
	}, nil
}

func prepareCreateHallV1Params(request model.PostHallV1Request) (database.CreateHallV1Params, error) {
	if request.BranchID <= 0 {
		return database.CreateHallV1Params{}, ErrHallValidation
	}
	body, err := prepareHallBodyParams(request)
	if err != nil {
		return database.CreateHallV1Params{}, err
	}

	return database.CreateHallV1Params{
		PBranchID: request.BranchID,
		PName:     body.name,
		PSubtitle: body.subtitle,
		PCapacity: body.capacity,
		PNote:     body.note,
	}, nil
}

func (s *hallService) mapHallRecordToResponse(ctx context.Context, record database.RootHall) (model.HallV1Response, error) {
	response := model.HallV1Response{
		ID:             record.ID,
		BranchID:       record.BranchID,
		Name:           record.Name,
		Capacity:       int(record.Capacity),
		ActivityStatus: string(record.ActivityStatus),
		CreatedAt:      record.CreatedAt,
		UpdatedAt:      record.UpdatedAt,
	}

	if record.Subtitle.Valid {
		value := record.Subtitle.String
		response.Subtitle = &value
	}
	if record.Note.Valid {
		value := record.Note.String
		response.Note = &value
	}
	if record.ArchivedAt.Valid {
		value := record.ArchivedAt.Time
		response.ArchivedAt = &value
	}

	branch, err := s.branchRepository.GetBranchV1(ctx, record.BranchID)
	if err == nil {
		name := branch.Name
		address := branch.Address
		branchStatus := string(branch.ActivityStatus)
		response.BranchName = &name
		response.BranchAddress = &address
		response.BranchActivityStatus = &branchStatus
	}

	return response, nil
}

func classifyHallError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case database.PgErrNoDataFound:
			if strings.Contains(strings.ToLower(pqErr.Message), "branch") {
				return ErrBranchNotFound
			}
			return ErrHallNotFound
		case database.PgErrUniqueViolation:
			return ErrHallNameConflict
		case database.PgErrInvalidParameterValue:
			if strings.Contains(strings.ToLower(pqErr.Message), "cannot activate hall while branch is inactive") {
				return ErrHallCannotActivateInactiveBranch
			}
			return ErrHallValidation
		}
	}
	return err
}

func normalizeHallsPagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = model.HallsListDefaultLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
