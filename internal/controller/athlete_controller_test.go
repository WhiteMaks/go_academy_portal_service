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
	getFunc     func(context.Context, model.GetAthletesV1Query) (model.GetAthletesV1Response, error)
	getByIDFunc func(context.Context, int64) (model.AthleteV1Response, error)
	updateFunc  func(context.Context, int64, model.PutAthleteV1Request) (model.AthleteV1Response, error)
	archiveFunc func(context.Context, int64) (model.AthleteV1Response, error)
	restoreFunc func(context.Context, int64) (model.AthleteV1Response, error)
	coachFunc   func(context.Context, int64, bool) (model.AthleteV1Response, error)
}

func (m *mockAthleteService) CreateAthleteV1(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
	return m.createFunc(ctx, request)
}

func (m *mockAthleteService) GetAthletesV1(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
	return m.getFunc(ctx, query)
}

func (m *mockAthleteService) GetAthleteByIDV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockAthleteService) UpdateAthleteV1(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
	return m.updateFunc(ctx, id, request)
}

func (m *mockAthleteService) ArchiveAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	return m.archiveFunc(ctx, id)
}

func (m *mockAthleteService) RestoreAthleteV1(ctx context.Context, id int64) (model.AthleteV1Response, error) {
	return m.restoreFunc(ctx, id)
}

func (m *mockAthleteService) SetAthleteCoachV1(ctx context.Context, id int64, isCoach bool) (model.AthleteV1Response, error) {
	return m.coachFunc(ctx, id, isCoach)
}

