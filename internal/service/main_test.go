package service

import (
	"context"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"os"
	"testing"
)

type mockUserRepository struct {
	createV1Func func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error)
}

func (m *mockUserRepository) CreateV1(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
	return m.createV1Func(ctx, params)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
