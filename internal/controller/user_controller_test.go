package controller

import (
	"context"
	"errors"
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
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestUserController_CreateUserV1_201(t *testing.T) {
	expectedStatusCode := http.StatusCreated
	expectedResponseBody := model.PostUserV1Response{
		ID: 1,
	}

	mock := &mockUserService{
		createUserV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
			return expectedResponseBody, nil
		},
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	request := tPreparePostRequest(
		route.ApiUserV1,
		model.PostUserV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Second))

	controller := NewUserController(mock)
	controller.CreateUserV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.PostUserV1Response
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_CreateUserV1_400(t *testing.T) {
	tests := []struct {
		name                 string
		expectedErrorMessage string
		errorRequestBody     model.PostUserV1Request
	}{
		{
			name:                 "Without Username",
			expectedErrorMessage: "Key: 'PostUserV1Request.Username' Error:Field validation for 'Username' failed on the 'required' tag",
			errorRequestBody: model.PostUserV1Request{
				Password: util.RandomString(64),
			},
		},
		{
			name:                 "Without Password",
			expectedErrorMessage: "Key: 'PostUserV1Request.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			errorRequestBody: model.PostUserV1Request{
				Username: util.RandomString(64),
			},
		},
		{
			name:                 "Min username length",
			expectedErrorMessage: "Key: 'PostUserV1Request.Username' Error:Field validation for 'Username' failed on the 'min' tag",
			errorRequestBody: model.PostUserV1Request{
				Username: util.RandomString(4),
				Password: util.RandomString(64),
			},
		},
		{
			name:                 "Min password length",
			expectedErrorMessage: "Key: 'PostUserV1Request.Password' Error:Field validation for 'Password' failed on the 'min' tag",
			errorRequestBody: model.PostUserV1Request{
				Username: util.RandomString(64),
				Password: util.RandomString(7),
			},
		},
		{
			name:                 "Max username length",
			expectedErrorMessage: "Key: 'PostUserV1Request.Username' Error:Field validation for 'Username' failed on the 'max' tag",
			errorRequestBody: model.PostUserV1Request{
				Username: util.RandomString(65),
				Password: util.RandomString(64),
			},
		},
		{
			name:                 "Max password length",
			expectedErrorMessage: "Key: 'PostUserV1Request.Password' Error:Field validation for 'Password' failed on the 'max' tag",
			errorRequestBody: model.PostUserV1Request{
				Username: util.RandomString(64),
				Password: util.RandomString(65),
			},
		},
		{
			name:                 "Username with special characters",
			expectedErrorMessage: "Key: 'PostUserV1Request.Username' Error:Field validation for 'Username' failed on the 'alphanum' tag",
			errorRequestBody: model.PostUserV1Request{
				Username: "123$%^&*авыва",
				Password: util.RandomString(64),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedStatusCode := http.StatusBadRequest
			expectedResponseBody := model.ErrorResponse{
				Message: test.expectedErrorMessage,
			}

			mock := &mockUserService{
				createUserV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
					return model.PostUserV1Response{}, nil
				},
				hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
					return true
				},
			}

			recorder := httptest.NewRecorder()
			testContext, _ := gin.CreateTestContext(recorder)
			request := tPreparePostRequest(route.ApiUserV1, test.errorRequestBody)

			testContext.Request = request
			testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Second))

			controller := NewUserController(mock)
			controller.CreateUserV1(testContext)

			actualStatusCode := recorder.Code

			var actualResponseBody model.ErrorResponse
			tPrepareResponse(recorder.Body, &actualResponseBody)

			require.Equal(t, expectedStatusCode, actualStatusCode)
			require.Equal(t, expectedResponseBody, actualResponseBody)
		})
	}
}

