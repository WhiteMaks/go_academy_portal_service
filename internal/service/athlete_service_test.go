package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/WhiteMaks/go_academy_portal_service/autogenerate/database"
	"github.com/WhiteMaks/go_academy_portal_service/internal/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type mockAthleteRepository struct {
	createFunc  func(context.Context, database.CreateAthleteV1Params) (database.RootAthlete, error)
	getFunc     func(context.Context, database.GetAthletesV1Params) ([]database.RootAthlete, error)
	countFunc   func(context.Context, database.CountAthletesV1Params) (int64, error)
	getByIDFunc func(context.Context, int64) (database.RootAthlete, error)
	updateFunc  func(context.Context, database.UpdateAthleteV1Params) (database.RootAthlete, error)
	archiveFunc func(context.Context, int64) (database.RootAthlete, error)
	restoreFunc func(context.Context, int64) (database.RootAthlete, error)
	coachFunc   func(context.Context, database.SetAthleteCoachV1Params) (database.RootAthlete, error)
}

func (m *mockAthleteRepository) CreateAthleteV1(ctx context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
	return m.createFunc(ctx, params)
}

func (m *mockAthleteRepository) GetAthletesV1(ctx context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error) {
	return m.getFunc(ctx, params)
}

func (m *mockAthleteRepository) CountAthletesV1(ctx context.Context, params database.CountAthletesV1Params) (int64, error) {
	return m.countFunc(ctx, params)
}

func (m *mockAthleteRepository) GetAthleteByIDV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockAthleteRepository) UpdateAthleteV1(ctx context.Context, params database.UpdateAthleteV1Params) (database.RootAthlete, error) {
	return m.updateFunc(ctx, params)
}

func (m *mockAthleteRepository) ArchiveAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return m.archiveFunc(ctx, id)
}

func (m *mockAthleteRepository) RestoreAthleteV1(ctx context.Context, id int64) (database.RootAthlete, error) {
	return m.restoreFunc(ctx, id)
}

func (m *mockAthleteRepository) SetAthleteCoachV1(ctx context.Context, params database.SetAthleteCoachV1Params) (database.RootAthlete, error) {
	return m.coachFunc(ctx, params)
}

func TestAthleteService_CreateAthleteV1_NormalizesAndMaps(t *testing.T) {
	blank := "   "
	phone := "+79991234567"
	parentPhone := "+441234567890"
	parentFirstName := "Анна"
	parentSecondName := "Иванова"
	updatedAt := time.Date(2025, 1, 2, 3, 4, 5, 123456000, time.UTC)
	var gotParams database.CreateAthleteV1Params
	repository := &mockAthleteRepository{
		createFunc: func(_ context.Context, params database.CreateAthleteV1Params) (database.RootAthlete, error) {
			gotParams = params
			return database.RootAthlete{
				ID: 9223372036854770000, SecondName: params.SecondName, FirstName: params.FirstName,
				MiddleName: params.MiddleName, BirthDate: params.BirthDate, PhotoUrl: params.PhotoUrl,
				Phone: params.Phone, ParentSecondName: params.ParentSecondName,
				ParentFirstName: params.ParentFirstName, ParentMiddleName: params.ParentMiddleName,
				ParentPhone: params.ParentPhone, IsActive: true, UpdatedAt: updatedAt,
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
	require.Equal(t, updatedAt, response.UpdatedAt)
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
	updatedAt := time.Date(2025, 1, 2, 3, 4, 5, 123456000, time.UTC)
	service := NewAthleteService(&mockAthleteRepository{getFunc: func(context.Context, database.GetAthletesV1Params) ([]database.RootAthlete, error) {
		return []database.RootAthlete{{ID: 1, FirstName: "A", SecondName: "B", BirthDate: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), UpdatedAt: updatedAt}}, nil
	}, countFunc: func(context.Context, database.CountAthletesV1Params) (int64, error) {
		return 1, nil
	}})
	response, err := service.GetAthletesV1(context.Background(), model.GetAthletesV1Query{IsActive: true, Limit: 50})
	require.NoError(t, err)
	require.Len(t, response.Athletes, 1)
	require.Nil(t, response.Athletes[0].Phone)
	require.Equal(t, utcDate(2000, time.January, 2), response.Athletes[0].BirthDate)
	require.Equal(t, updatedAt, response.Athletes[0].UpdatedAt)

	emptyService := NewAthleteService(&mockAthleteRepository{getFunc: func(context.Context, database.GetAthletesV1Params) ([]database.RootAthlete, error) {
		return nil, nil
	}, countFunc: func(context.Context, database.CountAthletesV1Params) (int64, error) {
		return 0, nil
	}})
	emptyResponse, err := emptyService.GetAthletesV1(context.Background(), model.GetAthletesV1Query{IsActive: true, Limit: 50})
	require.NoError(t, err)
	require.NotNil(t, emptyResponse.Athletes)
	require.Empty(t, emptyResponse.Athletes)
}

func TestNormalizeAthleteSearchQuery(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		want   string
		isDate bool
	}{
		{name: "empty", input: " \t ", want: ""},
		{name: "fio normalized", input: "  Ёлки\tИВАНОВ  ", want: "елки иванов"},
		{name: "date with dots and spaces", input: " 23. 04 ", want: "2304", isDate: true},
		{name: "date without separators", input: "2304", want: "2304", isDate: true},
		{name: "mixed remains fio", input: "Иванова 23.04", want: "иванова 23.04"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, isDate := NormalizeAthleteSearchQuery(tc.input)
			require.Equal(t, tc.want, got)
			require.Equal(t, tc.isDate, isDate)
		})
	}
}

