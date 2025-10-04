package database

import (
	"testing"

	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/stretchr/testify/require"
)

func PrepareRandomUserParams(role RootUserRole, isActive bool) CreateUserParams {
	return CreateUserParams{
		PUsername: util.RandomString(30),
		PEmail:    util.RandomString(30),
		PPassword: util.RandomString(30),
		PRole:     role,
		PIsActive: isActive,
	}
}

func PrepareUserRecord(t *testing.T, actual RootUser, queryParams CreateUserParams) RootUser {
	actual.Username = queryParams.PUsername
	actual.Email = queryParams.PEmail
	actual.Password = queryParams.PPassword
	actual.Role = queryParams.PRole
	actual.IsActive = queryParams.PIsActive

	require.NotZero(t, actual.ID)
	require.NotZero(t, actual.CreatedAt)
	require.NotZero(t, actual.UpdatedAt)

	return actual
}
