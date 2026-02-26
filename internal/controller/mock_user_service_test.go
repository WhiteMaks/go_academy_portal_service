package controller

import (
	"context"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/util"
)

type mockUserService struct {
	getUserV1Func           func(ctx context.Context, username string) (model.GetUserV1Response, error)
	createUserV1Func        func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error)
	createAdminV1Func       func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error)
	generateUserTokenV1Func func(ctx context.Context, request model.PostUserTokenV1Request) (model.PostUserTokenV1Response, error)
	hasAccessFunc           func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool
}

func (m *mockUserService) GetUserV1(ctx context.Context, username string) (model.GetUserV1Response, error) {
	return m.getUserV1Func(ctx, username)
}

func (m *mockUserService) CreateUserV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
	return m.createUserV1Func(ctx, request)
}

func (m *mockUserService) CreateAdminV1(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
	return m.createAdminV1Func(ctx, request)
}

func (m *mockUserService) GenerateUserTokenV1(ctx context.Context, request model.PostUserTokenV1Request) (model.PostUserTokenV1Response, error) {
	return m.generateUserTokenV1Func(ctx, request)
}

func (m *mockUserService) HasAccess(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
	return m.hasAccessFunc(ctx, payload, roles)
}
