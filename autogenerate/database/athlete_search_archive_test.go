package database

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const (
	athleteFeatureTestURLEnv       = "ATHLETE_SEARCH_ARCHIVE_TEST_DATABASE_URL"
	athleteFeatureTestConfirmEnv   = "ATHLETE_SEARCH_ARCHIVE_TEST_DISPOSABLE"
	athleteFeatureTestConfirmation = "I_UNDERSTAND_THIS_DATABASE_IS_DISPOSABLE"
	athleteFeatureTestDatabaseName = "academy_portal_service_disposable"
)

func openAthleteFeatureTestDB(t *testing.T) (*sql.DB, *Queries) {
	t.Helper()
	connectionString := os.Getenv(athleteFeatureTestURLEnv)
	if connectionString == "" {
		t.Skipf("set %s and %s=%s to run against a disposable database", athleteFeatureTestURLEnv, athleteFeatureTestConfirmEnv, athleteFeatureTestConfirmation)
	}
	require.Equal(t, athleteFeatureTestConfirmation, os.Getenv(athleteFeatureTestConfirmEnv), "explicit confirmation is required for a disposable database")

	db, err := sql.Open("postgres", connectionString)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.PingContext(context.Background()))
	var databaseName string
	require.NoError(t, db.QueryRowContext(context.Background(), "SELECT current_database()").Scan(&databaseName))
	require.Equal(t, athleteFeatureTestDatabaseName, databaseName, "refusing to run outside the explicitly named disposable database")
	return db, New(db)
}

func beginAthleteFeatureTestTx(t *testing.T, db *sql.DB) (*sql.Tx, *Queries) {
	t.Helper()
	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = tx.Rollback() })
	return tx, New(tx)
}

func createFeatureTestAthlete(t *testing.T, queries *Queries, secondName, firstName string, birthDate time.Time) RootAthlete {
	t.Helper()
	record, err := queries.CreateAthleteV1(context.Background(), CreateAthleteV1Params{
		SecondName: secondName,
		FirstName:  firstName,
		BirthDate:  birthDate,
	})
	require.NoError(t, err)
	return record
}

