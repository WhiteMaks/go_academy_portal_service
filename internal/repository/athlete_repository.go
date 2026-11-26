package repository

import (
	"context"
	"database/sql"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type AthleteRepository interface {
	CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error)
	GetAthletesV1(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error)
	CountAthletesV1(ctx context.Context, pQ sql.NullString) (int64, error)
	GetAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error)
}

type athleteRepository struct {
	store database.Store
}

func NewAthleteRepository(store database.Store) AthleteRepository {
	return &athleteRepository{store: store}
}

func (r *athleteRepository) CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
	return r.store.CreateAthleteV1(ctx, params)
}

func (r *athleteRepository) GetAthletesV1(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error) {
	return r.store.GetAthletesV1(ctx, params)
}

func (r *athleteRepository) CountAthletesV1(ctx context.Context, pQ sql.NullString) (int64, error) {
	return r.store.CountAthletesV1(ctx, pQ)
}

func (r *athleteRepository) GetAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return r.store.GetAthleteV1(ctx, id)
}
