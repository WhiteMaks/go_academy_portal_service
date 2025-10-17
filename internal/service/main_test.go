package service

import (
	"context"
	"database/sql"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"os"
	"testing"
)

type mockUserRepository struct {
	createUserV1Func             func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error)
	getUserByUsernameV1Func      func(ctx context.Context, username string) (database.RootUser, error)
	isAdminCreationAllowedV1Func func(ctx context.Context) (sql.NullBool, error)
}

func (m *mockUserRepository) CreateUserV1(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
	return m.createUserV1Func(ctx, params)
}

func (m *mockUserRepository) GetUserByUsernameV1(ctx context.Context, username string) (database.RootUser, error) {
	return m.getUserByUsernameV1Func(ctx, username)
}

func (m *mockUserRepository) IsAdminCreationAllowedV1(ctx context.Context) (sql.NullBool, error) {
	return m.isAdminCreationAllowedV1Func(ctx)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