func TestUserController_CreateUserV1_403(t *testing.T) {
	expectedStatusCode := http.StatusForbidden
	expectedResponseBody := service.PrepareForbiddenErrorResponse()

	mock := &mockUserService{
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return false
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	request := tPreparePostRequest(
		route.ApiUserV1,
		model.PostUserV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAthlete), time.Second))

	controller := NewUserController(mock)
	controller.CreateUserV1(testContext)

	actualStatusCode := recorder.Code
	var actualResponseBody model.ErrorResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_CreateUserV1_409(t *testing.T) {
	expectedStatusCode := http.StatusConflict
	expectedResponseBody := model.ErrorResponse{
		Message: "pq: user with the same username already exists",
	}

	mock := &mockUserService{
		createUserV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
			return model.PostUserV1Response{}, &pq.Error{Code: database.PgErrUniqueViolation, Message: "user with the same username already exists"}
		},
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	request := tPreparePostRequest(
		route.ApiUserV1,
		model.PostUserV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Second))

	controller := NewUserController(mock)
	controller.CreateUserV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.ErrorResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_CreateUserV1_500(t *testing.T) {
	expectedStatusCode := http.StatusInternalServerError
	expectedResponseBody := model.ErrorResponse{
		Message: "pq: some db error",
	}

	mock := &mockUserService{
		createUserV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
			return model.PostUserV1Response{}, &pq.Error{Code: "99999", Message: "some db error"}
		},
		hasAccessFunc: func(ctx context.Context, payload *util.Payload, roles []database.RootUserRole) bool {
			return true
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	request := tPreparePostRequest(
		route.ApiUserV1,
		model.PostUserV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request
	testContext.Set(middleware.AuthorizationPayload, util.NewTokenPayload("admin", string(database.RootUserRoleAdmin), time.Second))

	controller := NewUserController(mock)
	controller.CreateUserV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.ErrorResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_CreateAdminV1_201(t *testing.T) {
	expectedStatusCode := http.StatusCreated
	expectedResponseBody := model.PostUserV1Response{
		ID: 1,
	}

	mock := &mockUserService{
		createAdminV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
			return expectedResponseBody, nil
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	request := tPreparePostRequest(
		route.ApiUserAdminV1,
		model.PostUserV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request

	controller := NewUserController(mock)
	controller.CreateAdminV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.PostUserV1Response
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_CreateAdminV1_404(t *testing.T) {
	expectedStatusCode := http.StatusNotFound
	expectedResponseBody := service.PrepareEmptyResponse()

	mock := &mockUserService{
		createAdminV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
			return model.PostUserV1Response{}, service.ErrAdminCreationIsNotAllowed
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	request := tPreparePostRequest(
		route.ApiUserAdminV1,
		model.PostUserV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request

	controller := NewUserController(mock)
	controller.CreateAdminV1(testContext)

	actualStatusCode := recorder.Code
	var actualResponseBody model.EmptyResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_CreateAdminV1_500(t *testing.T) {
	expectedStatusCode := http.StatusInternalServerError
	expectedResponseBody := service.PrepareErrorResponse(errors.New("service error"))

	mock := &mockUserService{
		createAdminV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
			return model.PostUserV1Response{}, errors.New("service error")
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	request := tPreparePostRequest(
		route.ApiUserAdminV1,
		model.PostUserV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request

	controller := NewUserController(mock)
	controller.CreateAdminV1(testContext)

	actualStatusCode := recorder.Code
	var actualResponseBody model.ErrorResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_GenerateUserTokenV1_200(t *testing.T) {
	expectedStatusCode := http.StatusOK
	expectedResponseBody := model.PostUserTokenV1Response{
		Token: "token",
	}

	mock := &mockUserService{
		generateUserTokenV1Func: func(ctx context.Context, request model.PostUserTokenV1Request) (model.PostUserTokenV1Response, error) {
			return expectedResponseBody, nil
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	request := tPreparePostRequest(
		route.ApiUserTokenV1,
		model.PostUserTokenV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request

	controller := NewUserController(mock)
	controller.GenerateUserTokenV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.PostUserTokenV1Response
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestUserController_GenerateUserTokenV1_400(t *testing.T) {
	tests := []struct {
		name                 string
		expectedErrorMessage string
		errorRequestBody     model.PostUserTokenV1Request
	}{
		{
			name:                 "Without Username",
			expectedErrorMessage: "Key: 'PostUserTokenV1Request.Username' Error:Field validation for 'Username' failed on the 'required' tag",
			errorRequestBody: model.PostUserTokenV1Request{
				Password: util.RandomString(64),
			},
		},
		{
			name:                 "Without Password",
			expectedErrorMessage: "Key: 'PostUserTokenV1Request.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			errorRequestBody: model.PostUserTokenV1Request{
				Username: util.RandomString(64),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedStatusCode := http.StatusBadRequest
			expectedResponseBody := model.ErrorResponse{
				Message: test.expectedErrorMessage,
			}

			mock := &mockUserService{
				generateUserTokenV1Func: func(ctx context.Context, request model.PostUserTokenV1Request) (model.PostUserTokenV1Response, error) {
					return model.PostUserTokenV1Response{}, nil
				},
			}

			recorder := httptest.NewRecorder()
			testContext, _ := gin.CreateTestContext(recorder)
			request := tPreparePostRequest(route.ApiUserTokenV1, test.errorRequestBody)

			testContext.Request = request

			controller := NewUserController(mock)
			controller.GenerateUserTokenV1(testContext)

			actualStatusCode := recorder.Code

			var actualResponseBody model.ErrorResponse
			tPrepareResponse(recorder.Body, &actualResponseBody)

			require.Equal(t, expectedStatusCode, actualStatusCode)
			require.Equal(t, expectedResponseBody, actualResponseBody)
		})
	}
}

func TestUserController_GenerateUserTokenV1_401(t *testing.T) {
	expectedStatusCode := http.StatusUnauthorized
	expectedResponseBody := model.ErrorResponse{
		Message: "invalid credentials",
	}

	mock := &mockUserService{
		generateUserTokenV1Func: func(ctx context.Context, request model.PostUserTokenV1Request) (model.PostUserTokenV1Response, error) {
			return model.PostUserTokenV1Response{}, &pq.Error{Code: database.PgErrNoDataFound, Message: "user with given username not found"}
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	request := tPreparePostRequest(
		route.ApiUserTokenV1,
		model.PostUserTokenV1Request{
			Username: util.RandomString(64),
			Password: util.RandomString(64),
		},
	)

	testContext.Request = request

	controller := NewUserController(mock)
	controller.GenerateUserTokenV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.ErrorResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}
