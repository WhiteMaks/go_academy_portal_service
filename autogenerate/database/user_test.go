package database

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	roles := []RootUserRole{
		RootUserRoleAdmin,
		RootUserRoleCoach,
		RootUserRoleAthlete,
	}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			createUserParams := PrepareRandomUserParams(role, true)

			actualRecord, err := testQueries.CreateUser(context.Background(), createUserParams)
			require.NoError(t, err)
			require.NotEmpty(t, actualRecord)

			expectedRecord := PrepareUserRecord(t, actualRecord, createUserParams)
			require.Equal(t, expectedRecord, actualRecord)
		})
	}
}

func TestCreateUser_UserExist(t *testing.T) {
	createUserParams := PrepareRandomUserParams(RootUserRoleAdmin, true)

	actualRecord, err := testQueries.CreateUser(context.Background(), createUserParams)
	require.NoError(t, err)
	require.NotEmpty(t, actualRecord)

	actualRecord, err = testQueries.CreateUser(context.Background(), createUserParams)
	require.Error(t, err)
	require.Empty(t, actualRecord)

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	require.True(t, ok)
	require.Equal(t, PgErrUniqueViolation, pqErr.Code)
	require.Equal(t, "user with the same username or email already exists", pqErr.Message)
}
