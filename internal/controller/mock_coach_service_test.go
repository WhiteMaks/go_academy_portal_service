package controller

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
)

type mockCoachService struct {
	createCoachV1Func func(ctx context.Context, request model.PostCoachV1Request) (model.CoachV1Response, error)
	getCoachesV1Func  func(ctx context.Context, query model.GetCoachesV1Query) (model.GetCoachesV1Response, error)
	getCoachV1Func    func(ctx context.Context, id int64) (model.CoachV1Response, error)
	updateCoachV1Func func(ctx context.Context, id int64, request model.PutCoachV1Request) (model.CoachV1Response, error)
}

func (m *mockCoachService) CreateCoachV1(ctx context.Context, request model.PostCoachV1Request) (model.CoachV1Response, error) {
	if m.createCoachV1Func == nil {
		return model.CoachV1Response{}, nil
	}
	return m.createCoachV1Func(ctx, request)
}

func (m *mockCoachService) GetCoachesV1(ctx context.Context, query model.GetCoachesV1Query) (model.GetCoachesV1Response, error) {
	if m.getCoachesV1Func == nil {
		return model.GetCoachesV1Response{}, nil
	}
	return m.getCoachesV1Func(ctx, query)
}

func (m *mockCoachService) GetCoachV1(ctx context.Context, id int64) (model.CoachV1Response, error) {
	if m.getCoachV1Func == nil {
		return model.CoachV1Response{}, nil
	}
	return m.getCoachV1Func(ctx, id)
}

func (m *mockCoachService) UpdateCoachV1(ctx context.Context, id int64, request model.PutCoachV1Request) (model.CoachV1Response, error) {
	if m.updateCoachV1Func == nil {
		return model.CoachV1Response{}, nil
	}
	return m.updateCoachV1Func(ctx, id, request)
}
