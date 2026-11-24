package controller

import (
	"context"
	"encoding/json"
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
	createFunc  func(context.Context, model.PostAthleteV1Request) (model.AthleteV1Response, error)
	getFunc     func(context.Context) (model.GetAthletesV1Response, error)
	getByIDFunc func(context.Context, int64) (model.AthleteV1Response, error)
	updateFunc  func(context.Context, int64, model.PutAthleteV1Request) (model.AthleteV1Response, error)
}

func (m *mockAthleteService) CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
	return m.createFunc(ctx, request)
}

func (m *mockAthleteService) GetAthletesV1(ctx context.Context) (model.GetAthletesV1Response, error) {
	return m.getFunc(ctx)
}

func (m *mockAthleteService) GetAthleteByIDV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockAthleteService) UpdateAthleteV1(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
	return m.updateFunc(ctx, id, request)
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

func TestAthleteController_GetAthleteByIDV1_ValidatesIDAndMapsNotFound(t *testing.T) {
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	for _, id := range []string{"0", "01", "+1", "9223372036854775808", "abc"} {
		t.Run("invalid_"+id, func(t *testing.T) {
			called := false
			athleteService := &mockAthleteService{getByIDFunc: func(context.Context, int64) (model.AthleteV1Response, error) {
				called = true
				return model.AthleteV1Response{}, nil
			}}
			recorder, ctx := athleteRequestContext(http.MethodGet, "", id)
			NewAthleteController(athleteService, userService).GetAthleteByIDV1(ctx)
			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.False(t, called)
		})
	}
	athleteService := &mockAthleteService{getByIDFunc: func(_ context.Context, id int64) (model.AthleteV1Response, error) {
		require.Equal(t, int64(7), id)
		return model.AthleteV1Response{}, service.ErrAthleteNotFound
	}}
	recorder, ctx := athleteRequestContext(http.MethodGet, "", "7")
	NewAthleteController(athleteService, userService).GetAthleteByIDV1(ctx)
	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.JSONEq(t, `{"message":"athlete not found"}`, recorder.Body.String())
}

func TestAthleteController_UpdateAthleteV1_RejectsUnknownFieldsAndUpdates(t *testing.T) {
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	validBody := `{"firstName":"Иван","secondName":"Иванов","birthDate":"2010-01-02T00:00:00Z"}`
	for _, body := range []string{strings.TrimSuffix(validBody, "}") + `,"photoUrl":"https://example.com/a.jpg"}`, strings.TrimSuffix(validBody, "}") + `,"isActive":false}`, validBody + ` {}`} {
		t.Run("strict_body", func(t *testing.T) {
			called := false
			athleteService := &mockAthleteService{updateFunc: func(context.Context, int64, model.PutAthleteV1Request) (model.AthleteV1Response, error) {
				called = true
				return model.AthleteV1Response{}, nil
			}}
			recorder, ctx := athleteRequestContext(http.MethodPut, body, "7")
			NewAthleteController(athleteService, userService).UpdateAthleteV1(ctx)
			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.False(t, called)
		})
	}
	athleteService := &mockAthleteService{updateFunc: func(_ context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
		require.Equal(t, int64(7), id)
		require.Equal(t, "Иван", request.FirstName)
		return model.AthleteV1Response{ID: "7", FirstName: request.FirstName}, nil
	}}
	recorder, ctx := athleteRequestContext(http.MethodPut, validBody, "7")
	NewAthleteController(athleteService, userService).UpdateAthleteV1(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAthleteController_UpdateAthleteV1_ValidatesContractBoundaries(t *testing.T) {
	base := map[string]any{
		"firstName":  "Иван",
		"secondName": "Иванов",
		"birthDate":  "2010-01-02T00:00:00Z",
	}
	type testCase struct {
		name  string
		edit  func(map[string]any)
		valid bool
	}
	cases := []testCase{
		{name: "UTC midnight birthDate accepted", valid: true},
		{name: "birthDate offset", edit: func(body map[string]any) { body["birthDate"] = "2010-01-02T00:00:00+01:00" }},
		{name: "birthDate not midnight", edit: func(body map[string]any) { body["birthDate"] = "2010-01-02T00:00:01Z" }},
		{name: "date only birthDate", edit: func(body map[string]any) { body["birthDate"] = "2010-01-02" }},
		{name: "missing birthDate", edit: func(body map[string]any) { delete(body, "birthDate") }},
		{name: "required name whitespace only", edit: func(body map[string]any) { body["firstName"] = " \t " }},
		{name: "required name empty", edit: func(body map[string]any) { body["secondName"] = "" }},
		{name: "required name 100 runes accepted", edit: func(body map[string]any) { body["firstName"] = strings.Repeat("Я", 100) }, valid: true},
		{name: "required name 101 runes rejected", edit: func(body map[string]any) { body["firstName"] = strings.Repeat("Я", 101) }},
		{name: "optional name whitespace normalizes", edit: func(body map[string]any) { body["middleName"] = "   " }, valid: true},
		{name: "optional name 100 runes accepted", edit: func(body map[string]any) { body["parentFirstName"] = strings.Repeat("Я", 100) }, valid: true},
		{name: "optional name 101 runes rejected", edit: func(body map[string]any) { body["parentMiddleName"] = strings.Repeat("Я", 101) }},
		{name: "phone 11 characters accepted", edit: func(body map[string]any) { body["phone"] = "+1234567890" }, valid: true},
		{name: "phone 16 characters accepted", edit: func(body map[string]any) { body["parentPhone"] = "+123456789012345" }, valid: true},
		{name: "phone too short", edit: func(body map[string]any) { body["phone"] = "+123456789" }},
		{name: "phone too long", edit: func(body map[string]any) { body["phone"] = "+1234567890123456" }},
		{name: "phone whitespace only", edit: func(body map[string]any) { body["parentPhone"] = "   " }},
		{name: "phone formatted", edit: func(body map[string]any) { body["phone"] = "+1 234567890" }},
		{name: "phone invalid country prefix", edit: func(body map[string]any) { body["phone"] = "+0234567890" }},
	}
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			bodyValues := make(map[string]any, len(base)+2)
			for key, value := range base {
				bodyValues[key] = value
			}
			if test.edit != nil {
				test.edit(bodyValues)
			}
			encodedBody, err := json.Marshal(bodyValues)
			require.NoError(t, err)
			called := false
			athleteService := &mockAthleteService{updateFunc: func(_ context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
				called = true
				require.Equal(t, int64(7), id)
				return model.AthleteV1Response{}, nil
			}}
			recorder, ctx := athleteRequestContext(http.MethodPut, string(encodedBody), "7")
			NewAthleteController(athleteService, userService).UpdateAthleteV1(ctx)
			if test.valid {
				require.Equal(t, http.StatusOK, recorder.Code)
				require.True(t, called)
			} else {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.False(t, called)
			}
		})
	}
}

func athleteRequestContext(method, body, id string) (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = gin.Params{{Key: "id", Value: id}}
	ctx.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", "admin", time.Minute))
	ctx.Request = httptest.NewRequest(method, "/api/athlete/"+id+"/v1", strings.NewReader(body))
	return recorder, ctx
}

var _ service.AthleteService = (*mockAthleteService)(nil)
