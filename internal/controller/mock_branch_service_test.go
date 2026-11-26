package controller

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
)

type mockBranchService struct {
	createBranchV1Func func(ctx context.Context, request model.PostBranchV1Request) (model.BranchV1Response, error)
	getBranchesV1Func  func(ctx context.Context, query model.GetBranchesV1Query) (model.GetBranchesV1Response, error)
	getBranchV1Func    func(ctx context.Context, id int64) (model.BranchV1Response, error)
	updateBranchV1Func func(ctx context.Context, id int64, request model.PutBranchV1Request) (model.BranchV1Response, error)
}

func (m *mockBranchService) CreateBranchV1(ctx context.Context, request model.PostBranchV1Request) (model.BranchV1Response, error) {
	return m.createBranchV1Func(ctx, request)
}

func (m *mockBranchService) GetBranchesV1(ctx context.Context, query model.GetBranchesV1Query) (model.GetBranchesV1Response, error) {
	return m.getBranchesV1Func(ctx, query)
}

func (m *mockBranchService) GetBranchV1(ctx context.Context, id int64) (model.BranchV1Response, error) {
	return m.getBranchV1Func(ctx, id)
}

func (m *mockBranchService) UpdateBranchV1(ctx context.Context, id int64, request model.PutBranchV1Request) (model.BranchV1Response, error) {
	return m.updateBranchV1Func(ctx, id, request)
}
