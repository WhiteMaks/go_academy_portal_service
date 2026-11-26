package controller

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
)

type mockAthleteService struct {
	createAthleteV1Func func(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error)
	getAthletesV1Func   func(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error)
	getAthleteV1Func    func(ctx context.Context, id int64) (model.AthleteV1Response, error)
}

func (m *mockAthleteService) CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
	return m.createAthleteV1Func(ctx, request)
}

func (m *mockAthleteService) GetAthletesV1(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
	return m.getAthletesV1Func(ctx, query)
}

func (m *mockAthleteService) GetAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	return m.getAthleteV1Func(ctx, id)
}
