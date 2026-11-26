package service

import (
	"context"
	"database/sql"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
)

type mockCoachRepository struct {
	createCoachV1Func  func(ctx context.Context, params database.CreateCoachV1Params) (database.RootCoach, error)
	getCoachesV1Func   func(ctx context.Context, params database.GetCoachesV1Params) ([]database.RootCoach, error)
	countCoachesV1Func func(ctx context.Context, pQ sql.NullString) (int64, error)
	getCoachV1Func     func(ctx context.Context, id int64) (database.RootCoach, error)
	updateCoachV1Func  func(ctx context.Context, params database.UpdateCoachV1Params) (database.RootCoach, error)
}

func (m *mockCoachRepository) CreateCoachV1(ctx context.Context, params database.CreateCoachV1Params) (database.RootCoach, error) {
	return m.createCoachV1Func(ctx, params)
}

func (m *mockCoachRepository) GetCoachesV1(ctx context.Context, params database.GetCoachesV1Params) ([]database.RootCoach, error) {
	return m.getCoachesV1Func(ctx, params)
}

func (m *mockCoachRepository) CountCoachesV1(ctx context.Context, pQ sql.NullString) (int64, error) {
	return m.countCoachesV1Func(ctx, pQ)
}

func (m *mockCoachRepository) GetCoachV1(ctx context.Context, id int64) (database.RootCoach, error) {
	return m.getCoachV1Func(ctx, id)
}

func (m *mockCoachRepository) UpdateCoachV1(ctx context.Context, params database.UpdateCoachV1Params) (database.RootCoach, error) {
	return m.updateCoachV1Func(ctx, params)
}
