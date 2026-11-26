package service

import (
	"context"
	"database/sql"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type mockAthleteRepository struct {
	createAthleteV1Func func(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error)
	getAthletesV1Func   func(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error)
	countAthletesV1Func func(ctx context.Context, pQ sql.NullString) (int64, error)
	getAthleteV1Func    func(ctx context.Context, id int64) (database.RootAthlete, error)
}

func (m *mockAthleteRepository) CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
	return m.createAthleteV1Func(ctx, params)
}

func (m *mockAthleteRepository) GetAthletesV1(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error) {
	return m.getAthletesV1Func(ctx, params)
}

func (m *mockAthleteRepository) CountAthletesV1(ctx context.Context, pQ sql.NullString) (int64, error) {
	return m.countAthletesV1Func(ctx, pQ)
}

func (m *mockAthleteRepository) GetAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return m.getAthleteV1Func(ctx, id)
}
