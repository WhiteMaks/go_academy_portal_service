package database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestBranchAndHallV1_HappyPathAndCascade(t *testing.T) {
	ctx := context.Background()
	suffix := randomTestSuffix()

	branch, err := testQueries.CreateBranchV1(ctx, CreateBranchV1Params{
		PName:    "Branch " + suffix,
		PAddress: "Address " + suffix,
		PNote:    sql.NullString{},
	})
	require.NoError(t, err)
	require.NotZero(t, branch.ID)

	hall, err := testQueries.CreateHallV1(ctx, CreateHallV1Params{
		PBranchID: branch.ID,
		PName:     "Hall " + suffix,
		PSubtitle: sql.NullString{},
		PCapacity: 30,
		PNote:     sql.NullString{},
	})
	require.NoError(t, err)
	require.Equal(t, RootHallActivityStatusActive, hall.ActivityStatus)

	updated, err := testQueries.UpdateBranchV1(ctx, UpdateBranchV1Params{
		PID:             branch.ID,
		PName:           branch.Name,
		PAddress:        branch.Address,
		PNote:           sql.NullString{},
		PActivityStatus: "inactive",
	})
	require.NoError(t, err)
	require.Equal(t, RootBranchActivityStatusInactive, updated.ActivityStatus)

	hallAfter, err := testQueries.GetHallV1(ctx, hall.ID)
	require.NoError(t, err)
	require.Equal(t, RootHallActivityStatusInactive, hallAfter.ActivityStatus)

	reactivated, err := testQueries.UpdateBranchV1(ctx, UpdateBranchV1Params{
		PID:             branch.ID,
		PName:           branch.Name,
		PAddress:        branch.Address,
		PNote:           sql.NullString{},
		PActivityStatus: "active",
	})
	require.NoError(t, err)
	require.Equal(t, RootBranchActivityStatusActive, reactivated.ActivityStatus)

	hallAfterReactivate, err := testQueries.GetHallV1(ctx, hall.ID)
	require.NoError(t, err)
	require.Equal(t, RootHallActivityStatusInactive, hallAfterReactivate.ActivityStatus)
}

func TestCreateBranchV1_DuplicateName(t *testing.T) {
	ctx := context.Background()
	suffix := randomTestSuffix()
	name := "DupBranch " + suffix

	_, err := testQueries.CreateBranchV1(ctx, CreateBranchV1Params{
		PName:    name,
		PAddress: "Addr 1",
		PNote:    sql.NullString{},
	})
	require.NoError(t, err)

	_, err = testQueries.CreateBranchV1(ctx, CreateBranchV1Params{
		PName:    name,
		PAddress: "Addr 2",
		PNote:    sql.NullString{},
	})
	require.Error(t, err)
	var pqErr *pq.Error
	require.ErrorAs(t, err, &pqErr)
	require.Equal(t, PgErrUniqueViolation, pqErr.Code)
}

func TestUpdateHallV1_CannotActivateWhileBranchInactive(t *testing.T) {
	ctx := context.Background()
	suffix := randomTestSuffix()

	branch, err := testQueries.CreateBranchV1(ctx, CreateBranchV1Params{
		PName:    "HallAct " + suffix,
		PAddress: "Addr",
		PNote:    sql.NullString{},
	})
	require.NoError(t, err)

	hall, err := testQueries.CreateHallV1(ctx, CreateHallV1Params{
		PBranchID: branch.ID,
		PName:     "Hall",
		PCapacity: 10,
		PNote:     sql.NullString{},
	})
	require.NoError(t, err)

	_, err = testQueries.UpdateBranchV1(ctx, UpdateBranchV1Params{
		PID:             branch.ID,
		PName:           branch.Name,
		PAddress:        branch.Address,
		PNote:           sql.NullString{},
		PActivityStatus: "inactive",
	})
	require.NoError(t, err)

	_, err = testQueries.UpdateHallV1(ctx, UpdateHallV1Params{
		PID:             hall.ID,
		PName:           hall.Name,
		PSubtitle:       hall.Subtitle,
		PCapacity:       hall.Capacity,
		PNote:           hall.Note,
		PActivityStatus: "inactive",
	})
	require.NoError(t, err)

	_, err = testQueries.UpdateHallV1(ctx, UpdateHallV1Params{
		PID:             hall.ID,
		PName:           "Renamed",
		PSubtitle:       hall.Subtitle,
		PCapacity:       hall.Capacity,
		PNote:           hall.Note,
		PActivityStatus: "inactive",
	})
	require.NoError(t, err)

	_, err = testQueries.UpdateHallV1(ctx, UpdateHallV1Params{
		PID:             hall.ID,
		PName:           "Renamed",
		PSubtitle:       hall.Subtitle,
		PCapacity:       hall.Capacity,
		PNote:           hall.Note,
		PActivityStatus: "active",
	})
	require.Error(t, err)
	var pqErr *pq.Error
	require.ErrorAs(t, err, &pqErr)
	require.Equal(t, PgErrInvalidParameterValue, pqErr.Code)
}

func TestCreateHallV1_InactiveBranchAllowed(t *testing.T) {
	ctx := context.Background()
	suffix := randomTestSuffix()

	branch, err := testQueries.CreateBranchV1(ctx, CreateBranchV1Params{
		PName:    "InactiveParent " + suffix,
		PAddress: "Addr",
		PNote:    sql.NullString{},
	})
	require.NoError(t, err)

	_, err = testQueries.UpdateBranchV1(ctx, UpdateBranchV1Params{
		PID:             branch.ID,
		PName:           branch.Name,
		PAddress:        branch.Address,
		PNote:           sql.NullString{},
		PActivityStatus: "inactive",
	})
	require.NoError(t, err)

	hall, err := testQueries.CreateHallV1(ctx, CreateHallV1Params{
		PBranchID: branch.ID,
		PName:     "New hall",
		PCapacity: 10,
		PNote:     sql.NullString{},
	})
	require.NoError(t, err)
	require.Equal(t, branch.ID, hall.BranchID)
	require.Equal(t, RootHallActivityStatusInactive, hall.ActivityStatus)
}

func randomTestSuffix() string {
	return util.RandomString(8)
}
