package model

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestAthleteNameMaxValidator_usesAthleteNameMaxLen(t *testing.T) {
	t.Parallel()

	v := validator.New()
	v.SetTagName("binding")
	require.NoError(t, RegisterAthleteBindingValidators(v))

	tooLong := strings.Repeat("а", AthleteNameMaxLen+1)
	err := v.Struct(PostAthleteV1Request{
		LastName:  tooLong,
		FirstName: "Иван",
		BirthDate: DateOnly{},
		Phone:     "+79001234567",
	})
	require.Error(t, err)

	okLen := strings.Repeat("а", AthleteNameMaxLen)
	err = v.Struct(PostAthleteV1Request{
		LastName:  okLen,
		FirstName: "Иван",
		BirthDate: DateOnly{},
		Phone:     "+79001234567",
	})
	require.NoError(t, err)
}

func TestValidateGetAthletesV1Query_limits(t *testing.T) {
	t.Parallel()

	require.NoError(t, ValidateGetAthletesV1Query(GetAthletesV1Query{Limit: 0, Offset: 0}))
	require.NoError(t, ValidateGetAthletesV1Query(GetAthletesV1Query{Limit: AthletesListMaxLimit, Offset: 10}))
	require.Error(t, ValidateGetAthletesV1Query(GetAthletesV1Query{Limit: -1}))
	require.Error(t, ValidateGetAthletesV1Query(GetAthletesV1Query{Limit: AthletesListMaxLimit + 1}))
	require.Error(t, ValidateGetAthletesV1Query(GetAthletesV1Query{Offset: -1}))
}

func TestAthletesPageLimitValidator_usesConstant(t *testing.T) {
	t.Parallel()

	v := validator.New()
	v.SetTagName("binding")
	require.NoError(t, RegisterAthleteBindingValidators(v))

	err := v.Struct(GetAthletesV1Query{Limit: AthletesListMaxLimit + 1})
	require.Error(t, err)
}

func TestValidateGetAthletesV1Query_searchQTooLong(t *testing.T) {
	t.Parallel()

	q := strings.Repeat("q", AthletesSearchQMaxLen+1)
	require.Error(t, ValidateGetAthletesV1Query(GetAthletesV1Query{Q: q}))
}

func TestAthletesSearchQMaxValidator_usesConstant(t *testing.T) {
	t.Parallel()

	v := validator.New()
	v.SetTagName("binding")
	require.NoError(t, RegisterAthleteBindingValidators(v))

	err := v.Struct(GetAthletesV1Query{Q: strings.Repeat("z", AthletesSearchQMaxLen+1)})
	require.Error(t, err)
}

func TestValidateAthleteCreateFieldLengths_note(t *testing.T) {
	t.Parallel()

	note := strings.Repeat("x", AthleteNoteMaxLen+1)
	req := PostAthleteV1Request{
		LastName:  "Ivanov",
		FirstName: "Ivan",
		Note:      &note,
	}
	require.Error(t, ValidateAthleteCreateFieldLengths(req))
}
