package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/repository"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const disposableDatabaseMarker = "I_UNDERSTAND_THIS_IS_A_DISPOSABLE_TEST_DATABASE"

func TestAthleteV1PostgresIntegration(t *testing.T) {
	databaseURL := os.Getenv("MVP05_TEST_DATABASE_URL")
	if databaseURL == "" || os.Getenv("MVP05_TEST_DATABASE_DISPOSABLE") != disposableDatabaseMarker {
		t.Skip("requires explicit MVP05 disposable database URL and marker")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", databaseURL)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, db.PingContext(ctx))
	var databaseName string
	require.NoError(t, db.QueryRowContext(ctx, "SELECT current_database()").Scan(&databaseName))
	require.True(t, strings.HasSuffix(strings.ToLower(databaseName), "_test"), "disposable test database name must end in _test")

	store := database.NewStore(db)
	athleteService := service.NewAthleteService(repository.NewAthleteRepository(store))
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	controller := NewAthleteController(athleteService, userService)

	var id int64
	var before database.RootAthlete
	err = db.QueryRow(`INSERT INTO root.athletes (second_name, first_name, birth_date, photo_url)
		VALUES ('MVP05 integration', 'Transient athlete', DATE '2010-01-02', 'https://example.test/athlete.png')
		RETURNING id, second_name, first_name, middle_name, birth_date, photo_url, phone, parent_second_name,
		parent_first_name, parent_middle_name, parent_phone, is_active, created_at, updated_at`).Scan(
		&id, &before.SecondName, &before.FirstName, &before.MiddleName, &before.BirthDate, &before.PhotoUrl,
		&before.Phone, &before.ParentSecondName, &before.ParentFirstName, &before.ParentMiddleName,
		&before.ParentPhone, &before.IsActive, &before.CreatedAt, &before.UpdatedAt)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = db.Exec("DELETE FROM root.athletes WHERE id = $1", id) })

	request := model.PutAthleteV1Request{
		FirstName: "Updated athlete", SecondName: "MVP05 integration",
		BirthDate: time.Date(2010, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	body, err := json.Marshal(request)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(recorder)
	ginContext.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ginContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", "admin", time.Minute))
	ginContext.Request = httptest.NewRequest(http.MethodPut, "/api/athlete/"+strconv.FormatInt(id, 10)+"/v1", strings.NewReader(string(body)))
	controller.UpdateAthleteV1(ginContext)
	require.Equal(t, http.StatusOK, recorder.Code)

	var response model.AthleteV1Response
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	var after database.RootAthlete
	err = db.QueryRow(`SELECT id, second_name, first_name, middle_name, birth_date, photo_url, phone, parent_second_name,
		parent_first_name, parent_middle_name, parent_phone, is_active, created_at, updated_at
		FROM root.athletes WHERE id = $1`, id).Scan(
		&after.ID, &after.SecondName, &after.FirstName, &after.MiddleName, &after.BirthDate, &after.PhotoUrl,
		&after.Phone, &after.ParentSecondName, &after.ParentFirstName, &after.ParentMiddleName,
		&after.ParentPhone, &after.IsActive, &after.CreatedAt, &after.UpdatedAt)
	require.NoError(t, err)
	require.Equal(t, request.FirstName, after.FirstName)
	require.Equal(t, request.SecondName, after.SecondName)
	require.Equal(t, before.PhotoUrl, after.PhotoUrl)
	require.Equal(t, before.IsActive, after.IsActive)
	require.True(t, after.UpdatedAt.After(before.UpdatedAt))
	require.Equal(t, after.UpdatedAt, response.UpdatedAt)
}
