package repository

import (
	"context"
	"database/sql"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"os"
	"testing"
)

type mockStore struct {
	createUserV1Func             func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error)
	isAdminCreationAllowedV1Func func(ctx context.Context) (sql.NullBool, error)
	getUserByUsernameV1Func      func(ctx context.Context, pUsername string) (database.RootUser, error)
}

func (m *mockStore) CreateUserV1(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
	return m.createUserV1Func(ctx, params)
}

func (m *mockStore) IsAdminCreationAllowedV1(ctx context.Context) (sql.NullBool, error) {
	return m.isAdminCreationAllowedV1Func(ctx)
}

func (m *mockStore) GetUserByUsernameV1(ctx context.Context, pUsername string) (database.RootUser, error) {
	return m.getUserByUsernameV1Func(ctx, pUsername)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
