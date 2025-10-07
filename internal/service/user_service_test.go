package service

import (
	"context"
	"errors"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserRepository_CreateV1_Success(t *testing.T) {
	expectedResponse := model.PostUserV1Response{
		ID: 1,
	}

	mock := &mockUserRepository{
		createV1Func: func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
			return database.RootUser{ID: expectedResponse.ID}, nil
		},
	}

	service := NewUserService(mock)

	actualResponse, err := service.CreateUserV1(
		context.Background(),
		model.PostUserV1Request{},
	)

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestUserRepository_CreateV1_Error(t *testing.T) {
	mock := &mockUserRepository{
		createV1Func: func(ctx context.Context, params database.CreateUserV1Params) (database.RootUser, error) {
			return database.RootUser{}, errors.New("repository error")
		},
	}

	service := NewUserService(mock)

	_, err := service.CreateUserV1(
		context.Background(),
		model.PostUserV1Request{},
	)

	require.Error(t, err)
}