func TestAthleteService_GetAthletesV1MapsFiltersPaginationAndCount(t *testing.T) {
	isCoach := false
	query := model.GetAthletesV1Query{IsActive: false, SearchQuery: "иванова анна", IsCoach: &isCoach, Limit: 12, Offset: 24}
	coachFilter := sql.NullBool{Bool: false, Valid: true}
	wantList := database.GetAthletesV1Params{IsActive: false, SearchQuery: "иванова анна", IsCoach: coachFilter, LimitCount: 12, OffsetCount: 24}
	wantCount := database.CountAthletesV1Params{IsActive: false, SearchQuery: "иванова анна", IsCoach: coachFilter}
	var gotList database.GetAthletesV1Params
	var gotCount database.CountAthletesV1Params
	repository := &mockAthleteRepository{
		getFunc: func(_ context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error) {
			gotList = params
			return nil, nil
		},
		countFunc: func(_ context.Context, params database.CountAthletesV1Params) (int64, error) {
			gotCount = params
			return 37, nil
		},
	}
	response, err := NewAthleteService(repository).GetAthletesV1(context.Background(), query)
	require.NoError(t, err)
	require.Equal(t, wantList, gotList)
	require.Equal(t, wantCount, gotCount)
	require.Equal(t, int64(37), response.Total)
	require.Equal(t, int32(12), response.Limit)
	require.Equal(t, int32(24), response.Offset)
}

func TestAthleteService_SetAthleteCoachReturnsDesiredState(t *testing.T) {
	var got database.SetAthleteCoachV1Params
	service := NewAthleteService(&mockAthleteRepository{coachFunc: func(_ context.Context, params database.SetAthleteCoachV1Params) (database.RootAthlete, error) {
		got = params
		return database.RootAthlete{ID: params.ID, IsActive: false, IsCoach: params.IsCoach}, nil
	}})
	response, err := service.SetAthleteCoachV1(context.Background(), 9, true)
	require.NoError(t, err)
	require.Equal(t, database.SetAthleteCoachV1Params{ID: 9, IsCoach: true}, got)
	require.Equal(t, "9", response.ID)
	require.False(t, response.IsActive)
	require.True(t, response.IsCoach)
}

func TestAthleteService_SetAthleteCoachMapsMissingAthlete(t *testing.T) {
	for _, dbErr := range []error{
		sql.ErrNoRows,
		&pq.Error{Code: "P0002"},
	} {
		service := NewAthleteService(&mockAthleteRepository{coachFunc: func(context.Context, database.SetAthleteCoachV1Params) (database.RootAthlete, error) {
			return database.RootAthlete{}, dbErr
		}})
		_, err := service.SetAthleteCoachV1(context.Background(), 404, true)
		require.ErrorIs(t, err, ErrAthleteNotFound)
	}
}

func TestAthleteService_GetAthletesV1MapsDateQuery(t *testing.T) {
	var gotList database.GetAthletesV1Params
	var gotCount database.CountAthletesV1Params
	repository := &mockAthleteRepository{
		getFunc: func(_ context.Context, params database.GetAthletesV1Params) ([]database.RootAthlete, error) {
			gotList = params
			return nil, nil
		},
		countFunc: func(_ context.Context, params database.CountAthletesV1Params) (int64, error) {
			gotCount = params
			return 0, nil
		},
	}
	_, err := NewAthleteService(repository).GetAthletesV1(context.Background(), model.GetAthletesV1Query{
		IsActive: true, SearchQuery: "2304", IsBirthDateQuery: true, Limit: 50,
	})
	require.NoError(t, err)
	require.True(t, gotList.IsBirthDateQuery)
	require.True(t, gotCount.IsBirthDateQuery)
}

