package repository

import (
	"context"
	"database/sql"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"os"
	"testing"
)

type mockStore struct {
	createAdminV1Func            func(ctx context.Context, arg database.CreateAdminV1Params) (database.RootUser, error)
	createAthleteV1Func          func(ctx context.Context, arg database.CreateAthleteV1Params) (database.RootAthlete, error)
	createUserV1Func             func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error)
	isAdminCreationAllowedV1Func func(ctx context.Context) (sql.NullBool, error)
	getAthletesV1Func            func(ctx context.Context) ([]database.RootAthlete, error)
	getAthleteByIDV1Func         func(ctx context.Context, id int64) (database.RootAthlete, error)
	updateAthleteV1Func          func(ctx context.Context, arg database.UpdateAthleteV1Params) (database.RootAthlete, error)
	getUserByUsernameV1Func      func(ctx context.Context, pUsername string) (database.RootUser, error)
}

func (m *mockStore) CreateAdminV1(ctx context.Context, arg database.CreateAdminV1Params) (database.RootUser, error) {
	return m.createAdminV1Func(ctx, arg)
}

func (m *mockStore) CreateAthleteV1(ctx context.Context, arg database.CreateAthleteV1Params) (database.RootAthlete, error) {
	return m.createAthleteV1Func(ctx, arg)
}

func (m *mockStore) GetAthletesV1(ctx context.Context) ([]database.RootAthlete, error) {
	return m.getAthletesV1Func(ctx)
}

func (m *mockStore) GetAthleteByIDV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return m.getAthleteByIDV1Func(ctx, id)
}

func (m *mockStore) UpdateAthleteV1(ctx context.Context, arg database.UpdateAthleteV1Params) (database.RootAthlete, error) {
	return m.updateAthleteV1Func(ctx, arg)
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
