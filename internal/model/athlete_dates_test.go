package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateAthleteCreateBirthDates_FutureAthlete(t *testing.T) {
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	req := PostAthleteV1Request{
		BirthDate: DateOnly{Time: tomorrow},
	}
	err := ValidateAthleteCreateBirthDates(req)
	require.ErrorIs(t, err, ErrAthleteBirthDateFuture)
}

func TestValidateAthleteCreateBirthDates_TodayOk(t *testing.T) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	req := PostAthleteV1Request{
		BirthDate: DateOnly{Time: today},
	}
	require.NoError(t, ValidateAthleteCreateBirthDates(req))
}

func TestValidateAthleteCreateBirthDates_FutureParent(t *testing.T) {
	past := time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)
	future := time.Now().UTC().AddDate(1, 0, 0)
	req := PostAthleteV1Request{
		BirthDate:       DateOnly{Time: past},
		ParentBirthDate: &DateOnly{Time: future},
	}
	err := ValidateAthleteCreateBirthDates(req)
	require.ErrorIs(t, err, ErrAthleteBirthDateFuture)
}
