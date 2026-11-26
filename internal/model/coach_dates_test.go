package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateCoachCreateBirthDates_Future(t *testing.T) {
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	err := ValidateCoachCreateBirthDates(DateOnly{Time: tomorrow})
	require.ErrorIs(t, err, ErrCoachBirthDateFuture)
}

func TestValidateCoachCreateBirthDates_TodayOk(t *testing.T) {
	today := time.Now().UTC()
	err := ValidateCoachCreateBirthDates(DateOnly{Time: today})
	require.NoError(t, err)
}
