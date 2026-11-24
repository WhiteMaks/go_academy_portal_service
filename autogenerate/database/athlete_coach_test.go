package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/stretchr/testify/require"
)

func TestAthleteCoachSQL_DisposableDatabase(t *testing.T) {
	db, _ := openAthleteFeatureTestDB(t)
	tx, txQueries := beginAthleteFeatureTestTx(t, db)
	ctx := context.Background()
	var usersBefore int64
	require.NoError(t, tx.QueryRowContext(ctx, "SELECT count(*) FROM root.users").Scan(&usersBefore))
	suffix := util.RandomString(12)
	active := createFeatureTestAthlete(t, txQueries, "КоучФильтр"+suffix, "Анна"+suffix, time.Date(2012, 4, 23, 0, 0, 0, 0, time.UTC))
	archived := createFeatureTestAthlete(t, txQueries, "КоучАрхив"+suffix, "Борис"+suffix, time.Date(2011, 3, 12, 0, 0, 0, 0, time.UTC))

	require.False(t, active.IsCoach, "new athlete defaults to athlete classification")
	require.NoError(t, func() error {
		_, err := txQueries.ArchiveAthleteV1(ctx, archived.ID)
		return err
	}())

	updated, err := txQueries.SetAthleteCoachV1(ctx, SetAthleteCoachV1Params{ID: archived.ID, IsCoach: true})
	require.NoError(t, err)
	require.True(t, updated.IsCoach)
	require.False(t, updated.IsActive, "coach classification must work for archived athletes")
	require.Equal(t, archived.ID, updated.ID)
	require.Equal(t, archived.FirstName, updated.FirstName)
	require.Equal(t, archived.SecondName, updated.SecondName)
	require.Equal(t, archived.BirthDate, updated.BirthDate)
	require.Equal(t, archived.CreatedAt, updated.CreatedAt)
	updatedAt := updated.UpdatedAt

	unchanged, err := txQueries.SetAthleteCoachV1(ctx, SetAthleteCoachV1Params{ID: archived.ID, IsCoach: true})
	require.NoError(t, err)
	require.Equal(t, updatedAt, unchanged.UpdatedAt, "repeating desired state must not bump updated_at")

	coachFilter := GetAthletesV1Params{IsActive: false, SearchQuery: "коучархив" + suffix, IsCoach: sql.NullBool{Bool: true, Valid: true}, LimitCount: 50}
	coaches, err := txQueries.GetAthletesV1(ctx, coachFilter)
	require.NoError(t, err)
	require.Equal(t, []int64{archived.ID}, athleteIDs(coaches), "coach filter combines with search and archive status before pagination")

	athleteFilter := GetAthletesV1Params{IsActive: true, SearchQuery: "коучфильтр" + suffix, IsCoach: sql.NullBool{Bool: false, Valid: true}, LimitCount: 50}
	athletes, err := txQueries.GetAthletesV1(ctx, athleteFilter)
	require.NoError(t, err)
	require.Equal(t, []int64{active.ID}, athleteIDs(athletes))
	count, err := txQueries.CountAthletesV1(ctx, CountAthletesV1Params{IsActive: false, SearchQuery: "коучархив" + suffix, IsCoach: sql.NullBool{Bool: true, Valid: true}})
	require.NoError(t, err)
	require.EqualValues(t, 1, count, "count applies the same coach, status and search filters")

	_, err = txQueries.SetAthleteCoachV1(ctx, SetAthleteCoachV1Params{ID: archived.ID, IsCoach: false})
	require.NoError(t, err)
	demoted, err := txQueries.GetAthleteByIDV1(ctx, archived.ID)
	require.NoError(t, err)
	require.False(t, demoted.IsCoach, "the endpoint supports demotion as a desired state")
	require.False(t, demoted.IsActive)
	var usersAfter int64
	require.NoError(t, tx.QueryRowContext(ctx, "SELECT count(*) FROM root.users").Scan(&usersAfter))
	require.Equal(t, usersBefore, usersAfter, "coach classification must not create or remove user accounts")
}
