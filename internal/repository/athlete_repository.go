package repository

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type AthleteRepository interface {
	CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error)
	GetAthletesV1(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error)
	CountAthletesV1(ctx context.Context, params database.CountAthletesV1Params) (int64, error)
	GetAthleteByIDV1(ctx context.Context, id int64) (database.RootAthlete, error)
	UpdateAthleteV1(ctx context.Context, params database.UpdateAthleteV1Params) (database.RootAthlete, error)
	ArchiveAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error)
	RestoreAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error)
	SetAthleteCoachV1(ctx context.Context, params database.SetAthleteCoachV1Params) (database.RootAthlete, error)
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

func (r *athleteRepository) CountAthletesV1(ctx context.Context, params database.CountAthletesV1Params) (int64, error) {
	return r.store.CountAthletesV1(ctx, params)
}

func (r *athleteRepository) GetAthleteByIDV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return r.store.GetAthleteByIDV1(ctx, id)
}

func (r *athleteRepository) UpdateAthleteV1(ctx context.Context, params database.UpdateAthleteV1Params) (database.RootAthlete, error) {
	return r.store.UpdateAthleteV1(ctx, params)
}

func (r *athleteRepository) ArchiveAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return r.store.ArchiveAthleteV1(ctx, id)
}

func (r *athleteRepository) RestoreAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return r.store.RestoreAthleteV1(ctx, id)
}

func (r *athleteRepository) SetAthleteCoachV1(ctx context.Context, params database.SetAthleteCoachV1Params) (database.RootAthlete, error) {
	return r.store.SetAthleteCoachV1(ctx, params)
}
