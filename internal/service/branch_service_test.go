package service

import (
	"context"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestBranchService_CreateBranchV1_Success(t *testing.T) {
	branchID := int64(10)
	record := database.RootBranch{
		ID:             branchID,
		Name:           "Центральный",
		Address:        "Саранск, ул. Васенко, 9/2",
		ActivityStatus: database.RootBranchActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock := &mockBranchRepository{
		createBranchV1Func: func(ctx context.Context, params database.CreateBranchV1Params) (database.RootBranch, error) {
			require.Equal(t, "Центральный", params.PName)
			return record, nil
		},
	}

	svc := NewBranchService(mock)
	response, err := svc.CreateBranchV1(context.Background(), model.PostBranchV1Request{
		Name:    "Центральный",
		Address: "Саранск, ул. Васенко, 9/2",
	})

	require.NoError(t, err)
	require.Equal(t, branchID, response.ID)
	require.Equal(t, "active", response.ActivityStatus)
}

func TestBranchService_CreateBranchV1_NameConflict(t *testing.T) {
	mock := &mockBranchRepository{
		createBranchV1Func: func(ctx context.Context, params database.CreateBranchV1Params) (database.RootBranch, error) {
			return database.RootBranch{}, &pq.Error{Code: database.PgErrUniqueViolation}
		},
	}

	svc := NewBranchService(mock)
	_, err := svc.CreateBranchV1(context.Background(), model.PostBranchV1Request{
		Name:    "Дубль",
		Address: "Адрес",
	})

	require.ErrorIs(t, err, ErrBranchNameConflict)
}

func TestBranchService_GetBranchV1_NotFound(t *testing.T) {
	mock := &mockBranchRepository{
		getBranchV1Func: func(ctx context.Context, id int64) (database.RootBranch, error) {
			return database.RootBranch{}, &pq.Error{Code: database.PgErrNoDataFound}
		},
	}

	svc := NewBranchService(mock)
	_, err := svc.GetBranchV1(context.Background(), 999)

	require.ErrorIs(t, err, ErrBranchNotFound)
}
