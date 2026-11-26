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
	ErrBranchValidation   = errors.New("branch validation failed")
	ErrBranchNotFound     = errors.New("branch not found")
	ErrBranchNameConflict = errors.New("branch with the same name already exists")
	ErrBranchListQuery    = errors.New("invalid branches list query")
)

type BranchService interface {
	CreateBranchV1(ctx context.Context, request model.PostBranchV1Request) (model.BranchV1Response, error)
	GetBranchesV1(ctx context.Context, query model.GetBranchesV1Query) (model.GetBranchesV1Response, error)
	GetBranchV1(ctx context.Context, id int64) (model.BranchV1Response, error)
	UpdateBranchV1(ctx context.Context, id int64, request model.PutBranchV1Request) (model.BranchV1Response, error)
}

type branchService struct {
	branchRepository repository.BranchRepository
}

func NewBranchService(branchRepository repository.BranchRepository) BranchService {
	return &branchService{branchRepository: branchRepository}
}

func (s *branchService) CreateBranchV1(ctx context.Context, request model.PostBranchV1Request) (model.BranchV1Response, error) {
	params, err := prepareCreateBranchV1Params(request)
	if err != nil {
		return model.BranchV1Response{}, err
	}

	record, err := s.branchRepository.CreateBranchV1(ctx, params)
	if err != nil {
		return model.BranchV1Response{}, classifyBranchError(err)
	}

	return mapBranchRecordToResponse(record), nil
}

func (s *branchService) GetBranchesV1(ctx context.Context, query model.GetBranchesV1Query) (model.GetBranchesV1Response, error) {
	if err := model.ValidateGetBranchesV1Query(query); err != nil {
		return model.GetBranchesV1Response{}, ErrBranchListQuery
	}

	limit, offset := normalizeBranchesPagination(query.Limit, query.Offset)
	searchQuery := strings.TrimSpace(query.Q)

	listParams := database.GetBranchesV1Params{
		PLimit:  sqlNullInt32(limit),
		POffset: sqlNullInt32(offset),
	}
	if searchQuery != "" {
		listParams.PQ = sql.NullString{String: searchQuery, Valid: true}
	}

	records, err := s.branchRepository.GetBranchesV1(ctx, listParams)
	if err != nil {
		return model.GetBranchesV1Response{}, err
	}

	countQ := sql.NullString{}
	if searchQuery != "" {
		countQ = sql.NullString{String: searchQuery, Valid: true}
	}
	total, err := s.branchRepository.CountBranchesV1(ctx, countQ)
	if err != nil {
		return model.GetBranchesV1Response{}, err
	}

	branches := make([]model.BranchV1Response, 0, len(records))
	for _, record := range records {
		branches = append(branches, mapBranchRecordToResponse(record))
	}

	return model.GetBranchesV1Response{
		Branches: branches,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (s *branchService) GetBranchV1(ctx context.Context, id int64) (model.BranchV1Response, error) {
	record, err := s.branchRepository.GetBranchV1(ctx, id)
	if err != nil {
		return model.BranchV1Response{}, classifyBranchError(err)
	}
	return mapBranchRecordToResponse(record), nil
}

func (s *branchService) UpdateBranchV1(ctx context.Context, id int64, request model.PutBranchV1Request) (model.BranchV1Response, error) {
	params, err := prepareUpdateBranchV1Params(id, request)
	if err != nil {
		return model.BranchV1Response{}, err
	}

	record, err := s.branchRepository.UpdateBranchV1(ctx, params)
	if err != nil {
		return model.BranchV1Response{}, classifyBranchError(err)
	}

	return mapBranchRecordToResponse(record), nil
}

func prepareUpdateBranchV1Params(id int64, request model.PutBranchV1Request) (database.UpdateBranchV1Params, error) {
	createParams, err := prepareCreateBranchV1Params(request.AsCreateShape())
	if err != nil {
		return database.UpdateBranchV1Params{}, err
	}

	return database.UpdateBranchV1Params{
		PID:             id,
		PName:           createParams.PName,
		PAddress:        createParams.PAddress,
		PNote:           createParams.PNote,
		PActivityStatus: request.ActivityStatus,
	}, nil
}

func prepareCreateBranchV1Params(request model.PostBranchV1Request) (database.CreateBranchV1Params, error) {
	if err := model.ValidateBranchCreateFieldLengths(request); err != nil {
		return database.CreateBranchV1Params{}, ErrBranchValidation
	}

	name := strings.TrimSpace(request.Name)
	address := strings.TrimSpace(request.Address)
	if name == "" || address == "" {
		return database.CreateBranchV1Params{}, ErrBranchValidation
	}

	return database.CreateBranchV1Params{
		PName:    name,
		PAddress: address,
		PNote:    optionalTrimmedNullString(request.Note),
	}, nil
}

func mapBranchRecordToResponse(record database.RootBranch) model.BranchV1Response {
	response := model.BranchV1Response{
		ID:             record.ID,
		Name:           record.Name,
		Address:        record.Address,
		ActivityStatus: string(record.ActivityStatus),
		CreatedAt:      record.CreatedAt,
		UpdatedAt:      record.UpdatedAt,
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

func classifyBranchError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case database.PgErrNoDataFound:
			return ErrBranchNotFound
		case database.PgErrUniqueViolation:
			return ErrBranchNameConflict
		case database.PgErrInvalidParameterValue:
			return ErrBranchValidation
		}
	}
	return err
}

func normalizeBranchesPagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = model.BranchesListDefaultLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
