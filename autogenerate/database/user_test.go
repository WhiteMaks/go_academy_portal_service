package database

import (
	"context"
	"errors"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/lib/pq"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUserV1(t *testing.T) {
	roles := []RootUserRole{
		RootUserRoleAdmin,
		RootUserRoleCoach,
		RootUserRoleAthlete,
	}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			createUserParams := PrepareRandomUserV1Params(role, true)

			actualRecord, err := testQueries.CreateUserV1(context.Background(), createUserParams)
			require.NoError(t, err)
			require.NotEmpty(t, actualRecord)

			expectedRecord := PrepareUserRecord(t, actualRecord, createUserParams)
			require.Equal(t, expectedRecord, actualRecord)
		})
	}
}

func TestCreateUserV1_UserExist(t *testing.T) {
	createUserParams := PrepareRandomUserV1Params(RootUserRoleAdmin, true)

	actualRecord, err := testQueries.CreateUserV1(context.Background(), createUserParams)
	require.NoError(t, err)
	require.NotEmpty(t, actualRecord)

	actualRecord, err = testQueries.CreateUserV1(context.Background(), createUserParams)
	require.Error(t, err)
	require.Empty(t, actualRecord)

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)

	require.True(t, ok)
	require.Equal(t, PgErrUniqueViolation, pqErr.Code)
	require.Equal(t, "user with the same username already exists", pqErr.Message)
}

func TestGetUserByUsernameV1(t *testing.T) {
	createUserParams := PrepareRandomUserV1Params(RootUserRoleAdmin, true)

	expectedRecord, _ := testQueries.CreateUserV1(context.Background(), createUserParams)

	actualRecord, err := testQueries.GetUserByUsernameV1(context.Background(), expectedRecord.Username)

	require.NoError(t, err)
	require.Equal(t, expectedRecord, actualRecord)
}

func TestGetUserByUsernameV1_UserNotFound(t *testing.T) {
	actualRecord, err := testQueries.GetUserByUsernameV1(context.Background(), util.RandomString(30))

	require.Empty(t, actualRecord)

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)

	require.True(t, ok)
	require.Equal(t, PgErrNoDataFound, pqErr.Code)
	require.Equal(t, "user with given username not found", pqErr.Message)
}
