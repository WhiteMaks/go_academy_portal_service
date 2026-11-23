package repository

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type AthleteRepository interface {
	CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error)
	GetAthletesV1(ctx context.Context) ([]database.RootAthlete, error)
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
