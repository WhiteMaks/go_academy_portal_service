package service

import (
	"context"
	"errors"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestUserService_CreateUserV1_Success(t *testing.T) {
	request := model.PostUserV1Request{
		Username: util.RandomString(64),
		Password: util.RandomString(64),
	}

	expectedResponse := model.PostUserV1Response{
		ID: 1,
	}

	mock := &mockUserRepository{
		createUserV1Func: func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
			err := util.CheckPassword(request.Password, params.PPassword)

			require.NoError(t, err)

			return database.RootUser{ID: expectedResponse.ID}, nil
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	actualResponse, err := service.CreateUserV1(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestUserService_CreateUserV1_Error(t *testing.T) {
	mock := &mockUserRepository{
		createUserV1Func: func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
			return database.RootUser{}, errors.New("repository error")
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	_, err := service.CreateUserV1(
		context.Background(),
		model.PostUserV1Request{},
	)

	require.Error(t, err)
}

func TestUserService_GenerateUserTokenV1_Success(t *testing.T) {
	request := model.PostUserTokenV1Request{
		Username: util.RandomString(64),
		Password: util.RandomString(64),
	}

	hashPassword, _ := util.HashPassword(request.Password)

	userRecord := database.RootUser{
		ID:       1,
		Username: request.Username,
		Password: hashPassword,
		Role:     database.RootUserRoleAthlete,
	}

	mock := &mockUserRepository{
		getUserByUsernameV1Func: func(ctx context.Context, username string) (database.RootUser, error) {
			return userRecord, nil
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	actualResponse, err := service.GenerateUserTokenV1(context.Background(), request)

	require.NoError(t, err)
	require.NotEmpty(t, actualResponse)
	require.NotEmpty(t, actualResponse.Token)
}

func TestUserService_GenerateUserTokenV1_Error(t *testing.T) {
	mock := &mockUserRepository{
		getUserByUsernameV1Func: func(ctx context.Context, username string) (database.RootUser, error) {
			return database.RootUser{}, errors.New("repository error")
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	_, err := service.GenerateUserTokenV1(
		context.Background(),
		model.PostUserTokenV1Request{},
	)

	require.Error(t, err)
}
