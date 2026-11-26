package controller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
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

var coachCtrlBirthDate = model.DateOnly{Time: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)}

func TestCoachController_CreateCoachV1_201(t *testing.T) {
	expected := model.CoachV1Response{
		ID:             1,
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachCtrlBirthDate,
		Phone:          "+79001234567",
		Specialization: "Judo",
		ActivityStatus: "active",
	}

	coachMock := &mockCoachService{
		createCoachV1Func: func(ctx context.Context, request model.PostCoachV1Request) (model.CoachV1Response, error) {
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
	testContext.Request = tPreparePostRequest(route.ApiCoachV1, model.PostCoachV1Request{
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachCtrlBirthDate,
		Phone:          "+79001234567",
		Specialization: "Judo",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewCoachController(coachMock, userMock)
	controller.CreateCoachV1(testContext)

	require.Equal(t, http.StatusCreated, recorder.Code)
	var actual model.CoachV1Response
	tPrepareResponse(recorder.Body, &actual)
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, coachCtrlBirthDate, actual.BirthDate)
}

func TestCoachController_CreateCoachV1_400_missingBirthDate(t *testing.T) {
	coachMock := &mockCoachService{
		createCoachV1Func: func(ctx context.Context, request model.PostCoachV1Request) (model.CoachV1Response, error) {
			require.True(t, request.BirthDate.Time.IsZero())
			return model.CoachV1Response{}, service.ErrCoachValidation
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"last_name":"Sidorov","first_name":"Sergey","phone":"+79001234567","specialization":"Judo"}`)
	testContext.Request = httptest.NewRequest(http.MethodPost, route.ApiCoachV1, bytes.NewReader(body))
	testContext.Request.Header.Set("Content-Type", "application/json")
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewCoachController(coachMock, userMock)
	controller.CreateCoachV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCoachController_CreateCoachV1_403(t *testing.T) {
	coachMock := &mockCoachService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return false
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePostRequest(route.ApiCoachV1, model.PostCoachV1Request{
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachCtrlBirthDate,
		Phone:          "+79001234567",
		Specialization: "Judo",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("coach", string(database.RootUserRoleCoach), time.Minute))

	controller := NewCoachController(coachMock, userMock)
	controller.CreateCoachV1(testContext)

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestCoachController_UpdateCoachV1_404(t *testing.T) {
	coachMock := &mockCoachService{
		updateCoachV1Func: func(ctx context.Context, id int64, request model.PutCoachV1Request) (model.CoachV1Response, error) {
			return model.CoachV1Response{}, service.ErrCoachNotFound
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePutRequest(route.ApiCoachesV1+"/99", model.PutCoachV1Request{
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachCtrlBirthDate,
		Phone:          "+79001234567",
		Specialization: "Judo",
		ActivityStatus: "active",
	})
	testContext.Params = gin.Params{{Key: "id", Value: "99"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewCoachController(coachMock, userMock)
	controller.UpdateCoachV1(testContext)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestCoachController_UpdateCoachV1_409_phoneConflict(t *testing.T) {
	coachMock := &mockCoachService{
		updateCoachV1Func: func(ctx context.Context, id int64, request model.PutCoachV1Request) (model.CoachV1Response, error) {
			return model.CoachV1Response{}, service.ErrCoachPhoneConflict
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePutRequest(route.ApiCoachesV1+"/7", model.PutCoachV1Request{
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachCtrlBirthDate,
		Phone:          "+79009998877",
		Specialization: "Judo",
		ActivityStatus: "active",
	})
	testContext.Params = gin.Params{{Key: "id", Value: "7"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewCoachController(coachMock, userMock)
	controller.UpdateCoachV1(testContext)

	require.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCoachController_GetCoachV1_400_invalidID(t *testing.T) {
	coachMock := &mockCoachService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPrepareGetRequest(route.ApiCoachesV1 + "/not-id")
	testContext.Params = gin.Params{{Key: "id", Value: "not-id"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewCoachController(coachMock, userMock)
	controller.GetCoachV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
