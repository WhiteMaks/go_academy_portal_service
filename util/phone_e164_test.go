package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizePhoneE164(t *testing.T) {
	t.Parallel()

	phone, err := NormalizePhoneE164("+7 900 123-45-67")
	require.NoError(t, err)
	require.Equal(t, "+79001234567", phone)

	_, err = NormalizePhoneE164("12345")
	require.Error(t, err)

	require.True(t, IsPhoneE164("+79001234567"))
	require.False(t, IsPhoneE164("invalid"))
}