func TestAthleteSearchArchiveSQL_DisposableDatabase(t *testing.T) {
	db, queries := openAthleteFeatureTestDB(t)
	tx, txQueries := beginAthleteFeatureTestTx(t, db)
	ctx := context.Background()
	suffix := util.RandomString(12)
	activeSurname := "Смирнова" + suffix
	activeName := "Анна" + suffix
	activeOne := createFeatureTestAthlete(t, txQueries, activeSurname, activeName, time.Date(2012, 4, 23, 0, 0, 0, 0, time.UTC))
	activeTwo := createFeatureTestAthlete(t, txQueries, "Петров"+suffix, activeName, time.Date(2014, 6, 3, 0, 0, 0, 0, time.UTC))
	lowerSurnameMatch, err := txQueries.CreateAthleteV1(ctx, CreateAthleteV1Params{
		SecondName: "Алексеев" + suffix,
		FirstName:  "Иван" + suffix,
		MiddleName: sql.NullString{String: "Иванович" + suffix, Valid: true},
		BirthDate:  time.Date(2010, 2, 1, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	higherSurnameMatch, err := txQueries.CreateAthleteV1(ctx, CreateAthleteV1Params{
		SecondName: "Иванов" + suffix,
		FirstName:  "Иван" + suffix,
		MiddleName: sql.NullString{String: "Алексеевич" + suffix, Valid: true},
		BirthDate:  time.Date(2011, 3, 2, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	archived := createFeatureTestAthlete(t, txQueries, "Иванов"+suffix, "Иван"+suffix, time.Date(2012, 4, 23, 0, 0, 0, 0, time.UTC))
	yoName := createFeatureTestAthlete(t, txQueries, "Ёлкин"+suffix, "Иван"+suffix, time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC))
	literalWildcardName := createFeatureTestAthlete(t, txQueries, "Процент%ный"+suffix, "Виктор"+suffix, time.Date(2016, 2, 2, 0, 0, 0, 0, time.UTC))
	wildcardDecoy := createFeatureTestAthlete(t, txQueries, "ПроцентХный"+suffix, "Виктор"+suffix, time.Date(2016, 2, 3, 0, 0, 0, 0, time.UTC))
	_, err = txQueries.ArchiveAthleteV1(ctx, archived.ID)
	require.NoError(t, err)
	var reducedPatterns []string
	err = tx.QueryRowContext(ctx, "SELECT root.athlete_search_patterns_v1($1::TEXT[])", pq.Array([]string{"иван", "иванов"})).Scan(pq.Array(&reducedPatterns))
	require.NoError(t, err)
	require.Equal(t, []string{"%иванов%"}, reducedPatterns)
	err = tx.QueryRowContext(ctx, "SELECT root.athlete_search_patterns_v1($1::TEXT[])", pq.Array([]string{"иванов", "иван"})).Scan(pq.Array(&reducedPatterns))
	require.NoError(t, err)
	require.Equal(t, []string{"%иванов%"}, reducedPatterns)
	err = tx.QueryRowContext(ctx, "SELECT root.athlete_search_patterns_v1($1::TEXT[])", pq.Array([]string{"иван", "иван"})).Scan(pq.Array(&reducedPatterns))
	require.NoError(t, err)
	require.Equal(t, []string{"%иван%"}, reducedPatterns)

	// Token order, case and ё/е normalization are verified at the SQL boundary.
	activeMatches, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, SearchQuery: strings.ToLower(activeSurname + " " + activeName), LimitCount: 50})
	require.NoError(t, err)
	require.Equal(t, []int64{activeOne.ID}, athleteIDs(activeMatches))
	// Tokens can match any FIO component and in any order, but surname hits rank first.
	firstOrderMatches, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, SearchQuery: "иванов иван", LimitCount: 50})
	require.NoError(t, err)
	require.Equal(t, []int64{higherSurnameMatch.ID, lowerSurnameMatch.ID}, athleteIDs(firstOrderMatches))
	reversedOrderMatches, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, SearchQuery: "иван иванов", LimitCount: 50})
	require.NoError(t, err)
	require.Equal(t, athleteIDs(firstOrderMatches), athleteIDs(reversedOrderMatches))
	tooManyTokens, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, SearchQuery: "иванов иван иван иванович", LimitCount: 50})
	require.NoError(t, err)
	require.Empty(t, tooManyTokens)
	tooManyTokenCount, err := txQueries.CountAthletesV1(ctx, CountAthletesV1Params{IsActive: true, SearchQuery: "иванов иван иван иванович"})
	require.NoError(t, err)
	require.Zero(t, tooManyTokenCount)
	yoMatches, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, SearchQuery: strings.ToLower("елкин" + suffix), LimitCount: 50})
	require.NoError(t, err)
	require.Equal(t, []int64{yoName.ID}, athleteIDs(yoMatches))
	wildcardMatches, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, SearchQuery: "процент%", LimitCount: 50})
	require.NoError(t, err)
	require.Contains(t, athleteIDs(wildcardMatches), literalWildcardName.ID)
	require.NotContains(t, athleteIDs(wildcardMatches), wildcardDecoy.ID)
	dateMatches, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, SearchQuery: "2304", IsBirthDateQuery: true, LimitCount: 50})
	require.NoError(t, err)
	require.Contains(t, athleteIDs(dateMatches), activeOne.ID)
	archivedMatches, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: false, SearchQuery: strings.ToLower("иван" + suffix), LimitCount: 50})
	require.NoError(t, err)
	require.Equal(t, []int64{archived.ID}, athleteIDs(archivedMatches))

	count, err := txQueries.CountAthletesV1(ctx, CountAthletesV1Params{IsActive: true, SearchQuery: strings.ToLower(activeName)})
	require.NoError(t, err)
	require.EqualValues(t, 2, count)
	allActive, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, LimitCount: 50})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(allActive), 2)
	paged, err := txQueries.GetAthletesV1(ctx, GetAthletesV1Params{IsActive: true, LimitCount: 1, OffsetCount: 1})
	require.NoError(t, err)
	require.Len(t, paged, 1)
	require.Equal(t, allActive[1].ID, paged[0].ID)

	// Wrong-state attempts return their dedicated SQLSTATE and do not alter status.
	state := athleteSQLStateInSavepoint(t, tx, func() error {
		_, err := txQueries.ArchiveAthleteV1(ctx, archived.ID)
		return err
	})
	require.Equal(t, pq.ErrorCode("P4091"), state)
	archived, err = txQueries.GetAthleteByIDV1(ctx, archived.ID)
	require.NoError(t, err)
	require.False(t, archived.IsActive)
	archived, err = txQueries.RestoreAthleteV1(ctx, archived.ID)
	require.NoError(t, err)
	require.True(t, archived.IsActive)
	state = athleteSQLStateInSavepoint(t, tx, func() error {
		_, err := txQueries.RestoreAthleteV1(ctx, activeTwo.ID)
		return err
	})
	require.Equal(t, pq.ErrorCode("P4092"), state)
	activeTwo, err = txQueries.GetAthleteByIDV1(ctx, activeTwo.ID)
	require.NoError(t, err)
	require.True(t, activeTwo.IsActive)

	// Unknown IDs remain distinct from wrong-state conflicts.
	state = athleteSQLStateInSavepoint(t, tx, func() error {
		_, err := txQueries.ArchiveAthleteV1(ctx, 9223372036854775000)
		return err
	})
	require.Equal(t, pq.ErrorCode("P0002"), state)
	state = athleteSQLStateInSavepoint(t, tx, func() error {
		_, err := txQueries.RestoreAthleteV1(ctx, 9223372036854775000)
		return err
	})
	require.Equal(t, pq.ErrorCode("P0002"), state)

	// Confirm all work is transaction-scoped; rollback prevents fixture leakage.
	require.NoError(t, tx.Rollback())
	_, err = queries.GetAthleteByIDV1(ctx, activeOne.ID)
	require.Error(t, err)
}

func athleteIDs(records []RootAthlete) []int64 {
	ids := make([]int64, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.ID)
	}
	return ids
}

func athleteSQLStateInSavepoint(t *testing.T, tx *sql.Tx, invoke func() error) pq.ErrorCode {
	t.Helper()
	_, err := tx.ExecContext(context.Background(), "SAVEPOINT athlete_feature_expected_error")
	require.NoError(t, err)
	callErr := invoke()
	_, rollbackErr := tx.ExecContext(context.Background(), "ROLLBACK TO SAVEPOINT athlete_feature_expected_error")
	require.NoError(t, rollbackErr)
	_, releaseErr := tx.ExecContext(context.Background(), "RELEASE SAVEPOINT athlete_feature_expected_error")
	require.NoError(t, releaseErr)
	err = callErr
	require.Error(t, err)
	var pqErr *pq.Error
	require.True(t, errors.As(err, &pqErr), "expected PostgreSQL error")
	return pqErr.Code
}
