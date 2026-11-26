package service

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var coachTestBirthDate = model.DateOnly{Time: time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC)}

func TestCoachService_CreateCoachV1_Success(t *testing.T) {
	coachID := int64(42)
	record := database.RootCoach{
		ID:             coachID,
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachTestBirthDate.Time,
		Phone:          "+79001234567",
		Specialization: "Judo",
		ActivityStatus: database.RootCoachActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock := &mockCoachRepository{
		createCoachV1Func: func(ctx context.Context, params database.CreateCoachV1Params) (database.RootCoach, error) {
			require.Equal(t, "Sidorov", params.PLastName)
			require.Equal(t, "+79001234567", params.PPhone)
			require.Equal(t, "Judo", params.PSpecialization)
			require.Equal(t, coachTestBirthDate.Time, params.PBirthDate)
			return record, nil
		},
	}

	svc := NewCoachService(mock)
	response, err := svc.CreateCoachV1(context.Background(), model.PostCoachV1Request{
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachTestBirthDate,
		Phone:          "+7 900 123-45-67",
		Specialization: "Judo",
	})

	require.NoError(t, err)
	require.Equal(t, coachID, response.ID)
	require.Equal(t, "active", response.ActivityStatus)
	require.Equal(t, coachTestBirthDate, response.BirthDate)
}

func TestCoachService_CreateCoachV1_PhoneConflict(t *testing.T) {
	mock := &mockCoachRepository{
		createCoachV1Func: func(ctx context.Context, params database.CreateCoachV1Params) (database.RootCoach, error) {
			return database.RootCoach{}, &pq.Error{Code: database.PgErrUniqueViolation}
		},
	}

	svc := NewCoachService(mock)
	_, err := svc.CreateCoachV1(context.Background(), model.PostCoachV1Request{
		LastName:       "Petrov",
		FirstName:      "Petr",
		BirthDate:      coachTestBirthDate,
		Phone:          "+79001112233",
		Specialization: "Boxing",
	})

	require.ErrorIs(t, err, ErrCoachPhoneConflict)
}

func TestCoachService_GetCoachesV1_WithQuery(t *testing.T) {
	record := database.RootCoach{
		ID:             1,
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachTestBirthDate.Time,
		Phone:          "+79001234567",
		Specialization: "Judo",
		ActivityStatus: database.RootCoachActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock := &mockCoachRepository{
		getCoachesV1Func: func(ctx context.Context, params database.GetCoachesV1Params) ([]database.RootCoach, error) {
			require.True(t, params.PQ.Valid)
			require.Equal(t, "Judo", params.PQ.String)
			require.Equal(t, int32(model.CoachesListDefaultLimit), params.PLimit.Int32)
			return []database.RootCoach{record}, nil
		},
		countCoachesV1Func: func(ctx context.Context, pQ sql.NullString) (int64, error) {
			require.True(t, pQ.Valid)
			require.Equal(t, "Judo", pQ.String)
			return 1, nil
		},
	}

	svc := NewCoachService(mock)
	response, err := svc.GetCoachesV1(context.Background(), model.GetCoachesV1Query{Q: "Judo"})
	require.NoError(t, err)
	require.Equal(t, int64(1), response.Total)
	require.Len(t, response.Coaches, 1)
	require.Equal(t, model.CoachesListDefaultLimit, response.Limit)
	require.Equal(t, coachTestBirthDate, response.Coaches[0].BirthDate)
}

func TestCoachService_GetCoachesV1_SearchByBirthDateYear(t *testing.T) {
	record := database.RootCoach{
		ID:             2,
		LastName:       "Petrov",
		FirstName:      "Petr",
		BirthDate:      time.Date(1985, 11, 3, 0, 0, 0, 0, time.UTC),
		Phone:          "+79001112233",
		Specialization: "Boxing",
		ActivityStatus: database.RootCoachActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock := &mockCoachRepository{
		getCoachesV1Func: func(ctx context.Context, params database.GetCoachesV1Params) ([]database.RootCoach, error) {
			require.True(t, params.PQ.Valid)
			require.Equal(t, "1985", params.PQ.String)
			return []database.RootCoach{record}, nil
		},
		countCoachesV1Func: func(ctx context.Context, pQ sql.NullString) (int64, error) {
			return 1, nil
		},
	}

	svc := NewCoachService(mock)
	response, err := svc.GetCoachesV1(context.Background(), model.GetCoachesV1Query{Q: "1985"})
	require.NoError(t, err)
	require.Len(t, response.Coaches, 1)
	require.Equal(t, model.DateOnly{Time: record.BirthDate}, response.Coaches[0].BirthDate)
}

func TestCoachService_CreateCoachV1_FutureBirthDate(t *testing.T) {
	svc := NewCoachService(&mockCoachRepository{})
	future := time.Now().UTC().AddDate(0, 0, 1)
	_, err := svc.CreateCoachV1(context.Background(), model.PostCoachV1Request{
		LastName:       "Petrov",
		FirstName:      "Petr",
		BirthDate:      model.DateOnly{Time: future},
		Phone:          "+79001234567",
		Specialization: "Judo",
	})
	require.ErrorIs(t, err, ErrCoachValidation)
}

func TestCoachService_CreateCoachV1_MissingBirthDate(t *testing.T) {
	svc := NewCoachService(&mockCoachRepository{})
	_, err := svc.CreateCoachV1(context.Background(), model.PostCoachV1Request{
		LastName:       "Petrov",
		FirstName:      "Petr",
		Phone:          "+79001234567",
		Specialization: "Judo",
	})
	require.ErrorIs(t, err, ErrCoachValidation)
}

func TestCoachService_GetCoachV1_NotFound(t *testing.T) {
	mock := &mockCoachRepository{
		getCoachV1Func: func(ctx context.Context, id int64) (database.RootCoach, error) {
			return database.RootCoach{}, &pq.Error{Code: database.PgErrNoDataFound}
		},
	}

	svc := NewCoachService(mock)
	_, err := svc.GetCoachV1(context.Background(), 99)
	require.ErrorIs(t, err, ErrCoachNotFound)
}

func TestCoachService_UpdateCoachV1_Success(t *testing.T) {
	coachID := int64(7)
	updatedBirth := time.Date(1988, 7, 20, 0, 0, 0, 0, time.UTC)
	record := database.RootCoach{
		ID:             coachID,
		LastName:       "Petrov",
		FirstName:      "Petr",
		BirthDate:      updatedBirth,
		Phone:          "+79001234567",
		Specialization: "Karate",
		ActivityStatus: database.RootCoachActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock := &mockCoachRepository{
		updateCoachV1Func: func(ctx context.Context, params database.UpdateCoachV1Params) (database.RootCoach, error) {
			require.Equal(t, coachID, params.PID)
			require.Equal(t, "Petrov", params.PLastName)
			require.Equal(t, "active", params.PActivityStatus)
			require.Equal(t, updatedBirth, params.PBirthDate)
			return record, nil
		},
	}

	svc := NewCoachService(mock)
	response, err := svc.UpdateCoachV1(context.Background(), coachID, model.PutCoachV1Request{
		LastName:       "Petrov",
		FirstName:      "Petr",
		BirthDate:      model.DateOnly{Time: updatedBirth},
		Phone:          "+79001234567",
		Specialization: "Karate",
		ActivityStatus: "active",
	})

	require.NoError(t, err)
	require.Equal(t, coachID, response.ID)
	require.Equal(t, model.DateOnly{Time: updatedBirth}, response.BirthDate)
}

func TestCoachService_UpdateCoachV1_NotFound(t *testing.T) {
	mock := &mockCoachRepository{
		updateCoachV1Func: func(ctx context.Context, params database.UpdateCoachV1Params) (database.RootCoach, error) {
			return database.RootCoach{}, &pq.Error{Code: database.PgErrNoDataFound}
		},
	}

	svc := NewCoachService(mock)
	_, err := svc.UpdateCoachV1(context.Background(), 99, model.PutCoachV1Request{
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachTestBirthDate,
		Phone:          "+79001234567",
		Specialization: "Judo",
		ActivityStatus: "active",
	})

	require.ErrorIs(t, err, ErrCoachNotFound)
}

func TestCoachService_UpdateCoachV1_PhoneConflict(t *testing.T) {
	mock := &mockCoachRepository{
		updateCoachV1Func: func(ctx context.Context, params database.UpdateCoachV1Params) (database.RootCoach, error) {
			return database.RootCoach{}, &pq.Error{Code: database.PgErrUniqueViolation}
		},
	}

	svc := NewCoachService(mock)
	_, err := svc.UpdateCoachV1(context.Background(), 7, model.PutCoachV1Request{
		LastName:       "Sidorov",
		FirstName:      "Sergey",
		BirthDate:      coachTestBirthDate,
		Phone:          "+79005554433",
		Specialization: "Judo",
		ActivityStatus: "active",
	})

	require.ErrorIs(t, err, ErrCoachPhoneConflict)
}

func TestCoachService_CreateCoachV1_FieldTooLong(t *testing.T) {
	svc := NewCoachService(&mockCoachRepository{})
	long := strings.Repeat("x", model.CoachNameMaxLen+1)
	_, err := svc.CreateCoachV1(context.Background(), model.PostCoachV1Request{
		LastName:       long,
		FirstName:      "Ivan",
		BirthDate:      coachTestBirthDate,
		Phone:          "+79001234567",
		Specialization: "Judo",
	})
	require.ErrorIs(t, err, ErrCoachValidation)
}

func TestCoachService_GetCoachesV1_InvalidQuery(t *testing.T) {
	svc := NewCoachService(&mockCoachRepository{})
	_, err := svc.GetCoachesV1(context.Background(), model.GetCoachesV1Query{Limit: model.CoachesListMaxLimit + 1})
	require.ErrorIs(t, err, ErrCoachListQuery)
}
