package repository

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type HallRepository interface {
	CreateHallV1(ctx context.Context, params database.CreateHallV1Params) (database.RootHall, error)
	GetBranchHallsV1(ctx context.Context, params database.GetBranchHallsV1Params) ([]database.RootHall, error)
	CountBranchHallsV1(ctx context.Context, params database.CountBranchHallsV1Params) (int64, error)
	GetHallV1(ctx context.Context, id int64) (database.RootHall, error)
	UpdateHallV1(ctx context.Context, params database.UpdateHallV1Params) (database.RootHall, error)
}

type hallRepository struct {
	store database.Store
}

func NewHallRepository(store database.Store) HallRepository {
	return &hallRepository{store: store}
}

func (r *hallRepository) CreateHallV1(ctx context.Context, params database.CreateHallV1Params) (database.RootHall, error) {
	return r.store.CreateHallV1(ctx, params)
}

func (r *hallRepository) GetBranchHallsV1(ctx context.Context, params database.GetBranchHallsV1Params) ([]database.RootHall, error) {
	return r.store.GetBranchHallsV1(ctx, params)
}

func (r *hallRepository) CountBranchHallsV1(ctx context.Context, params database.CountBranchHallsV1Params) (int64, error) {
	return r.store.CountBranchHallsV1(ctx, params)
}

func (r *hallRepository) GetHallV1(ctx context.Context, id int64) (database.RootHall, error) {
	return r.store.GetHallV1(ctx, id)
}

func (r *hallRepository) UpdateHallV1(ctx context.Context, params database.UpdateHallV1Params) (database.RootHall, error) {
	return r.store.UpdateHallV1(ctx, params)
}
