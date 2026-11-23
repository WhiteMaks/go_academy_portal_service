package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type mockAthleteService struct {
	createFunc func(context.Context, model.PostAthleteV1Request) (model.AthleteV1Response, error)
	getFunc    func(context.Context) (model.GetAthletesV1Response, error)
}

func (m *mockAthleteService) CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
	return m.createFunc(ctx, request)
}

func (m *mockAthleteService) GetAthletesV1(ctx context.Context) (model.GetAthletesV1Response, error) {
	return m.getFunc(ctx)
}

func TestAthleteController_GetAthletesV1_RequiresAdmin(t *testing.T) {
	called := false
	athleteService := &mockAthleteService{getFunc: func(context.Context) (model.GetAthletesV1Response, error) {
		called = true
		return model.GetAthletesV1Response{}, nil
	}}
	userService := &mockUserService{hasAccessFunc: func(_ context.Context, _ *util.Payload, roles []database.RootUserRole) bool {
		require.Equal(t, []database.RootUserRole{database.RootUserRoleAdmin}, roles)
		return false
	}}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set(middleware.AuthorizationPayload, util.NewTokenPayload("coach", "coach", time.Minute))

	NewAthleteController(athleteService, userService).GetAthletesV1(ctx)
	require.Equal(t, http.StatusForbidden, recorder.Code)
	require.JSONEq(t, `{"message":"forbidden"}`, recorder.Body.String())
	require.False(t, called)
}

func TestAthleteController_CreateAthleteV1_BadRequest(t *testing.T) {
	called := false
	athleteService := &mockAthleteService{createFunc: func(context.Context, model.PostAthleteV1Request) (model.AthleteV1Response, error) {
		called = true
		return model.AthleteV1Response{}, nil
	}}
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", "admin", time.Minute))
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/athlete/v1", strings.NewReader(`{"lastName":`))

	NewAthleteController(athleteService, userService).CreateAthleteV1(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.False(t, called)
}

func TestAthleteController_CreateAthleteV1_RejectsInvalidRequestFields(t *testing.T) {
	cases := map[string]string{
		"missing second name (old key submitted)": `{"firstName":"Иван","lastName":"Иванов","birthDate":"2010-01-02T00:00:00Z"}`,
		"missing first name":                      `{"secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z"}`,
		"empty first name":                        `{"firstName":"","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z"}`,
		"empty second name":                       `{"firstName":"Иван","secondName":"","birthDate":"2010-01-02T00:00:00Z"}`,
		"first name too long":                     `{"firstName":"` + strings.Repeat("Я", 101) + `","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z"}`,
		"second name too long":                    `{"firstName":"Иван","secondName":"` + strings.Repeat("Я", 101) + `","birthDate":"2010-01-02T00:00:00Z"}`,
		"middle name too long":                    `{"firstName":"Иван","secondName":"Иванов","middleName":"` + strings.Repeat("Я", 101) + `","birthDate":"2010-01-02T00:00:00Z"}`,
		"parent first name too long":              `{"firstName":"Иван","secondName":"Иванов","parentFirstName":"` + strings.Repeat("Я", 101) + `","birthDate":"2010-01-02T00:00:00Z"}`,
		"missing date":                            `{"firstName":"Иван","secondName":"Иванов"}`,
		"null date":                               `{"firstName":"Иван","secondName":"Иванов","birthDate":null}`,
		"number date":                             `{"firstName":"Иван","secondName":"Иванов","birthDate":20100102}`,
		"date-only string":                        `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02"}`,
		"impossible date":                         `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-02-30T00:00:00Z"}`,
		"non UTC midnight date":                   `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00+01:00"}`,
		"non midnight date":                       `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T01:00:00Z"}`,
		"fractional second date":                  `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00.000000001Z"}`,
		"invalid photo URL":                       `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","photoUrl":"ftp://example.com/image.png"}`,
		"photo URL too long":                      `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","photoUrl":"https://example.com/` + strings.Repeat("a", 2048) + `"}`,
		"formatted phone":                         `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"+7 9991234567"}`,
		"phone without plus":                      `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","parentPhone":"79991234567"}`,
		"phone with zero country prefix":          `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"+09991234567"}`,
		"phone too short":                         `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"+123456789"}`,
		"phone too long":                          `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"+1234567890123456"}`,
		"phone contains letters":                  `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"+7999abc4567"}`,
		"parent phone malformed":                  `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","parentPhone":"   "}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			called := false
			athleteService := &mockAthleteService{createFunc: func(context.Context, model.PostAthleteV1Request) (model.AthleteV1Response, error) {
				called = true
				return model.AthleteV1Response{}, nil
			}}
			userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", "admin", time.Minute))
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/athlete/v1", strings.NewReader(body))

			NewAthleteController(athleteService, userService).CreateAthleteV1(ctx)
			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.False(t, called)
		})
	}
}

func TestAthleteController_CreateAthleteV1_AcceptsValidOptionalPhonesAndEmptyValues(t *testing.T) {
	for _, body := range []string{
		`{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z"}`,
		`{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":null,"parentPhone":null}`,
		`{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"","parentPhone":""}`,
		`{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"+1234567890"}`,
		`{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","phone":"+123456789012345"}`,
		`{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","parentPhone":"+1234567890"}`,
		`{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z","parentPhone":"+123456789012345"}`,
	} {
		t.Run(body, func(t *testing.T) {
			called := false
			athleteService := &mockAthleteService{createFunc: func(context.Context, model.PostAthleteV1Request) (model.AthleteV1Response, error) {
				called = true
				return model.AthleteV1Response{}, nil
			}}
			userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", "admin", time.Minute))
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/athlete/v1", strings.NewReader(body))

			NewAthleteController(athleteService, userService).CreateAthleteV1(ctx)
			require.Equal(t, http.StatusCreated, recorder.Code)
			require.True(t, called)
		})
	}
}

func TestAthleteController_GetAthletesV1_HidesInternalError(t *testing.T) {
	athleteService := &mockAthleteService{getFunc: func(context.Context) (model.GetAthletesV1Response, error) {
		return model.GetAthletesV1Response{}, errors.New("database password leaked")
	}}
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", "admin", time.Minute))

	NewAthleteController(athleteService, userService).GetAthletesV1(ctx)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.JSONEq(t, `{"message":"internal server error"}`, recorder.Body.String())
}

var _ service.AthleteService = (*mockAthleteService)(nil)
