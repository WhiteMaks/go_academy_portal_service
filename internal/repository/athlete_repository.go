package repository

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type AthleteRepository interface {
	CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error)
	GetAthletesV1(ctx context.Context) ([]database.RootAthlete, error)
	GetAthleteByIDV1(ctx context.Context, id int64) (database.RootAthlete, error)
	UpdateAthleteV1(ctx context.Context, params database.UpdateAthleteV1Params) (database.RootAthlete, error)
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

func (r *athleteRepository) GetAthletesV1(ctx context.Context) ([]database.RootAthlete, error) {
	return r.store.GetAthletesV1(ctx)
}

func (r *athleteRepository) GetAthleteByIDV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return r.store.GetAthleteByIDV1(ctx, id)
}

func (r *athleteRepository) UpdateAthleteV1(ctx context.Context, params database.UpdateAthleteV1Params) (database.RootAthlete, error) {
	return r.store.UpdateAthleteV1(ctx, params)
}
