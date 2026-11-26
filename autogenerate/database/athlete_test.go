package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateAthleteV1_AndSearchByBirthDateFragment(t *testing.T) {
	ctx := context.Background()
	phone := "+7900" + time.Now().UTC().Format("150405")
	birthDate := time.Date(2017, 3, 15, 0, 0, 0, 0, time.UTC)

	created, err := testQueries.CreateAthleteV1(ctx, CreateAthleteV1Params{
		PLastName:   "Searchov",
		PFirstName:  "Sergey",
		PBirthDate:  birthDate,
		PPhone:      phone,
		PPatronymic: sql.NullString{},
		PNote:       sql.NullString{},
	})
	require.NoError(t, err)
	require.Equal(t, RootAthleteActivityStatusActive, created.ActivityStatus)
	require.NotZero(t, created.ID)

	duplicate, err := testQueries.CreateAthleteV1(ctx, CreateAthleteV1Params{
		PLastName:  "Other",
		PFirstName: "Athlete",
		PBirthDate: birthDate,
		PPhone:     phone,
	})
	require.Error(t, err)
	var pqErr *pq.Error
	require.ErrorAs(t, err, &pqErr)
	require.Equal(t, PgErrUniqueViolation, pqErr.Code)
	require.Empty(t, duplicate.ID)

	list, err := testQueries.GetAthletesV1(ctx, GetAthletesV1Params{
		PQ:      sql.NullString{String: "15.03", Valid: true},
		PLimit:  sql.NullInt32{Int32: 50, Valid: true},
		POffset: sql.NullInt32{Valid: true},
	})
	require.NoError(t, err)
	require.NotEmpty(t, list)
}
