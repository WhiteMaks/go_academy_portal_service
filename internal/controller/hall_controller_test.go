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

func TestHallController_CreateHallV1_400_capacity(t *testing.T) {
	hallMock := &mockHallService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	body := []byte(`{"branch_id":1,"name":"Зал","capacity":0}`)
	testContext.Request = httptest.NewRequest(http.MethodPost, route.ApiHallV1, bytes.NewReader(body))
	testContext.Request.Header.Set("Content-Type", "application/json")
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewHallController(hallMock, userMock)
	controller.CreateHallV1(testContext)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHallController_CreateHallV1_201(t *testing.T) {
	expected := model.HallV1Response{
		ID:             3,
		BranchID:       1,
		Name:           "Основной",
		Capacity:       40,
		ActivityStatus: "active",
	}

	hallMock := &mockHallService{
		createHallV1Func: func(ctx context.Context, request model.PostHallV1Request) (model.HallV1Response, error) {
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
	testContext.Request = tPreparePostRequest(route.ApiHallV1, model.PostHallV1Request{
		BranchID: 1,
		Name:     "Основной",
		Capacity: 40,
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewHallController(hallMock, userMock)
	controller.CreateHallV1(testContext)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestHallController_UpdateHallV1_409(t *testing.T) {
	hallMock := &mockHallService{
		updateHallV1Func: func(ctx context.Context, id int64, request model.PutHallV1Request) (model.HallV1Response, error) {
			return model.HallV1Response{}, service.ErrHallNameConflict
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePostRequest(route.ApiHallByIDV1, model.PutHallV1Request{
		Name:           "Дубль",
		Capacity:       20,
		ActivityStatus: "active",
	})
	testContext.Request.Method = http.MethodPut
	testContext.Params = gin.Params{{Key: "id", Value: "1"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewHallController(hallMock, userMock)
	controller.UpdateHallV1(testContext)

	require.Equal(t, http.StatusConflict, recorder.Code)
}
