package model

import "errors"

// ErrCoachBirthDateFuture is returned when birth_date is after today.
var ErrCoachBirthDateFuture = errors.New("birth_date cannot be in the future")

// ValidateCoachCreateBirthDates rejects a future coach birth date.
func ValidateCoachCreateBirthDates(birthDate DateOnly) error {
	if !BirthDateNotAfterToday(birthDate) {
		return ErrCoachBirthDateFuture
	}
	return nil
}
