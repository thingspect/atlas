// +build !unit

package device

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dev := random.Device("dao-device", createOrg.Id)
		dev.Tags = nil
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgId, createDev.OrgId)
		require.Equal(t, dev.UniqId, createDev.UniqId)
		require.Equal(t, dev.Name, createDev.Name)
		require.Equal(t, dev.Status, createDev.Status)
		require.Equal(t, dev.Decoder, createDev.Decoder)
		require.Equal(t, dev.Tags, createDev.Tags)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with tags", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dev := random.Device("dao-device", createOrg.Id)
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgId, createDev.OrgId)
		require.Equal(t, dev.UniqId, createDev.UniqId)
		require.Equal(t, dev.Name, createDev.Name)
		require.Equal(t, dev.Status, createDev.Status)
		require.Equal(t, dev.Decoder, createDev.Decoder)
		require.Equal(t, dev.Tags, createDev.Tags)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with uppercase UniqId", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dev := random.Device("dao-device", createOrg.Id)
		dev.UniqId = strings.ToUpper(dev.UniqId)
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgId, createDev.OrgId)
		require.Equal(t, strings.ToLower(dev.UniqId), createDev.UniqId)
		require.Equal(t, dev.Name, createDev.Name)
		require.Equal(t, dev.Status, createDev.Status)
		require.Equal(t, dev.Decoder, createDev.Decoder)
		require.Equal(t, dev.Tags, createDev.Tags)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dev := random.Device("dao-device", createOrg.Id)
		dev.UniqId = "dao-device-" + random.String(40)
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
		createOrg.Id))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, createDev.Id, createDev.OrgId)
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDev)
	})

	t.Run("Read device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, uuid.NewString(),
			uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, createDev.Id,
			uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read device by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
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

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
		createOrg.Id))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.ReadByUniqID(ctx, createDev.UniqId)
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDev)
	})

	t.Run("Read device by unknown UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.ReadByUniqID(ctx, uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.Id))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "dao-device-" + random.String(16)
		createDev.Name = "dao-device-" + random.String(10)
		createDev.Status = api.Status_DISABLED
		createDev.Decoder = api.Decoder_GATEWAY
		createDev.Tags = nil

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev.UniqId, updateDev.UniqId)
		require.Equal(t, createDev.Name, updateDev.Name)
		require.Equal(t, createDev.Status, updateDev.Status)
		require.Equal(t, createDev.Decoder, updateDev.Decoder)
		require.Equal(t, createDev.Tags, updateDev.Tags)
		require.Equal(t, createDev.CreatedAt, updateDev.CreatedAt)
		require.True(t, updateDev.UpdatedAt.AsTime().After(
			updateDev.CreatedAt.AsTime()))
		require.WithinDuration(t, createDev.CreatedAt.AsTime(),
			updateDev.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Update unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateDev, err := globalDevDAO.Update(ctx, random.Device("dao-device",
			createOrg.Id))
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.Id))
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.Id))
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

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.Id))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDevDAO.Delete(ctx, createDev.Id, createOrg.Id)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readDev, err := globalDevDAO.Read(ctx, createDev.Id,
				createOrg.Id)
			t.Logf("readDev, err: %+v, %v", readDev, err)
			require.Nil(t, readDev)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalDevDAO.Delete(ctx, uuid.NewString(), createOrg.Id)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.Id))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDevDAO.Delete(ctx, createDev.Id, uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	devIDs := []string{}
	devStatuses := []api.Status{}
	devDecoders := []api.Decoder{}
	devTags := [][]string{}
	devTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.Id))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		devIDs = append(devIDs, createDev.Id)
		devStatuses = append(devStatuses, createDev.Status)
		devDecoders = append(devDecoders, createDev.Decoder)
		devTags = append(devTags, createDev.Tags)
		devTSes = append(devTSes, createDev.CreatedAt.AsTime())
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 0, "")
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, dev := range listDevs {
			if dev.Id == devIDs[len(devIDs)-1] &&
				dev.Status == devStatuses[len(devIDs)-1] &&
				dev.Decoder == devDecoders[len(devIDs)-1] &&
				reflect.DeepEqual(dev.Tags, devTags[len(devIDs)-1]) {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.Id,
			devTSes[0], devIDs[0], 5, "")
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, dev := range listDevs {
			if dev.Id == devIDs[len(devIDs)-1] &&
				dev.Status == devStatuses[len(devIDs)-1] &&
				dev.Decoder == devDecoders[len(devIDs)-1] &&
				reflect.DeepEqual(dev.Tags, devTags[len(devIDs)-1]) {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 1, "")
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 1)
		require.Equal(t, int32(3), listCount)
	})

	t.Run("List devices with tag filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 5, devTags[2][0])
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, devIDs[len(devIDs)-1], listDevs[0].Id)
		require.Equal(t, devStatuses[len(devIDs)-1], listDevs[0].Status)
		require.Equal(t, devDecoders[len(devIDs)-1], listDevs[0].Decoder)
		require.Equal(t, devTags[len(devIDs)-1], listDevs[0].Tags)
	})

	t.Run("List devices with tag filter and pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.Id,
			devTSes[0], devIDs[0], 5, devTags[2][0])
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, devIDs[len(devIDs)-1], listDevs[0].Id)
		require.Equal(t, devStatuses[len(devIDs)-1], listDevs[0].Status)
		require.Equal(t, devDecoders[len(devIDs)-1], listDevs[0].Decoder)
		require.Equal(t, devTags[len(devIDs)-1], listDevs[0].Tags)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, uuid.NewString(),
			time.Time{}, "", 0, "")
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("List devices by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, random.String(10),
			time.Time{}, "", 0, "")
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.Nil(t, listDevs)
		require.Equal(t, int32(0), listCount)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
