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

func TestAthleteService_CreateAthleteV1_Success(t *testing.T) {
	athleteID := int64(42)
	record := database.RootAthlete{
		ID:             athleteID,
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC),
		Phone:          "+79001234567",
		ActivityStatus: database.RootAthleteActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock := &mockAthleteRepository{
		createAthleteV1Func: func(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
			require.Equal(t, "Ivanov", params.PLastName)
			require.Equal(t, "+79001234567", params.PPhone)
			return record, nil
		},
	}

	svc := NewAthleteService(mock)
	response, err := svc.CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		LastName:  "Ivanov",
		FirstName: "Ivan",
		BirthDate: model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:     "+7 900 123-45-67",
	})

	require.NoError(t, err)
	require.Equal(t, athleteID, response.ID)
	require.Equal(t, "active", response.ActivityStatus)
}

func TestAthleteService_CreateAthleteV1_PhoneConflict(t *testing.T) {
	mock := &mockAthleteRepository{
		createAthleteV1Func: func(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
			return database.RootAthlete{}, &pq.Error{Code: database.PgErrUniqueViolation}
		},
	}

	svc := NewAthleteService(mock)
	_, err := svc.CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		LastName:  "Petrov",
		FirstName: "Petr",
		BirthDate: model.DateOnly{Time: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)},
		Phone:     "+79001112233",
	})

	require.ErrorIs(t, err, ErrAthletePhoneConflict)
}

func TestAthleteService_CreateAthleteV1_ParentPartial(t *testing.T) {
	svc := NewAthleteService(&mockAthleteRepository{})
	parentLast := "Sidorov"
	_, err := svc.CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:          "+79001234567",
		ParentLastName: &parentLast,
	})

	require.ErrorIs(t, err, ErrAthleteParentPartial)
}

func TestAthleteService_GetAthletesV1_WithQuery(t *testing.T) {
	record := database.RootAthlete{
		ID:             1,
		LastName:       "Ivanov",
		FirstName:      "Ivan",
		BirthDate:      time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC),
		Phone:          "+79001234567",
		ActivityStatus: database.RootAthleteActivityStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	mock := &mockAthleteRepository{
		getAthletesV1Func: func(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error) {
			require.True(t, params.PQ.Valid)
			require.Equal(t, "15.03", params.PQ.String)
			require.Equal(t, int32(50), params.PLimit.Int32)
			return []database.RootAthlete{record}, nil
		},
		countAthletesV1Func: func(ctx context.Context, pQ sql.NullString) (int64, error) {
			require.True(t, pQ.Valid)
			require.Equal(t, "15.03", pQ.String)
			return 1, nil
		},
	}

	svc := NewAthleteService(mock)
	response, err := svc.GetAthletesV1(context.Background(), model.GetAthletesV1Query{Q: "15.03"})
	require.NoError(t, err)
	require.Equal(t, int64(1), response.Total)
	require.Len(t, response.Athletes, 1)
	require.Equal(t, 50, response.Limit)
}

func TestAthleteService_CreateAthleteV1_FutureBirthDate(t *testing.T) {
	svc := NewAthleteService(&mockAthleteRepository{})
	future := time.Now().UTC().AddDate(0, 0, 1)
	_, err := svc.CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		LastName:  "Ivanov",
		FirstName: "Ivan",
		BirthDate: model.DateOnly{Time: future},
		Phone:     "+79001234567",
	})
	require.ErrorIs(t, err, ErrAthleteValidation)
}

func TestAthleteService_CreateAthleteV1_FieldTooLong(t *testing.T) {
	svc := NewAthleteService(&mockAthleteRepository{})
	long := strings.Repeat("x", model.AthleteNameMaxLen+1)
	_, err := svc.CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		LastName:  long,
		FirstName: "Ivan",
		BirthDate: model.DateOnly{Time: time.Date(2015, 3, 15, 0, 0, 0, 0, time.UTC)},
		Phone:     "+79001234567",
	})
	require.ErrorIs(t, err, ErrAthleteValidation)
}

func TestAthleteService_GetAthletesV1_InvalidQuery(t *testing.T) {
	svc := NewAthleteService(&mockAthleteRepository{})
	_, err := svc.GetAthletesV1(context.Background(), model.GetAthletesV1Query{Limit: model.AthletesListMaxLimit + 1})
	require.ErrorIs(t, err, ErrAthleteListQuery)
}

func TestAthleteService_GetAthletesV1_CapsLimitAtMax(t *testing.T) {
	mock := &mockAthleteRepository{
		getAthletesV1Func: func(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error) {
			require.Equal(t, int32(model.AthletesListMaxLimit), params.PLimit.Int32)
			return nil, nil
		},
		countAthletesV1Func: func(ctx context.Context, pQ sql.NullString) (int64, error) {
			return 0, nil
		},
	}
	svc := NewAthleteService(mock)
	response, err := svc.GetAthletesV1(context.Background(), model.GetAthletesV1Query{Limit: model.AthletesListMaxLimit})
	require.NoError(t, err)
	require.Equal(t, model.AthletesListMaxLimit, response.Limit)
}
