package controller

import (
	"context"
	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/route"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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
			}

			recorder := httptest.NewRecorder()
			testContext, _ := gin.CreateTestContext(recorder)
			request := tPreparePostRequest(route.ApiUserV1, test.errorRequestBody)

			testContext.Request = request

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

func TestUserController_CreateUserV1_409(t *testing.T) {
	expectedStatusCode := http.StatusConflict
	expectedResponseBody := model.ErrorResponse{
		Message: "pq: user with the same username or email already exists",
	}

	mock := &mockUserService{
		createUserV1Func: func(ctx context.Context, request model.PostUserV1Request) (model.PostUserV1Response, error) {
			return model.PostUserV1Response{}, &pq.Error{Code: database.PgErrUniqueViolation, Message: "user with the same username or email already exists"}
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

	controller := NewUserController(mock)
	controller.CreateUserV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.ErrorResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}
