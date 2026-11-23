package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/stretchr/testify/require"
)

type mockAthleteRepository struct {
	createFunc func(context.Context, database.CreateAthleteV1Params) (database.RootAthlete, error)
	getFunc    func(context.Context) ([]database.RootAthlete, error)
}

func (m *mockAthleteRepository) CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
	return m.createFunc(ctx, params)
}

func (m *mockAthleteRepository) GetAthletesV1(ctx context.Context) ([]database.RootAthlete, error) {
	return m.getFunc(ctx)
}

func TestAthleteService_CreateAthleteV1_NormalizesAndMaps(t *testing.T) {
	blank := "   "
	phone := "+79991234567"
	parentPhone := "+441234567890"
	parentFirstName := "Анна"
	parentSecondName := "Иванова"
	var gotParams database.CreateAthleteV1Params
	repository := &mockAthleteRepository{
		createFunc: func(_ context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
			gotParams = params
			return database.RootAthlete{
				ID: 9223372036854770000, SecondName: params.SecondName, FirstName: params.FirstName,
				MiddleName: params.MiddleName, BirthDate: params.BirthDate, PhotoUrl: params.PhotoUrl,
				Phone: params.Phone, ParentSecondName: params.ParentSecondName,
				ParentFirstName: params.ParentFirstName, ParentMiddleName: params.ParentMiddleName,
				ParentPhone: params.ParentPhone, IsActive: true,
			}, nil
		},
	}
	service := NewAthleteService(repository)
	response, err := service.CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		SecondName: " Иванов ", FirstName: "Иван", BirthDate: utcDate(2010, time.February, 3),
		MiddleName: &blank, Phone: phone, ParentPhone: parentPhone,
		ParentFirstName: &parentFirstName, ParentSecondName: &parentSecondName,
	})
	require.NoError(t, err)
	require.Equal(t, "Иванов", gotParams.SecondName)
	require.Equal(t, "Иван", gotParams.FirstName)
	require.False(t, gotParams.MiddleName.Valid)
	require.Equal(t, "+79991234567", gotParams.Phone.String)
	require.True(t, gotParams.Phone.Valid)
	require.Equal(t, "Иванова", gotParams.ParentSecondName.String)
	require.Equal(t, "Анна", gotParams.ParentFirstName.String)
	require.Equal(t, "+441234567890", gotParams.ParentPhone.String)
	require.True(t, gotParams.ParentPhone.Valid)
	require.Equal(t, "9223372036854770000", response.ID)
	require.Equal(t, utcDate(2010, time.February, 3), response.BirthDate)
	require.Nil(t, response.MiddleName)
	require.Equal(t, "Иванова", *response.ParentSecondName)
	require.Equal(t, "+441234567890", *response.ParentPhone)
	require.True(t, response.IsActive)
}

func TestAthleteService_CreateAthleteV1_NormalizesEmptyOptionalPhones(t *testing.T) {
	repository := &mockAthleteRepository{createFunc: func(_ context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
		require.False(t, params.Phone.Valid)
		require.False(t, params.ParentPhone.Valid)
		return database.RootAthlete{ID: 1, FirstName: params.FirstName, SecondName: params.SecondName,
			BirthDate: params.BirthDate, IsActive: true}, nil
	}}
	_, err := NewAthleteService(repository).CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		FirstName: "Иван", SecondName: "Иванов", BirthDate: utcDate(2010, time.January, 2), Phone: "", ParentPhone: "",
	})
	require.NoError(t, err)
}

func TestAthleteService_CreateAthleteV1_AcceptsPhoneLengthBoundaries(t *testing.T) {
	for _, phone := range []string{"+1234567890", "+123456789012345"} {
		t.Run(phone, func(t *testing.T) {
			repository := &mockAthleteRepository{createFunc: func(_ context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
				require.Equal(t, phone, params.Phone.String)
				require.True(t, params.Phone.Valid)
				return database.RootAthlete{ID: 1, FirstName: params.FirstName, SecondName: params.SecondName,
					BirthDate: params.BirthDate, Phone: params.Phone, IsActive: true}, nil
			}}
			_, err := NewAthleteService(repository).CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
				FirstName: "A", SecondName: "B", BirthDate: utcDate(2000, time.January, 1), Phone: phone,
			})
			require.NoError(t, err)
		})
	}
}

func TestAthleteService_CreateAthleteV1_RepositoryError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	repository := &mockAthleteRepository{createFunc: func(context.Context, database.CreateAthleteV1Params) (database.RootAthlete, error) {
		return database.RootAthlete{}, wantErr
	}}
	_, err := NewAthleteService(repository).CreateAthleteV1(context.Background(), model.PostAthleteV1Request{
		FirstName: "A", SecondName: "B", BirthDate: utcDate(2000, time.January, 1),
	})
	require.ErrorIs(t, err, wantErr)
}

func TestAthleteService_GetAthletesV1ReturnsEmptyArrayAndNullableFields(t *testing.T) {
	service := NewAthleteService(&mockAthleteRepository{getFunc: func(context.Context) ([]database.RootAthlete, error) {
		return []database.RootAthlete{{ID: 1, FirstName: "A", SecondName: "B", BirthDate: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)}}, nil
	}})
	response, err := service.GetAthletesV1(context.Background())
	require.NoError(t, err)
	require.Len(t, response.Athletes, 1)
	require.Nil(t, response.Athletes[0].Phone)
	require.Equal(t, utcDate(2000, time.January, 2), response.Athletes[0].BirthDate)

	emptyService := NewAthleteService(&mockAthleteRepository{getFunc: func(context.Context) ([]database.RootAthlete, error) {
		return nil, nil
	}})
	emptyResponse, err := emptyService.GetAthletesV1(context.Background())
	require.NoError(t, err)
	require.NotNil(t, emptyResponse.Athletes)
	require.Empty(t, emptyResponse.Athletes)
}

func strPtr(value string) *string { return &value }

func utcDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
