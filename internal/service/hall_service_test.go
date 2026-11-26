package service

import (
	"context"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/stretchr/testify/require"
)

func TestHallService_CreateHallV1_Success(t *testing.T) {
	hallID := int64(5)
	branchID := int64(1)
	hallRecord := database.RootHall{
		ID:             hallID,
		BranchID:       branchID,
		Name:           "Основной зал",
		Capacity:       50,
		ActivityStatus: database.RootHallActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	branchRecord := database.RootBranch{
		ID:      branchID,
		Name:    "Центральный",
		Address: "Саранск",
	}

	hallMock := &mockHallRepository{
		createHallV1Func: func(ctx context.Context, params database.CreateHallV1Params) (database.RootHall, error) {
			require.Equal(t, branchID, params.PBranchID)
			require.Equal(t, int32(50), params.PCapacity)
			return hallRecord, nil
		},
	}
	branchMock := &mockBranchRepository{
		getBranchV1Func: func(ctx context.Context, id int64) (database.RootBranch, error) {
			return branchRecord, nil
		},
	}

	svc := NewHallService(hallMock, branchMock)
	response, err := svc.CreateHallV1(context.Background(), model.PostHallV1Request{
		BranchID: branchID,
		Name:     "Основной зал",
		Capacity: 50,
	})

	require.NoError(t, err)
	require.Equal(t, hallID, response.ID)
	require.NotNil(t, response.BranchName)
	require.Equal(t, "Центральный", *response.BranchName)
}

func TestHallService_UpdateHallV1_SuccessWithoutBranchIDInBody(t *testing.T) {
	hallID := int64(3)
	branchID := int64(1)
	hallRecord := database.RootHall{
		ID:             hallID,
		BranchID:       branchID,
		Name:           "Малый зал",
		Capacity:       25,
		ActivityStatus: database.RootHallActivityStatusInactive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	branchRecord := database.RootBranch{ID: branchID, Name: "Центральный", Address: "Саранск"}

	hallMock := &mockHallRepository{
		updateHallV1Func: func(ctx context.Context, params database.UpdateHallV1Params) (database.RootHall, error) {
			require.Equal(t, hallID, params.PID)
			require.Equal(t, "Малый зал", params.PName)
			require.Equal(t, int32(25), params.PCapacity)
			require.Equal(t, "inactive", params.PActivityStatus)
			return hallRecord, nil
		},
	}
	branchMock := &mockBranchRepository{
		getBranchV1Func: func(ctx context.Context, id int64) (database.RootBranch, error) {
			return branchRecord, nil
		},
	}

	svc := NewHallService(hallMock, branchMock)
	response, err := svc.UpdateHallV1(context.Background(), hallID, model.PutHallV1Request{
		Name:           "Малый зал",
		Capacity:       25,
		ActivityStatus: "inactive",
	})

	require.NoError(t, err)
	require.Equal(t, hallID, response.ID)
	require.Equal(t, "inactive", response.ActivityStatus)
}

func TestHallService_CreateHallV1_InvalidCapacity(t *testing.T) {
	svc := NewHallService(&mockHallRepository{}, &mockBranchRepository{})
	_, err := svc.CreateHallV1(context.Background(), model.PostHallV1Request{
		BranchID: 1,
		Name:     "Зал",
		Capacity: 0,
	})

	require.ErrorIs(t, err, ErrHallValidation)
}
