package repository

import (
	"context"
	"database/sql"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type CoachRepository interface {
	CreateCoachV1(ctx context.Context, params database.CreateCoachV1Params) (database.RootCoach, error)
	GetCoachesV1(ctx context.Context, params database.GetCoachesV1Params) ([]database.RootCoach, error)
	CountCoachesV1(ctx context.Context, pQ sql.NullString) (int64, error)
	GetCoachV1(ctx context.Context, id int64) (database.RootCoach, error)
	UpdateCoachV1(ctx context.Context, params database.UpdateCoachV1Params) (database.RootCoach, error)
}

type coachRepository struct {
	store database.Store
}

func NewCoachRepository(store database.Store) CoachRepository {
	return &coachRepository{store: store}
}

func (r *coachRepository) CreateCoachV1(ctx context.Context, params database.CreateCoachV1Params) (database.RootCoach, error) {
	return r.store.CreateCoachV1(ctx, params)
}

func (r *coachRepository) GetCoachesV1(ctx context.Context, params database.GetCoachesV1Params) ([]database.RootCoach, error) {
	return r.store.GetCoachesV1(ctx, params)
}

func (r *coachRepository) CountCoachesV1(ctx context.Context, pQ sql.NullString) (int64, error) {
	return r.store.CountCoachesV1(ctx, pQ)
}

func (r *coachRepository) GetCoachV1(ctx context.Context, id int64) (database.RootCoach, error) {
	return r.store.GetCoachV1(ctx, id)
}

func (r *coachRepository) UpdateCoachV1(ctx context.Context, params database.UpdateCoachV1Params) (database.RootCoach, error) {
	return r.store.UpdateCoachV1(ctx, params)
}
