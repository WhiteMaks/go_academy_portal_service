package repository

import (
	"context"
	"errors"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserRepository_CreateUserV1_Success(t *testing.T) {
	expectedUser := database.RootUser{
		ID:       1,
		Username: util.RandomString(64),
		Password: util.RandomString(255),
		Role:     database.RootUserRoleAthlete,
		IsActive: true,
	}

	mock := &mockStore{
		createUserV1Func: func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
			return expectedUser, nil
		},
	}

	repository := NewUserRepository(mock)

	actualUser, err := repository.CreateUserV1(
		context.Background(),
		database.CreateUserV1Params{
			PUsername: expectedUser.Username,
			PPassword: expectedUser.Password,
			PRole:     expectedUser.Role,
			PIsActive: expectedUser.IsActive,
		},
	)

	require.NoError(t, err)
	require.Equal(t, expectedUser, actualUser)
}

func TestUserRepository_CreateUserV1_Error(t *testing.T) {
	mock := &mockStore{
		createUserV1Func: func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
			return database.RootUser{}, errors.New("db error")
		},
	}

	repo := NewUserRepository(mock)

	_, err := repo.CreateUserV1(
		context.Background(),
		database.CreateUserV1Params{},
	)

	require.Error(t, err)
}
