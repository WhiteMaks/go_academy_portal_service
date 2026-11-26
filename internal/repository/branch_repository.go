package repository

import (
	"context"
	"database/sql"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type BranchRepository interface {
	CreateBranchV1(ctx context.Context, params database.CreateBranchV1Params) (database.RootBranch, error)
	GetBranchesV1(ctx context.Context, params database.GetBranchesV1Params) ([]database.RootBranch, error)
	CountBranchesV1(ctx context.Context, pQ sql.NullString) (int64, error)
	GetBranchV1(ctx context.Context, id int64) (database.RootBranch, error)
	UpdateBranchV1(ctx context.Context, params database.UpdateBranchV1Params) (database.RootBranch, error)
}

type branchRepository struct {
	store database.Store
}

func NewBranchRepository(store database.Store) BranchRepository {
	return &branchRepository{store: store}
}

func (r *branchRepository) CreateBranchV1(ctx context.Context, params database.CreateBranchV1Params) (database.RootBranch, error) {
	return r.store.CreateBranchV1(ctx, params)
}

func (r *branchRepository) GetBranchesV1(ctx context.Context, params database.GetBranchesV1Params) ([]database.RootBranch, error) {
	return r.store.GetBranchesV1(ctx, params)
}

func (r *branchRepository) CountBranchesV1(ctx context.Context, pQ sql.NullString) (int64, error) {
	return r.store.CountBranchesV1(ctx, pQ)
}

func (r *branchRepository) GetBranchV1(ctx context.Context, id int64) (database.RootBranch, error) {
	return r.store.GetBranchV1(ctx, id)
}

func (r *branchRepository) UpdateBranchV1(ctx context.Context, params database.UpdateBranchV1Params) (database.RootBranch, error) {
	return r.store.UpdateBranchV1(ctx, params)
}
