package service

import (
	"context"
	"database/sql"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type mockBranchRepository struct {
	createBranchV1Func  func(ctx context.Context, params database.CreateBranchV1Params) (database.RootBranch, error)
	getBranchesV1Func   func(ctx context.Context, params database.GetBranchesV1Params) ([]database.RootBranch, error)
	countBranchesV1Func func(ctx context.Context, pQ sql.NullString) (int64, error)
	getBranchV1Func     func(ctx context.Context, id int64) (database.RootBranch, error)
	updateBranchV1Func  func(ctx context.Context, params database.UpdateBranchV1Params) (database.RootBranch, error)
}

func (m *mockBranchRepository) CreateBranchV1(ctx context.Context, params database.CreateBranchV1Params) (database.RootBranch, error) {
	return m.createBranchV1Func(ctx, params)
}

func (m *mockBranchRepository) GetBranchesV1(ctx context.Context, params database.GetBranchesV1Params) ([]database.RootBranch, error) {
	return m.getBranchesV1Func(ctx, params)
}

func (m *mockBranchRepository) CountBranchesV1(ctx context.Context, pQ sql.NullString) (int64, error) {
	return m.countBranchesV1Func(ctx, pQ)
}

func (m *mockBranchRepository) GetBranchV1(ctx context.Context, id int64) (database.RootBranch, error) {
	return m.getBranchV1Func(ctx, id)
}

func (m *mockBranchRepository) UpdateBranchV1(ctx context.Context, params database.UpdateBranchV1Params) (database.RootBranch, error) {
	return m.updateBranchV1Func(ctx, params)
}
