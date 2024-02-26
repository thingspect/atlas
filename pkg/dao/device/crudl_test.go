//go:build !unit

package device

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 8 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("dao-device", createOrg.GetId())
		createDev, _ := proto.Clone(dev).(*api.Device)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, createDev)
		t.Logf("dev, createDev, err: %+v, %+v, %v", dev, createDev, err)
		require.NoError(t, err)
		require.NotEqual(t, dev.GetId(), createDev.GetId())
		require.NotEqual(t, dev.GetToken(), createDev.GetToken())
		require.WithinDuration(t, time.Now(), createDev.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with uppercase UniqId", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("dao-device", createOrg.GetId())
		dev.UniqId = strings.ToUpper(dev.GetUniqId())
		createDev, _ := proto.Clone(dev).(*api.Device)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, createDev)
		t.Logf("dev, createDev, err: %+v, %+v, %v", dev, createDev, err)
		require.NoError(t, err)
		require.NotEqual(t, dev.GetId(), createDev.GetId())
		require.Equal(t, strings.ToLower(dev.GetUniqId()), createDev.GetUniqId())
		require.NotEqual(t, dev.GetToken(), createDev.GetToken())
		require.WithinDuration(t, time.Now(), createDev.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("dao-device", createOrg.GetId())
		dev.UniqId = "dao-device-" + random.String(40)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("dev, createDev, err: %+v, %+v, %v", dev, createDev, err)
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
		createOrg.GetId()))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.Read(ctx, createDev.GetId(), createDev.GetOrgId())
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

		readDev, err := globalDevDAO.Read(ctx, createDev.GetId(),
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
			createDev.GetOrgId())
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
		createOrg.GetId()))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.ReadByUniqID(ctx, createDev.GetUniqId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDev)
	})

	t.Run("Read device by unknown UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAO.ReadByUniqID(ctx, random.String(16))
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestReadUpdateDeleteCache(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-device"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Read cached device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAOCache.Create(ctx,
			random.Device("dao-device", createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		readDev, err := globalDevDAOCache.Read(ctx, createDev.GetId(),
			createDev.GetOrgId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, readDev)
		}

		readDev, err = globalDevDAOCache.Read(ctx, createDev.GetId(),
			createDev.GetOrgId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, readDev)
		}
	})

	t.Run("Read cached device by valid UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAOCache.Create(ctx,
			random.Device("dao-device", createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		readDev, err := globalDevDAOCache.ReadByUniqID(ctx, createDev.GetUniqId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, readDev)
		}

		readDev, err = globalDevDAOCache.ReadByUniqID(ctx, createDev.GetUniqId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, readDev)
		}
	})

	t.Run("Read updated device by valid ID and UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAOCache.Create(ctx,
			random.Device("dao-device", createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "dao-device-" + random.String(16)
		createDev.Name = "dao-device-" + random.String(10)
		createDev.Status = api.Status_DISABLED
		createDev.Decoder = api.Decoder_GATEWAY
		createDev.Tags = nil
		updateDev, _ := proto.Clone(createDev).(*api.Device)

		updateDev, err = globalDevDAOCache.Update(ctx, updateDev)
		t.Logf("createDev, updateDev, err: %+v, %+v, %v", createDev, updateDev,
			err)
		require.NoError(t, err)
		require.Equal(t, createDev.GetUniqId(), updateDev.GetUniqId())
		require.Equal(t, createDev.GetName(), updateDev.GetName())
		require.Equal(t, createDev.GetStatus(), updateDev.GetStatus())
		require.Equal(t, createDev.GetDecoder(), updateDev.GetDecoder())
		require.Equal(t, createDev.GetTags(), updateDev.GetTags())
		require.True(t, updateDev.GetUpdatedAt().AsTime().After(
			updateDev.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createDev.GetCreatedAt().AsTime(),
			updateDev.GetUpdatedAt().AsTime(), 2*time.Second)

		readDev, err := globalDevDAOCache.Read(ctx, createDev.GetId(),
			createDev.GetOrgId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateDev, readDev)
		}

		readDev, err = globalDevDAOCache.ReadByUniqID(ctx, createDev.GetUniqId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateDev, readDev)
		}
	})

	t.Run("Read delete device by valid ID and UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAOCache.Create(ctx,
			random.Device("dao-device", createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		readDev, err := globalDevDAOCache.Read(ctx, createDev.GetId(),
			createDev.GetOrgId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, readDev)
		}

		readDev, err = globalDevDAOCache.ReadByUniqID(ctx, createDev.GetUniqId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, readDev)
		}

		err = globalDevDAOCache.Delete(ctx, createDev.GetId(), createOrg.GetId())
		t.Logf("err: %v", err)
		require.NoError(t, err)

		readDev, err = globalDevDAOCache.Read(ctx, createDev.GetId(), createOrg.GetId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)

		readDev, err = globalDevDAOCache.ReadByUniqID(ctx, createDev.GetUniqId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAOCache.Read(ctx, uuid.NewString(),
			uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAOCache.Create(ctx,
			random.Device("dao-device", createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		readDev, err := globalDevDAOCache.Read(ctx, createDev.GetId(),
			createDev.GetOrgId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, readDev)
		}

		readDev, err = globalDevDAOCache.Read(ctx, createDev.GetId(),
			uuid.NewString())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read device by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAOCache.Read(ctx, random.String(10),
			createOrg.GetId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})

	t.Run("Read device by unknown UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readDev, err := globalDevDAOCache.ReadByUniqID(ctx, random.String(16))
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
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "dao-device-" + random.String(16)
		createDev.Name = "dao-device-" + random.String(10)
		createDev.Status = api.Status_DISABLED
		createDev.Decoder = api.Decoder_GATEWAY
		createDev.Tags = nil
		updateDev, _ := proto.Clone(createDev).(*api.Device)

		updateDev, err = globalDevDAO.Update(ctx, updateDev)
		t.Logf("createDev, updateDev, err: %+v, %+v, %v", createDev, updateDev,
			err)
		require.NoError(t, err)
		require.Equal(t, createDev.GetUniqId(), updateDev.GetUniqId())
		require.Equal(t, createDev.GetName(), updateDev.GetName())
		require.Equal(t, createDev.GetStatus(), updateDev.GetStatus())
		require.Equal(t, createDev.GetDecoder(), updateDev.GetDecoder())
		require.Equal(t, createDev.GetTags(), updateDev.GetTags())
		require.True(t, updateDev.GetUpdatedAt().AsTime().After(
			updateDev.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createDev.GetCreatedAt().AsTime(),
			updateDev.GetUpdatedAt().AsTime(), 2*time.Second)

		readDev, err := globalDevDAO.Read(ctx, createDev.GetId(), createDev.GetOrgId())
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)
		require.Equal(t, updateDev, readDev)
	})

	t.Run("Update unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateDev, err := globalDevDAO.Update(ctx, random.Device("dao-device",
			createOrg.GetId()))
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.OrgId = uuid.NewString()
		createDev.UniqId = "dao-device-" + random.String(16)

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("createDev, updateDev, err: %+v, %+v, %v", createDev, updateDev,
			err)
		require.Nil(t, updateDev)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "dao-device-" + random.String(40)

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("createDev, updateDev, err: %+v, %+v, %v", createDev, updateDev,
			err)
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
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDevDAO.Delete(ctx, createDev.GetId(), createOrg.GetId())
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readDev, err := globalDevDAO.Read(ctx, createDev.GetId(),
				createOrg.GetId())
			t.Logf("readDev, err: %+v, %v", readDev, err)
			require.Nil(t, readDev)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalDevDAO.Delete(ctx, uuid.NewString(), createOrg.GetId())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDevDAO.Delete(ctx, createDev.GetId(), uuid.NewString())
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
	devNames := []string{}
	devStatuses := []api.Status{}
	devDecoders := []api.Decoder{}
	devTags := [][]string{}
	devTSes := []time.Time{}
	for range 3 {
		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-device",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		devIDs = append(devIDs, createDev.GetId())
		devNames = append(devNames, createDev.GetName())
		devStatuses = append(devStatuses, createDev.GetStatus())
		devDecoders = append(devDecoders, createDev.GetDecoder())
		devTags = append(devTags, createDev.GetTags())
		devTSes = append(devTSes, createDev.GetCreatedAt().AsTime())
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.GetId(),
			time.Time{}, "", 0, "")
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, dev := range listDevs {
			if dev.GetId() == devIDs[len(devIDs)-1] &&
				dev.GetName() == devNames[len(devNames)-1] &&
				dev.GetStatus() == devStatuses[len(devStatuses)-1] &&
				dev.GetDecoder() == devDecoders[len(devDecoders)-1] &&
				reflect.DeepEqual(dev.GetTags(), devTags[len(devTags)-1]) {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.GetId(),
			devTSes[0], devIDs[0], 5, "")
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, dev := range listDevs {
			if dev.GetId() == devIDs[len(devIDs)-1] &&
				dev.GetName() == devNames[len(devNames)-1] &&
				dev.GetStatus() == devStatuses[len(devStatuses)-1] &&
				dev.GetDecoder() == devDecoders[len(devDecoders)-1] &&
				reflect.DeepEqual(dev.GetTags(), devTags[len(devTags)-1]) {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.GetId(),
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

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.GetId(),
			time.Time{}, "", 5, devTags[2][0])
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, devIDs[len(devIDs)-1], listDevs[0].GetId())
		require.Equal(t, devNames[len(devNames)-1], listDevs[0].GetName())
		require.Equal(t, devStatuses[len(devStatuses)-1], listDevs[0].GetStatus())
		require.Equal(t, devDecoders[len(devDecoders)-1], listDevs[0].GetDecoder())
		require.Equal(t, devTags[len(devTags)-1], listDevs[0].GetTags())
	})

	t.Run("List devices with tag filter and pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listDevs, listCount, err := globalDevDAO.List(ctx, createOrg.GetId(),
			devTSes[0], devIDs[0], 5, devTags[2][0])
		t.Logf("listDevs, listCount, err: %+v, %v, %v", listDevs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listDevs, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, devIDs[len(devIDs)-1], listDevs[0].GetId())
		require.Equal(t, devNames[len(devNames)-1], listDevs[0].GetName())
		require.Equal(t, devStatuses[len(devStatuses)-1], listDevs[0].GetStatus())
		require.Equal(t, devDecoders[len(devDecoders)-1], listDevs[0].GetDecoder())
		require.Equal(t, devTags[len(devTags)-1], listDevs[0].GetTags())
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
		require.Empty(t, listDevs)
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
