package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/stretchr/testify/require"
)

func TestMicroserviceService_GetMicroserviceV1_Success(t *testing.T) {
	expectedResponse := model.GetMicroserviceV1Response{
		ReadyToUse: true,
	}

	mock := &mockUserRepository{
		isAdminCreationAllowedV1Func: func(ctx context.Context) (sql.NullBool, error) {
			return sql.NullBool{Bool: !expectedResponse.ReadyToUse, Valid: true}, nil
		},
	}

	service := NewMicroserviceService(mock)

	actualResponse, err := service.GetMicroserviceV1(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedResponse, actualResponse)
}

func TestMicroserviceService_GetMicroserviceV1_Error(t *testing.T) {
	mock := &mockUserRepository{
		isAdminCreationAllowedV1Func: func(ctx context.Context) (sql.NullBool, error) {
			return sql.NullBool{}, errors.New("repository error")
		},
	}

	service := NewMicroserviceService(mock)

	_, err := service.GetMicroserviceV1(context.Background())

	require.Error(t, err)
}
