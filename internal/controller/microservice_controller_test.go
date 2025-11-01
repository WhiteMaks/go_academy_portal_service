package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/WhiteMaks/go_academy_portal_service/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMicroserviceController_GetMicroserviceV1_200(t *testing.T) {
	expectedStatusCode := http.StatusOK
	expectedResponseBody := model.GetMicroserviceV1Response{
		ReadyToUse: true,
	}

	mock := &mockMicroserviceService{
		getMicroserviceV1Func: func(ctx context.Context) (model.GetMicroserviceV1Response, error) {
			return expectedResponseBody, nil
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	request := tPrepareGetRequest(route.ApiMicroserviceV1)

	testContext.Request = request

	controller := NewMicroserviceController(mock)
	controller.GetMicroserviceV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.GetMicroserviceV1Response
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}

func TestMicroserviceController_GetMicroserviceV1_500(t *testing.T) {
	expectedStatusCode := http.StatusInternalServerError
	expectedResponseBody := model.ErrorResponse{
		Message: "internal Server Error",
	}

	mock := &mockMicroserviceService{
		getMicroserviceV1Func: func(ctx context.Context) (model.GetMicroserviceV1Response, error) {
			return model.GetMicroserviceV1Response{}, errors.New("internal Server Error")
		},
	}

	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)

	request := tPrepareGetRequest(route.ApiMicroserviceV1)

	testContext.Request = request

	controller := NewMicroserviceController(mock)
	controller.GetMicroserviceV1(testContext)

	actualStatusCode := recorder.Code

	var actualResponseBody model.ErrorResponse
	tPrepareResponse(recorder.Body, &actualResponseBody)

	require.Equal(t, expectedStatusCode, actualStatusCode)
	require.Equal(t, expectedResponseBody, actualResponseBody)
}
