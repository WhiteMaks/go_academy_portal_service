package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashPassword_Success(t *testing.T) {
	password := RandomString(64)

	hashedPassword, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)
}

func TestHashPassword_Error(t *testing.T) {
	password := RandomString(64)

	hashedPassword, err := HashPassword(password)

	notValidPassword := RandomString(64)

	err = CheckPassword(notValidPassword, hashedPassword)
	require.Error(t, err)
}
