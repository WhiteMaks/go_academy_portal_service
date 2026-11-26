package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateCoachV1_AndSearchByBirthDateFragment(t *testing.T) {
	ctx := context.Background()
	phone := "+7900" + time.Now().UTC().Format("1504059")
	birthDate := time.Date(1975, 3, 15, 0, 0, 0, 0, time.UTC)

	created, err := testQueries.CreateCoachV1(ctx, CreateCoachV1Params{
		PLastName:       "Searchov",
		PFirstName:      "Coach",
		PBirthDate:      birthDate,
		PPhone:          phone,
		PPatronymic:     sql.NullString{},
		PSpecialization: "Swimming",
		PNote:           sql.NullString{},
	})
	require.NoError(t, err)
	require.Equal(t, RootCoachActivityStatusActive, created.ActivityStatus)
	require.Equal(t, birthDate.UTC(), created.BirthDate.UTC())
	require.NotZero(t, created.ID)

	duplicate, err := testQueries.CreateCoachV1(ctx, CreateCoachV1Params{
		PLastName:       "Other",
		PFirstName:      "Coach",
		PBirthDate:      birthDate,
		PPhone:          phone,
		PSpecialization: "Boxing",
		PNote:           sql.NullString{},
	})
	require.Error(t, err)
	var pqErr *pq.Error
	require.ErrorAs(t, err, &pqErr)
	require.Equal(t, PgErrUniqueViolation, pqErr.Code)
	require.Empty(t, duplicate.ID)

	list, err := testQueries.GetCoachesV1(ctx, GetCoachesV1Params{
		PQ:      sql.NullString{String: "15.03", Valid: true},
		PLimit:  sql.NullInt32{Int32: 50, Valid: true},
		POffset: sql.NullInt32{Valid: true},
	})
	require.NoError(t, err)
	require.NotEmpty(t, list)
}

func TestCreateCoachV1_FutureBirthDateRejected(t *testing.T) {
	ctx := context.Background()
	future := time.Now().UTC().AddDate(1, 0, 0)
	phone := "+7900" + time.Now().UTC().Format("1504058")

	_, err := testQueries.CreateCoachV1(ctx, CreateCoachV1Params{
		PLastName:       "Future",
		PFirstName:      "Coach",
		PBirthDate:      future,
		PPhone:          phone,
		PSpecialization: "Judo",
	})
	require.Error(t, err)
	var pqErr *pq.Error
	require.ErrorAs(t, err, &pqErr)
	require.Equal(t, PgErrInvalidParameterValue, pqErr.Code)
}