func TestAthleteService_ArchiveRestoreClassifiesDatabaseStates(t *testing.T) {
	cases := []struct {
		name      string
		operation func(*athleteService) error
		dbCode    pq.ErrorCode
		wantErr   error
	}{
		{name: "archive already archived", operation: func(s *athleteService) error { _, err := s.ArchiveAthleteV1(context.Background(), 1); return err }, dbCode: "P4091", wantErr: ErrAthleteAlreadyArchived},
		{name: "restore already active", operation: func(s *athleteService) error { _, err := s.RestoreAthleteV1(context.Background(), 1); return err }, dbCode: "P4092", wantErr: ErrAthleteAlreadyActive},
		{name: "unknown archive id", operation: func(s *athleteService) error { _, err := s.ArchiveAthleteV1(context.Background(), 1); return err }, dbCode: "P0002", wantErr: ErrAthleteNotFound},
		{name: "unknown restore id", operation: func(s *athleteService) error { _, err := s.RestoreAthleteV1(context.Background(), 1); return err }, dbCode: "P0002", wantErr: ErrAthleteNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dbErr := &pq.Error{Code: tc.dbCode}
			repository := &mockAthleteRepository{
				archiveFunc: func(context.Context, int64) (database.RootAthlete, error) { return database.RootAthlete{}, dbErr },
				restoreFunc: func(context.Context, int64) (database.RootAthlete, error) { return database.RootAthlete{}, dbErr },
			}
			require.ErrorIs(t, tc.operation(&athleteService{repository: repository}), tc.wantErr)
		})
	}
}

func TestAthleteService_ArchiveRestoreMapSuccessfulState(t *testing.T) {
	archived := database.RootAthlete{ID: 7, IsActive: false}
	active := database.RootAthlete{ID: 7, IsActive: true}
	repository := &mockAthleteRepository{
		archiveFunc: func(_ context.Context, id int64) (database.RootAthlete, error) {
			require.Equal(t, int64(7), id)
			return archived, nil
		},
		restoreFunc: func(_ context.Context, id int64) (database.RootAthlete, error) {
			require.Equal(t, int64(7), id)
			return active, nil
		},
	}
	service := NewAthleteService(repository)
	archivedResponse, err := service.ArchiveAthleteV1(context.Background(), 7)
	require.NoError(t, err)
	require.False(t, archivedResponse.IsActive)
	activeResponse, err := service.RestoreAthleteV1(context.Background(), 7)
	require.NoError(t, err)
	require.True(t, activeResponse.IsActive)
}

func TestAthleteService_UpdateAthleteV1PreservesArchivedState(t *testing.T) {
	service := NewAthleteService(&mockAthleteRepository{updateFunc: func(_ context.Context, params database.UpdateAthleteV1Params) (database.RootAthlete, error) {
		return database.RootAthlete{ID: params.ID, FirstName: params.FirstName, SecondName: params.SecondName, BirthDate: params.BirthDate, IsActive: false}, nil
	}})
	response, err := service.UpdateAthleteV1(context.Background(), 7, validUpdateRequest())
	require.NoError(t, err)
	require.False(t, response.IsActive)
}

func TestAthleteService_UpdateAthleteV1_MapsResponse(t *testing.T) {
	old := time.Date(2025, 1, 2, 3, 4, 5, 123456000, time.UTC)
	advanced := old.Add(time.Microsecond)
	request := validUpdateRequest()
	photo := sql.NullString{String: "https://example.com/photo.jpg", Valid: true}
	gotParams := database.UpdateAthleteV1Params{}
	service := NewAthleteService(&mockAthleteRepository{updateFunc: func(_ context.Context, params database.UpdateAthleteV1Params) (database.RootAthlete, error) {
		gotParams = params
		return database.RootAthlete{ID: 7, FirstName: params.FirstName, SecondName: params.SecondName, BirthDate: params.BirthDate,
			PhotoUrl: photo, IsActive: true, UpdatedAt: advanced}, nil
	}})
	response, err := service.UpdateAthleteV1(context.Background(), 7, request)
	require.NoError(t, err)
	require.Equal(t, "Иван", gotParams.FirstName)
	require.Equal(t, "Иванов", gotParams.SecondName)
	require.Equal(t, advanced, response.UpdatedAt)
	require.Equal(t, photo.String, *response.PhotoURL)
	require.True(t, response.IsActive)
}

func TestAthleteService_GetAthleteByIDV1_MapsInactiveOrMissingToNotFound(t *testing.T) {
	service := NewAthleteService(&mockAthleteRepository{getByIDFunc: func(context.Context, int64) (database.RootAthlete, error) {
		return database.RootAthlete{}, sql.ErrNoRows
	}})
	_, err := service.GetAthleteByIDV1(context.Background(), 7)
	require.ErrorIs(t, err, ErrAthleteNotFound)
}

func validUpdateRequest() model.PutAthleteV1Request {
	return model.PutAthleteV1Request{FirstName: " Иван ", SecondName: " Иванов ",
		BirthDate: utcDate(2010, time.January, 2)}
}

func strPtr(value string) *string { return &value }

func utcDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
