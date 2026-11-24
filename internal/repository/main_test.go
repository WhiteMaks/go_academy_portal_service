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
	getAthletesV1Func            func(ctx context.Context, arg database.GetAthletesV1Params) ([]database.RootAthlete, error)
	countAthletesV1Func          func(ctx context.Context, arg database.CountAthletesV1Params) (int64, error)
	archiveAthleteV1Func         func(ctx context.Context, id int64) (database.RootAthlete, error)
	restoreAthleteV1Func         func(ctx context.Context, id int64) (database.RootAthlete, error)
	getAthleteByIDV1Func         func(ctx context.Context, id int64) (database.RootAthlete, error)
	updateAthleteV1Func          func(ctx context.Context, arg database.UpdateAthleteV1Params) (database.RootAthlete, error)
	getUserByUsernameV1Func      func(ctx context.Context, pUsername string) (database.RootUser, error)
	setAthleteCoachV1Func        func(ctx context.Context, arg database.SetAthleteCoachV1Params) (database.RootAthlete, error)
}

func (m *mockStore) CreateAdminV1(ctx context.Context, arg database.CreateAdminV1Params) (database.RootUser, error) {
	return m.createAdminV1Func(ctx, arg)
}

func (m *mockStore) CreateAthleteV1(ctx context.Context, arg database.CreateAthleteV1Params) (database.RootAthlete, error) {
	return m.createAthleteV1Func(ctx, arg)
}

func (m *mockStore) GetAthletesV1(ctx context.Context, arg database.GetAthletesV1Params) ([]database.RootAthlete, error) {
	return m.getAthletesV1Func(ctx, arg)
}

func (m *mockStore) CountAthletesV1(ctx context.Context, arg database.CountAthletesV1Params) (int64, error) {
	return m.countAthletesV1Func(ctx, arg)
}

func (m *mockStore) ArchiveAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return m.archiveAthleteV1Func(ctx, id)
}

func (m *mockStore) RestoreAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return m.restoreAthleteV1Func(ctx, id)
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

func (m *mockStore) SetAthleteCoachV1(ctx context.Context, arg database.SetAthleteCoachV1Params) (database.RootAthlete, error) {
	return m.setAthleteCoachV1Func(ctx, arg)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
