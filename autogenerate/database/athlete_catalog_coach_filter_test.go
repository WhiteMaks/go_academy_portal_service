package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/stretchr/testify/require"
)

func TestAthleteCatalogCoachFilterSQL_DisposableDatabase(t *testing.T) {
	db, _ := openAthleteFeatureTestDB(t)
	_, queries := beginAthleteFeatureTestTx(t, db)
	ctx := context.Background()
	suffix := util.RandomString(12)
	activeAthlete := createFeatureTestAthlete(t, queries, "AthleteCatalog"+suffix, "Анна"+suffix, time.Date(2013, 2, 1, 0, 0, 0, 0, time.UTC))
	activeCoach := createFeatureTestAthlete(t, queries, "CoachCatalog"+suffix, "Борис"+suffix, time.Date(2012, 3, 2, 0, 0, 0, 0, time.UTC))
	archivedAthlete := createFeatureTestAthlete(t, queries, "AthleteArchived"+suffix, "Вера"+suffix, time.Date(2011, 4, 3, 0, 0, 0, 0, time.UTC))
	archivedCoach := createFeatureTestAthlete(t, queries, "CoachArchived"+suffix, "Глеб"+suffix, time.Date(2010, 5, 4, 0, 0, 0, 0, time.UTC))
	_, err := queries.SetAthleteCoachV1(ctx, SetAthleteCoachV1Params{ID: activeCoach.ID, IsCoach: true})
	require.NoError(t, err)
	_, err = queries.ArchiveAthleteV1(ctx, archivedAthlete.ID)
	require.NoError(t, err)
	_, err = queries.SetAthleteCoachV1(ctx, SetAthleteCoachV1Params{ID: archivedCoach.ID, IsCoach: true})
	require.NoError(t, err)
	_, err = queries.ArchiveAthleteV1(ctx, archivedCoach.ID)
	require.NoError(t, err)

	for _, tc := range []struct {
		name      string
		isActive  bool
		isCoach   sql.NullBool
		includeID int64
		excludeID int64
	}{
		{name: "active coaches", isActive: true, isCoach: sql.NullBool{Bool: true, Valid: true}, includeID: activeCoach.ID, excludeID: activeAthlete.ID},
		{name: "active athletes", isActive: true, isCoach: sql.NullBool{Bool: false, Valid: true}, includeID: activeAthlete.ID, excludeID: activeCoach.ID},
		{name: "archived coaches", isActive: false, isCoach: sql.NullBool{Bool: true, Valid: true}, includeID: archivedCoach.ID, excludeID: archivedAthlete.ID},
		{name: "archived athletes", isActive: false, isCoach: sql.NullBool{Bool: false, Valid: true}, includeID: archivedAthlete.ID, excludeID: archivedCoach.ID},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := GetAthletesV1Params{IsActive: tc.isActive, IsCoach: tc.isCoach, LimitCount: 100}
			filtered, err := queries.GetAthletesV1(ctx, params)
			require.NoError(t, err)
			ids := athleteIDs(filtered)
			require.Contains(t, ids, tc.includeID)
			require.NotContains(t, ids, tc.excludeID)

			total, err := queries.CountAthletesV1(ctx, CountAthletesV1Params{IsActive: tc.isActive, IsCoach: tc.isCoach})
			require.NoError(t, err)
			require.GreaterOrEqual(t, total, int64(len(filtered)))
			require.GreaterOrEqual(t, total, int64(1))

			index := indexOfAthlete(ids, tc.includeID)
			require.GreaterOrEqual(t, index, 0)
			page, err := queries.GetAthletesV1(ctx, GetAthletesV1Params{
				IsActive: tc.isActive, IsCoach: tc.isCoach, LimitCount: 1, OffsetCount: int32(index),
			})
			require.NoError(t, err)
			require.Equal(t, []int64{tc.includeID}, athleteIDs(page), "filtering must precede offset pagination")
		})
	}

	// The omitted filter remains unfiltered and returns both profile kinds.
	unfiltered, err := queries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, LimitCount: 100})
	require.NoError(t, err)
	require.Contains(t, athleteIDs(unfiltered), activeAthlete.ID)
	require.Contains(t, athleteIDs(unfiltered), activeCoach.ID)
}

func indexOfAthlete(ids []int64, expected int64) int {
	for index, id := range ids {
		if id == expected {
			return index
		}
	}
	return -1
}
