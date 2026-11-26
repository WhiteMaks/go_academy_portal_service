package controller

import (
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

func TestBranchController_CreateBranchV1_201(t *testing.T) {
	expected := model.BranchV1Response{
		ID:             1,
		Name:           "Центральный",
		Address:        "Саранск",
		ActivityStatus: "active",
	}

	branchMock := &mockBranchService{
		createBranchV1Func: func(ctx context.Context, request model.PostBranchV1Request) (model.BranchV1Response, error) {
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
	testContext.Request = tPreparePostRequest(route.ApiBranchV1, model.PostBranchV1Request{
		Name:    "Центральный",
		Address: "Саранск",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewBranchController(branchMock, userMock)
	controller.CreateBranchV1(testContext)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestBranchController_CreateBranchV1_403(t *testing.T) {
	branchMock := &mockBranchService{}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return false
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = tPreparePostRequest(route.ApiBranchV1, model.PostBranchV1Request{
		Name:    "Центральный",
		Address: "Саранск",
	})
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("coach", string(database.RootUserRoleCoach), time.Minute))

	controller := NewBranchController(branchMock, userMock)
	controller.CreateBranchV1(testContext)

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestBranchController_GetBranchV1_404(t *testing.T) {
	branchMock := &mockBranchService{
		getBranchV1Func: func(ctx context.Context, id int64) (model.BranchV1Response, error) {
			return model.BranchV1Response{}, service.ErrBranchNotFound
		},
	}
	userMock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = httptest.NewRequest(http.MethodGet, route.ApiBranchByIDV1, nil)
	testContext.Params = gin.Params{{Key: "id", Value: "404"}}
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Minute))

	controller := NewBranchController(branchMock, userMock)
	controller.GetBranchV1(testContext)

	require.Equal(t, http.StatusNotFound, recorder.Code)
}
