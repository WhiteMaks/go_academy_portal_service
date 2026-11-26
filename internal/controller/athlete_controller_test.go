package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/middleware"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/route"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAthleteController_CreateAthleteV1_201(t *testing.T) {
	expected := model.AthleteV1Response{
		ID:             1,
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:          "+79001234567",
		ActivityStatus: "active",
	}

	athleteMock := &mockAthleteService{
		createAthleteV1Func: func(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
			return expected, nil
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePostRequest(route.ApiAthleteV1, model.PostAthleteV1Request{
		LastName:  "Ivanov",
		FirstName: "Ivan",
		BirthDate: model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:     "+79001234567",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.CreateAthleteV1(testContext)

	require.Equal(t, http.StatusCreated, recorder.Code)
	var actual model.AthleteV1Response
	tPrepareResponse(recorder.Body, &actual)
	require.Equal(t, expected.ID, actual.ID)
}

func TestAthleteController_CreateAthleteV1_403(t *testing.T) {
	athleteMock := &mockAthleteService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return false
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePostRequest(route.ApiAthleteV1, model.PostAthleteV1Request{
		LastName:  "Ivanov",
		FirstName: "Ivan",
		BirthDate: model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:     "+79001234567",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("coach", string(database.RootUserRoleCoach), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.CreateAthleteV1(testContext)

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestAthleteController_CreateAthleteV1_400Validation(t *testing.T) {
	athleteMock := &mockAthleteService{
		createAthleteV1Func: func(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
			return model.AthleteV1Response{}, service.ErrAthleteValidation
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePostRequest(route.ApiAthleteV1, model.PostAthleteV1Request{
		LastName:  "Ivanov",
		FirstName: "Ivan",
		BirthDate: model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:     "invalid",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.CreateAthleteV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAthleteController_GetAthletesV1_200(t *testing.T) {
	expected := model.GetAthletesV1Response{
		Athletes: []model.AthleteV1Response{},
		Total:    0,
		Limit:    50,
		Offset:   0,
	}

	athleteMock := &mockAthleteService{
		getAthletesV1Func: func(ctx context.Context, query model.GetAthletesV1Query) (model.GetAthletesV1Response, error) {
			return expected, nil
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPrepareGetRequest(route.ApiAthletesV1 + "?q=2017")
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.GetAthletesV1(testContext)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAthleteController_GetAthletesV1_400_invalidLimit(t *testing.T) {
	athleteMock := &mockAthleteService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPrepareGetRequest(route.ApiAthletesV1 + "?limit=999")
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.GetAthletesV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAthleteController_CreateAthleteV1_400_nameTooLong(t *testing.T) {
	athleteMock := &mockAthleteService{
		createAthleteV1Func: func(ctx context.Context, request model.PostAthleteV1Request) (model.AthleteV1Response, error) {
			t.Fatal("service must not be called when binding fails")
			return model.AthleteV1Response{}, nil
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePostRequest(route.ApiAthleteV1, model.PostAthleteV1Request{
		LastName:  strings.Repeat("а", model.AthleteNameMaxLen+1),
		FirstName: "Ivan",
		BirthDate: model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:     "+79001234567",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.CreateAthleteV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAthleteController_GetAthletesV1_400_invalidOffset(t *testing.T) {
	athleteMock := &mockAthleteService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPrepareGetRequest(route.ApiAthletesV1 + "?offset=-1")
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.GetAthletesV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestAthleteController_GetAthleteV1_200(t *testing.T) {
	expected := model.AthleteV1Response{ID: 7, LastName: "Ivanov", FirstName: "Ivan", ActivityStatus: "active"}

	athleteMock := &mockAthleteService{
		getAthleteV1Func: func(ctx context.Context, id int64) (model.AthleteV1Response, error) {
			require.Equal(t, int64(7), id)
			return expected, nil
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPrepareGetRequest(route.ApiAthletesV1 + "/7")
	testContext.Params = gin.Params{{Key: "id", Value: "7"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.GetAthleteV1(testContext)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestAthleteController_UpdateAthleteV1_200(t *testing.T) {
	expected := model.AthleteV1Response{
		ID:             7,
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:          "+79001234567",
		ActivityStatus: "inactive",
	}

	athleteMock := &mockAthleteService{
		updateAthleteV1Func: func(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
			require.Equal(t, int64(7), id)
			require.Equal(t, "inactive", request.ActivityStatus)
			return expected, nil
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePutRequest(route.ApiAthletesV1+"/7", model.PutAthleteV1Request{
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:          "+79001234567",
		ActivityStatus: "inactive",
	})
	testContext.Params = gin.Params{{Key: "id", Value: "7"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.UpdateAthleteV1(testContext)

	require.Equal(t, http.StatusOK, recorder.Code)
	var actual model.AthleteV1Response
	tPrepareResponse(recorder.Body, &actual)
	require.Equal(t, "inactive", actual.ActivityStatus)
}

func TestAthleteController_UpdateAthleteV1_404(t *testing.T) {
	athleteMock := &mockAthleteService{
		updateAthleteV1Func: func(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
			return model.AthleteV1Response{}, service.ErrAthleteNotFound
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePutRequest(route.ApiAthletesV1+"/99", model.PutAthleteV1Request{
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:          "+79001234567",
		ActivityStatus: "active",
	})
	testContext.Params = gin.Params{{Key: "id", Value: "99"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.UpdateAthleteV1(testContext)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestAthleteController_UpdateAthleteV1_403(t *testing.T) {
	athleteMock := &mockAthleteService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return false
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePutRequest(route.ApiAthletesV1+"/7", model.PutAthleteV1Request{
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:          "+79001234567",
		ActivityStatus: "active",
	})
	testContext.Params = gin.Params{{Key: "id", Value: "7"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("coach", string(database.RootUserRoleCoach), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.UpdateAthleteV1(testContext)

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestAthleteController_UpdateAthleteV1_409_phoneConflict(t *testing.T) {
	athleteMock := &mockAthleteService{
		updateAthleteV1Func: func(ctx context.Context, id int64, request model.PutAthleteV1Request) (model.AthleteV1Response, error) {
			return model.AthleteV1Response{}, service.ErrAthletePhoneConflict
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePutRequest(route.ApiAthletesV1+"/7", model.PutAthleteV1Request{
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:          "+79009998877",
		ActivityStatus: "active",
	})
	testContext.Params = gin.Params{{Key: "id", Value: "7"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.UpdateAthleteV1(testContext)

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestAthleteController_GetAthleteV1_400_invalidID(t *testing.T) {
	athleteMock := &mockAthleteService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPrepareGetRequest(route.ApiAthletesV1 + "/not-id")
	testContext.Params = gin.Params{{Key: "id", Value: "not-id"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewAthleteController(athleteMock, userMock)
	controller.GetAthleteV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
