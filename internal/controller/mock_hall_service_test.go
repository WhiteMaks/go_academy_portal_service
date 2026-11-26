package controller

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
)

type mockHallService struct {
	createHallV1Func     func(ctx context.Context, request model.PostHallV1Request) (model.HallV1Response, error)
	getBranchHallsV1Func func(ctx context.Context, branchID int64, query model.GetBranchHallsV1Query) (model.GetBranchHallsV1Response, error)
	getHallV1Func        func(ctx context.Context, id int64) (model.HallV1Response, error)
	updateHallV1Func     func(ctx context.Context, id int64, request model.PutHallV1Request) (model.HallV1Response, error)
}

func (m *mockHallService) CreateHallV1(ctx context.Context, request model.PostHallV1Request) (model.HallV1Response, error) {
	return m.createHallV1Func(ctx, request)
}

func (m *mockHallService) GetBranchHallsV1(ctx context.Context, branchID int64, query model.GetBranchHallsV1Query) (model.GetBranchHallsV1Response, error) {
	return m.getBranchHallsV1Func(ctx, branchID, query)
}

func (m *mockHallService) GetHallV1(ctx context.Context, id int64) (model.HallV1Response, error) {
	return m.getHallV1Func(ctx, id)
}

func (m *mockHallService) UpdateHallV1(ctx context.Context, id int64, request model.PutHallV1Request) (model.HallV1Response, error) {
	return m.updateHallV1Func(ctx, id, request)
}
