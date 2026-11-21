package database

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// This test requires an explicitly configured, disposable PostgreSQL database.
// It intentionally leaves the sole admin record behind, so the database can be
// discarded after the test instead of deleting data from a shared environment.
func TestCreateAdminV1_ConcurrentBootstrap(t *testing.T) {
	connectionString := os.Getenv("ADMIN_BOOTSTRAP_TEST_DATABASE_URL")
	if connectionString == "" {
		t.Skip("set ADMIN_BOOTSTRAP_TEST_DATABASE_URL and ADMIN_BOOTSTRAP_TEST_DISPOSABLE=I_UNDERSTAND_THIS_DATABASE_WILL_RETAIN_ONE_ADMIN to an empty disposable PostgreSQL database")
	}
	require.Equal(t, "I_UNDERSTAND_THIS_DATABASE_WILL_RETAIN_ONE_ADMIN", os.Getenv("ADMIN_BOOTSTRAP_TEST_DISPOSABLE"), "explicit confirmation is required because the test leaves one admin row")

	db, err := sql.Open("postgres", connectionString)
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, db.PingContext(context.Background()))

	var databaseName string
	require.NoError(t, db.QueryRowContext(context.Background(), "SELECT current_database()").Scan(&databaseName))
	require.Equal(t, "academy_portal_service_disposable", databaseName, "refusing to run outside the explicitly named disposable database")
	var userCount int
	require.NoError(t, db.QueryRowContext(context.Background(), "SELECT count(*) FROM root.users").Scan(&userCount))
	require.Zero(t, userCount, "test database must contain no users because the test intentionally leaves one admin")

	db.SetMaxOpenConns(2)

	params := CreateAdminV1Params{
		PUsername: util.RandomString(20),
		PPassword: util.RandomString(24),
	}
	start := make(chan struct{})
	ready := make(chan struct{}, 2)
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for i := 0; i < 2; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			conn, connErr := db.Conn(context.Background())
			if connErr != nil {
				ready <- struct{}{}
				results <- connErr
				return
			}
			defer conn.Close()
			queries := New(conn)
			ready <- struct{}{}
			<-start
			_, callErr := queries.CreateAdminV1(context.Background(), params)
			results <- callErr
		}()
	}
	<-ready
	<-ready
	close(start)
	workers.Wait()
	close(results)

	createdCount := 0
	adminExistsCount := 0
	for callErr := range results {
		if callErr == nil {
			createdCount++
			continue
		}
		var pqErr *pq.Error
		if errors.As(callErr, &pqErr) && pqErr.Code == PgErrNoDataFound {
			adminExistsCount++
			continue
		}
		require.NoError(t, callErr)
	}

	require.Equal(t, 1, createdCount)
	require.Equal(t, 1, adminExistsCount)
	var adminCount int
	err = db.QueryRowContext(context.Background(), "SELECT count(*) FROM root.users WHERE role = 'admin'::root.user_role").Scan(&adminCount)
	require.NoError(t, err)
	require.Equal(t, 1, adminCount)
}

func TestCreateUserV1(t *testing.T) {
	roles := []RootUserRole{
		RootUserRoleAdmin,
		RootUserRoleCoach,
		RootUserRoleAthlete,
	}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			createUserParams := PrepareRandomUserV1Params(role, true)

			actualRecord, err := testQueries.CreateUserV1(context.Background(), createUserParams)
			require.NoError(t, err)
			require.NotEmpty(t, actualRecord)

			expectedRecord := PrepareUserRecord(t, actualRecord, createUserParams)
			require.Equal(t, expectedRecord, actualRecord)
		})
	}
}

func TestCreateUserV1_UserExist(t *testing.T) {
	createUserParams := PrepareRandomUserV1Params(RootUserRoleAdmin, true)

	actualRecord, err := testQueries.CreateUserV1(context.Background(), createUserParams)
	require.NoError(t, err)
	require.NotEmpty(t, actualRecord)

	actualRecord, err = testQueries.CreateUserV1(context.Background(), createUserParams)
	require.Error(t, err)
	require.Empty(t, actualRecord)

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)

	require.True(t, ok)
	require.Equal(t, PgErrUniqueViolation, pqErr.Code)
	require.Equal(t, "user with the same username already exists", pqErr.Message)
}

func TestGetUserByUsernameV1(t *testing.T) {
	createUserParams := PrepareRandomUserV1Params(RootUserRoleAdmin, true)

	expectedRecord, _ := testQueries.CreateUserV1(context.Background(), createUserParams)

	actualRecord, err := testQueries.GetUserByUsernameV1(context.Background(), expectedRecord.Username)

	require.NoError(t, err)
	require.Equal(t, expectedRecord, actualRecord)
}

func TestGetUserByUsernameV1_UserNotFound(t *testing.T) {
	actualRecord, err := testQueries.GetUserByUsernameV1(context.Background(), util.RandomString(30))

	require.Empty(t, actualRecord)

	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)

	require.True(t, ok)
	require.Equal(t, PgErrNoDataFound, pqErr.Code)
	require.Equal(t, "user with given username not found", pqErr.Message)
}