func TestAthleteController_GetAthletesV1_RequiresAdmin(t *testing.T) {
	called := false
	athleteService := &mockAthleteService{getFunc: func(context.Context, model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
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
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/athlete/v1", nil)
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
	athleteService := &mockAthleteService{getFunc: func(context.Context, model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
		return model.GetAthletesV1Response{}, errors.New("database password leaked")
	}}
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", "admin", time.Minute))
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/athlete/v1", nil)

	NewAthleteController(athleteService, userService).GetAthletesV1(ctx)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.JSONEq(t, `{"message":"internal server error"}`, recorder.Body.String())
}

func TestParseAthleteListQueryDefaultsAndBounds(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  model.GetAthletesV1Query
	}{
		{name: "defaults", want: model.GetAthletesV1Query{IsActive: true, Limit: 50}},
		{name: "archived and bounds", query: "status=archived&limit=100&offset=2147483647", want: model.GetAthletesV1Query{IsActive: false, Limit: 100, Offset: 2147483647}},
		{name: "minimum limit", query: "limit=1", want: model.GetAthletesV1Query{IsActive: true, Limit: 1}},
		{name: "coach false", query: "isCoach=false", want: model.GetAthletesV1Query{IsActive: true, IsCoach: boolPointer(false), Limit: 50}},
		{name: "coach true with other filters", query: "status=archived&q=anna&isCoach=true", want: model.GetAthletesV1Query{IsActive: false, SearchQuery: "anna", IsCoach: boolPointer(true), Limit: 50}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseAthleteListQuery(tc.query)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseAthleteListQueryRejectsInvalidValues(t *testing.T) {
	for _, query := range []string{
		"status=deleted", "status=active&status=archived", "limit=0", "limit=101", "limit=1.5",
		"offset=-1", "offset=2147483648", "unknown=value", "q=%FF", "q=" + strings.Repeat("a", 201),
		"isCoach=1", "isCoach=TRUE", "isCoach=", "isCoach=true&isCoach=false",
	} {
		t.Run(query, func(t *testing.T) {
			_, err := parseAthleteListQuery(query)
			require.Error(t, err)
		})
	}
}

func TestAthleteController_GetAthletesV1ParsesQueryAndReturnsEnvelope(t *testing.T) {
	var gotQuery model.GetAthletesV1Query
	athleteService := &mockAthleteService{getFunc: func(_ context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
		gotQuery = query
		return model.GetAthletesV1Response{Athletes: []model.AthleteV1Response{}, Total: 42, Limit: query.Limit, Offset: query.Offset}, nil
	}}
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	recorder, ctx := athleteRequestContext(http.MethodGet, "", "")
	ctx.Request.URL.RawQuery = "status=archived&q=23.04&limit=25&offset=50"
	NewAthleteController(athleteService, userService).GetAthletesV1(ctx)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, model.GetAthletesV1Query{IsActive: false, SearchQuery: "2304", IsBirthDateQuery: true, Limit: 25, Offset: 50}, gotQuery)
	require.JSONEq(t, `{"athletes":[],"total":42,"limit":25,"offset":50}`, recorder.Body.String())
}

func TestAthleteController_GetAthletesV1RejectsInvalidQuery(t *testing.T) {
	for _, query := range []string{"status=unknown", "limit=0", "offset=-1", "extra=1", "status=active&status=archived"} {
		t.Run(query, func(t *testing.T) {
			called := false
			athleteService := &mockAthleteService{getFunc: func(context.Context, model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
				called = true
				return model.GetAthletesV1Response{}, nil
			}}
			userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
			recorder, ctx := athleteRequestContext(http.MethodGet, "", "")
			ctx.Request.URL.RawQuery = query
			NewAthleteController(athleteService, userService).GetAthletesV1(ctx)
			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.JSONEq(t, `{"message":"invalid request query"}`, recorder.Body.String())
			require.False(t, called)
		})
	}
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

func TestAthleteController_ArchiveRestoreResponsesAndErrors(t *testing.T) {
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	for _, tc := range []struct {
		name       string
		method     func(service.AthleteService, *gin.Context)
		result     model.AthleteV1Response
		serviceErr error
		status     int
		message    string
	}{
		{name: "archive success", method: func(s service.AthleteService, c *gin.Context) {
			NewAthleteController(s, userService).ArchiveAthleteV1(c)
		}, result: model.AthleteV1Response{ID: "7", IsActive: false}, status: http.StatusOK},
		{name: "restore success", method: func(s service.AthleteService, c *gin.Context) {
			NewAthleteController(s, userService).RestoreAthleteV1(c)
		}, result: model.AthleteV1Response{ID: "7", IsActive: true}, status: http.StatusOK},
		{name: "archive conflict", method: func(s service.AthleteService, c *gin.Context) {
			NewAthleteController(s, userService).ArchiveAthleteV1(c)
		}, serviceErr: service.ErrAthleteAlreadyArchived, status: http.StatusConflict, message: "athlete already archived"},
		{name: "restore conflict", method: func(s service.AthleteService, c *gin.Context) {
			NewAthleteController(s, userService).RestoreAthleteV1(c)
		}, serviceErr: service.ErrAthleteAlreadyActive, status: http.StatusConflict, message: "athlete is already active"},
		{name: "unknown id", method: func(s service.AthleteService, c *gin.Context) {
			NewAthleteController(s, userService).ArchiveAthleteV1(c)
		}, serviceErr: service.ErrAthleteNotFound, status: http.StatusNotFound, message: "athlete not found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			athleteService := &mockAthleteService{
				archiveFunc: func(_ context.Context, id int64) (model.AthleteV1Response, error) {
					require.Equal(t, int64(7), id)
					return tc.result, tc.serviceErr
				},
				restoreFunc: func(_ context.Context, id int64) (model.AthleteV1Response, error) {
					require.Equal(t, int64(7), id)
					return tc.result, tc.serviceErr
				},
			}
			recorder, ctx := athleteRequestContext(http.MethodPost, "", "7")
			tc.method(athleteService, ctx)
			require.Equal(t, tc.status, recorder.Code)
			if tc.message != "" {
				require.JSONEq(t, `{"message":"`+tc.message+`"}`, recorder.Body.String())
			}
		})
	}
}

func TestAthleteController_ArchiveRejectsInvalidIDAndBody(t *testing.T) {
	userService := &mockUserService{hasAccessFunc: func(context.Context, *util.Payload, []database.RootUserRole) bool { return true }}
	for _, tc := range []struct{ id, body, message string }{{id: "01", message: "invalid athlete id"}, {id: "7", body: "{}", message: "invalid request query"}} {
		called := false
		athleteService := &mockAthleteService{archiveFunc: func(context.Context, int64) (model.AthleteV1Response, error) {
			called = true
			return model.AthleteV1Response{}, nil
		}}
		recorder, ctx := athleteRequestContext(http.MethodPost, tc.body, tc.id)
		NewAthleteController(athleteService, userService).ArchiveAthleteV1(ctx)
		require.Equal(t, http.StatusBadRequest, recorder.Code)
		require.JSONEq(t, `{"message":"`+tc.message+`"}`, recorder.Body.String())
		require.False(t, called)
	}
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

func boolPointer(value bool) *bool { return &value }

func TestAthleteController_SetAthleteCoachV1StrictBodyAndDesiredState(t *testing.T) {
	for _, tc := range []struct {
		body   string
		status int
		want   bool
	}{
		{body: `{"isCoach":true}`, status: http.StatusOK, want: true},
		{body: `{"isCoach":false}`, status: http.StatusOK, want: false},
		{body: `{}`, status: http.StatusBadRequest},
		{body: `{"isCoach":null}`, status: http.StatusBadRequest},
		{body: `{"isCoach":"true"}`, status: http.StatusBadRequest},
		{body: `{"isCoach":true,"extra":1}`, status: http.StatusBadRequest},
		{body: `{"isCoach":true} {}`, status: http.StatusBadRequest},
	} {
		t.Run(tc.body, func(t *testing.T) {
			called := false
			athleteService := &mockAthleteService{coachFunc: func(_ context.Context, id int64, isCoach bool) (model.AthleteV1Response, error) {
				called = true
				require.EqualValues(t, 23, id)
				require.Equal(t, tc.want, isCoach)
				return model.AthleteV1Response{ID: "23", IsCoach: isCoach}, nil
			}}
			userService := &mockUserService{hasAccessFunc: func(_ context.Context, _ *util.Payload, roles []database.RootUserRole) bool {
				require.Equal(t, []database.RootUserRole{database.RootUserRoleAdmin}, roles)
				return true
			}}
			recorder, ctx := athleteRequestContext(http.MethodPut, tc.body, "23")
			NewAthleteController(athleteService, userService).SetAthleteCoachV1(ctx)
			require.Equal(t, tc.status, recorder.Code)
			require.Equal(t, tc.status == http.StatusOK, called)
		})
	}
}

var _ service.AthleteService = (*mockAthleteService)(nil)
