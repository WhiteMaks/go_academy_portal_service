package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type failingUserTokenMaker struct {
	err error
}

func (m failingUserTokenMaker) CreateToken(username string, role string) (string, error) {
	return "", m.err
}

func (m failingUserTokenMaker) VerifyToken(token string) (*util.Payload, error) {
	return nil, m.err
}

func TestUserService_GetUserV1_Success(t *testing.T) {
	userRecord := database.RootUser{
		ID:       util.RandomInt(0, 1_000_000),
		Username: util.RandomString(8),
		Password: "hashed_password",
		Role:     database.RootUserRoleAdmin,
	}

	expectedResponse := model.GetUserV1Response{
		ID:       userRecord.ID,
		Username: userRecord.Username,
		Role:     string(userRecord.Role),
	}

	mock := &mockUserRepository{
		getUserByUsernameV1Func: func(ctx context.Context, username string) (database.RootUser, error) {
			return userRecord, nil
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	actualResponse, err := service.GetUserV1(context.Background(), expectedResponse.Username)

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestUserService_GetUserV1_Error(t *testing.T) {
	mock := &mockUserRepository{
		getUserByUsernameV1Func: func(ctx context.Context, username string) (database.RootUser, error) {
			return database.RootUser{}, errors.New("repository error")
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	_, err := service.GetUserV1(context.Background(), util.RandomString(16))

	require.Error(t, err)
}

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

func TestUserService_CreateAdminV1_Success(t *testing.T) {
	request := model.PostUserV1Request{
		Username: util.RandomString(64),
		Password: util.RandomString(64),
	}

	expectedResponse := model.PostUserV1Response{
		ID: 1,
	}

	mock := &mockUserRepository{
		createAdminV1Func: func(ctx context.Context, username string, password string) (database.RootUser, error) {
			require.Equal(t, request.Username, username)
			err := util.CheckPassword(request.Password, password)

			require.NoError(t, err)

			return database.RootUser{ID: expectedResponse.ID, Role: database.RootUserRoleAdmin, IsActive: true}, nil
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	actualResponse, err := service.CreateAdminV1(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestUserService_CreateAdminV1_AlreadyCreated(t *testing.T) {
	mock := &mockUserRepository{
		createAdminV1Func: func(ctx context.Context, username string, password string) (database.RootUser, error) {
			return database.RootUser{}, &pq.Error{Code: database.PgErrNoDataFound}
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	_, err := service.CreateAdminV1(
		context.Background(),
		model.PostUserV1Request{},
	)

	require.ErrorIs(t, err, ErrAdminCreationIsNotAllowed)
}

func TestUserService_CreateAdminV1_RepositoryError(t *testing.T) {
	mock := &mockUserRepository{
		createAdminV1Func: func(ctx context.Context, username string, password string) (database.RootUser, error) {
			return database.RootUser{}, errors.New("repository error")
		},
	}

	service := NewUserService(mock, util.NewJWTTokenMaker("secret", time.Minute))

	_, err := service.CreateAdminV1(
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
	request := model.PostUserTokenV1Request{
		Username: util.RandomString(16),
		Password: "submitted-password",
	}
	validHash, err := util.HashPassword("stored-password")
	require.NoError(t, err)
	validSubmittedHash, err := util.HashPassword(request.Password)
	require.NoError(t, err)

	tests := []struct {
		name             string
		user             database.RootUser
		repositoryError  error
		tokenMaker       util.TokenMaker
		wantInvalidCreds bool
	}{
		{
			name:             "unknown username",
			repositoryError:  &pq.Error{Code: database.PgErrNoDataFound, Message: "username not found"},
			wantInvalidCreds: true,
		},
		{
			name:             "password mismatch",
			user:             database.RootUser{Username: request.Username, Password: validHash, Role: database.RootUserRoleAdmin},
			wantInvalidCreds: true,
		},
		{
			name:            "repository failure",
			repositoryError: errors.New("database connection detail"),
		},
		{
			name: "invalid stored password hash",
			user: database.RootUser{Username: request.Username, Password: "invalid-hash", Role: database.RootUserRoleAdmin},
		},
		{
			name:       "token signing failure",
			user:       database.RootUser{Username: request.Username, Password: validSubmittedHash, Role: database.RootUserRoleAdmin},
			tokenMaker: failingUserTokenMaker{err: errors.New("token signing detail")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &mockUserRepository{
				getUserByUsernameV1Func: func(ctx context.Context, username string) (database.RootUser, error) {
					return test.user, test.repositoryError
				},
			}
			tokenMaker := test.tokenMaker
			if tokenMaker == nil {
				tokenMaker = util.NewJWTTokenMaker("secret", time.Minute)
			}

			userService := NewUserService(mock, tokenMaker)
			_, err := userService.GenerateUserTokenV1(context.Background(), request)

			require.Error(t, err)
			require.Equal(t, test.wantInvalidCreds, errors.Is(err, ErrInvalidCredentials))
		})
	}
}
