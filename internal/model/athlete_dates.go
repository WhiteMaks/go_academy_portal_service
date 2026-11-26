package model

import (
	"errors"
	"time"
)

// ErrAthleteBirthDateFuture is returned when birth_date or parent_birth_date is after today.
var ErrAthleteBirthDateFuture = errors.New("birth_date cannot be in the future")

// calendarDateUTC normalizes t to midnight UTC (calendar date from JSON DateOnly).
func calendarDateUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func todayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// BirthDateNotAfterToday reports whether birth_date is today or earlier (UTC calendar).
func BirthDateNotAfterToday(d DateOnly) bool {
	if d.Time.IsZero() {
		return false
	}
	return !calendarDateUTC(d.Time).After(todayUTC())
}

// ValidateAthleteCreateBirthDates rejects future athlete or parent birth dates.
func ValidateAthleteCreateBirthDates(req PostAthleteV1Request) error {
	if !BirthDateNotAfterToday(req.BirthDate) {
		return ErrAthleteBirthDateFuture
	}
	if req.ParentBirthDate != nil && !req.ParentBirthDate.Time.IsZero() {
		if !BirthDateNotAfterToday(*req.ParentBirthDate) {
			return ErrAthleteBirthDateFuture
		}
	}
	return nil
}
