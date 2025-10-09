package util

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateToken_Success(t *testing.T) {
	tokenMaker := NewJWTTokenMaker(RandomString(32))

	username := RandomString(64)
	role := RandomString(5)
	duration := time.Minute

	token, err := tokenMaker.CreateToken(username, role, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := tokenMaker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.IssuedAt)
	require.NotZero(t, payload.ExpiredAt)

	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(duration), payload.ExpiredAt, time.Second)
}

func TestVerifyToken_Expired(t *testing.T) {
	tokenMaker := NewJWTTokenMaker(RandomString(32))

	username := RandomString(64)
	role := RandomString(5)
	duration := -time.Minute

	token, err := tokenMaker.CreateToken(username, role, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := tokenMaker.VerifyToken(token)

	require.Error(t, err)
	require.Empty(t, payload)
	require.EqualError(t, err, ErrTokenExpired.Error())
}
