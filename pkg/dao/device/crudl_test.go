// +build !unit

package device

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-device-" +
			random.String(16), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgId, createDev.OrgId)
		require.Equal(t, dev.UniqId, createDev.UniqId)
		require.Equal(t, dev.Status, createDev.Status)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with uppercase UniqId", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID,
			UniqId: strings.ToUpper("dao-device-" + random.String(16)),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgId, createDev.OrgId)
		require.Equal(t, strings.ToLower(dev.UniqId), createDev.UniqId)
		require.Equal(t, dev.Status, createDev.Status)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID,
			UniqId: "dao-device-" + random.String(40)}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-device-" +
		random.String(16), Status: []common.Status{common.Status_ACTIVE,
		common.Status_DISABLED}[random.Intn(2)]}
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, createDev.Id, createDev.OrgId)
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDev)
	})

	t.Run("Read device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, uuid.NewString(),
			uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, createDev.Id,
			uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read device by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, random.String(10),
			createDev.OrgId)
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestReadByUniqID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-device-" +
		random.String(16), Status: []common.Status{common.Status_ACTIVE,
		common.Status_DISABLED}[random.Intn(2)]}
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDev, err := globalDevDAO.ReadByUniqID(ctx, createDev.UniqId)
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDev)
	})

	t.Run("Read device by unknown UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDev, err := globalDevDAO.ReadByUniqID(ctx, uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID,
			UniqId: "dao-device-" + random.String(16)}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "dao-device-" + random.String(16)
		createDev.Status = []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev.UniqId, updateDev.UniqId)
		require.Equal(t, createDev.Status, updateDev.Status)
		require.Equal(t, createDev.CreatedAt, updateDev.CreatedAt)
		require.True(t, updateDev.UpdatedAt.AsTime().After(
			updateDev.CreatedAt.AsTime()))
		require.WithinDuration(t, createDev.CreatedAt.AsTime(),
			updateDev.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Update unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownDevice := &api.Device{Id: uuid.NewString(),
			OrgId: createOrg.ID, UniqId: "dao-device-" + random.String(16),
			Token: uuid.NewString()}
		updateDev, err := globalDevDAO.Update(ctx, unknownDevice)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-device-" +
			random.String(16), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.OrgId = uuid.NewString()
		createDev.UniqId = "dao-device-" + random.String(16)

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-device-" +
			random.String(16), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "dao-device-" + random.String(40)

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-device-" +
			random.String(16), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDevDAO.Delete(ctx, createDev.Id, createOrg.ID)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			readDev, err := globalDevDAO.Read(ctx, createDev.Id,
				createOrg.ID)
			t.Logf("readDev, err: %+v, %v", readDev, err)
			require.Nil(t, readDev)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalDevDAO.Delete(ctx, uuid.NewString(), createOrg.ID)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID,
			UniqId: "dao-device-" + random.String(16), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDevDAO.Delete(ctx, createDev.Id, uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	devIDs := []string{}
	devStatuses := []common.Status{}
	devTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-device-" +
			random.String(16), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		devIDs = append(devIDs, createDev.Id)
		devStatuses = append(devStatuses, createDev.Status)
		devTSes = append(devTSes, createDev.CreatedAt.AsTime())
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.ID,
			time.Time{}, "", 0)
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, dev := range listDevs {
			if dev.Id == devIDs[len(devIDs)-1] &&
				dev.Status == devStatuses[len(devIDs)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.ID,
			devTSes[0], devIDs[0], 5)
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, dev := range listDevs {
			if dev.Id == devIDs[len(devIDs)-1] &&
				dev.Status == devStatuses[len(devIDs)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.ID,
			time.Time{}, "", 1)
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 1)
		require.Equal(t, int32(3), listCount)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, uuid.NewString(),
			time.Time{}, "", 0)
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("List devices by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, random.String(10),
			time.Time{}, "", 0)
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, int32(0), listCount)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
