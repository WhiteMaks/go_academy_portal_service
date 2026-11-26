package service

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type mockHallRepository struct {
	createHallV1Func     func(ctx context.Context, params database.CreateHallV1Params) (database.RootHall, error)
	getBranchHallsV1Func func(ctx context.Context, params database.GetBranchHallsV1Params) ([]database.RootHall, error)
	countBranchHallsFunc func(ctx context.Context, params database.CountBranchHallsV1Params) (int64, error)
	getHallV1Func        func(ctx context.Context, id int64) (database.RootHall, error)
	updateHallV1Func     func(ctx context.Context, params database.UpdateHallV1Params) (database.RootHall, error)
}

func (m *mockHallRepository) CreateHallV1(ctx context.Context, params database.CreateHallV1Params) (database.RootHall, error) {
	return m.createHallV1Func(ctx, params)
}

func (m *mockHallRepository) GetBranchHallsV1(ctx context.Context, params database.GetBranchHallsV1Params) ([]database.RootHall, error) {
	return m.getBranchHallsV1Func(ctx, params)
}

func (m *mockHallRepository) CountBranchHallsV1(ctx context.Context, params database.CountBranchHallsV1Params) (int64, error) {
	return m.countBranchHallsFunc(ctx, params)
}

func (m *mockHallRepository) GetHallV1(ctx context.Context, id int64) (database.RootHall, error) {
	return m.getHallV1Func(ctx, id)
}

func (m *mockHallRepository) UpdateHallV1(ctx context.Context, params database.UpdateHallV1Params) (database.RootHall, error) {
	return m.updateHallV1Func(ctx, params)
}
