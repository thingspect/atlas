// +build !unit

package device

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16),
			IsDisabled: []bool{true, false}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgId, createDev.OrgId)
		require.Equal(t, dev.UniqId, createDev.UniqId)
		require.Equal(t, dev.IsDisabled, createDev.IsDisabled)
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
			UniqId:     strings.ToUpper(random.String(16)),
			IsDisabled: []bool{true, false}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgId, createDev.OrgId)
		require.Equal(t, strings.ToLower(dev.UniqId), createDev.UniqId)
		require.Equal(t, dev.IsDisabled, createDev.IsDisabled)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(41)}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.EqualError(t, err, "ERROR: value too long for type character "+
			"varying(40) (SQLSTATE 22001)")
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16),
		IsDisabled: []bool{true, false}[random.Intn(2)]}
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDevDAO.Read(ctx, createDev.Id, createDev.OrgId)
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDevice)
	})

	t.Run("Read device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDevDAO.Read(ctx, uuid.New().String(),
			uuid.New().String())
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.Nil(t, readDevice)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDevDAO.Read(ctx, createDev.Id,
			uuid.New().String())
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.Nil(t, readDevice)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Read device by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDevDAO.Read(ctx, random.String(10),
			createDev.OrgId)
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.Nil(t, readDevice)
		require.Contains(t, err.Error(),
			"ERROR: invalid input syntax for type uuid")
	})
}

func TestReadByUniqID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16),
		IsDisabled: []bool{true, false}[random.Intn(2)]}
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDevDAO.ReadByUniqID(ctx, createDev.UniqId)
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDevice)
	})

	t.Run("Read device by unknown UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDevDAO.ReadByUniqID(ctx, uuid.New().String())
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.Nil(t, readDevice)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func TestUpdateDevice(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16)}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = random.String(16)
		createDev.IsDisabled = []bool{true, false}[random.Intn(2)]

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev.UniqId, updateDev.UniqId)
		require.Equal(t, createDev.IsDisabled, updateDev.IsDisabled)
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

		unknownDevice := &api.Device{Id: uuid.New().String(),
			OrgId: createOrg.ID, UniqId: random.String(16),
			Token: uuid.New().String()}
		updateDev, err := globalDevDAO.Update(ctx, unknownDevice)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16)}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.OrgId = uuid.New().String()
		createDev.UniqId = random.String(16)

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16)}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = random.String(41)

		updateDev, err := globalDevDAO.Update(ctx, createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "ERROR: value too long for type character "+
			"varying(40) (SQLSTATE 22001)")
	})
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16)}
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

			readDevice, err := globalDevDAO.Read(ctx, createDev.Id,
				createOrg.ID)
			t.Logf("readDevice, err: %+v, %v", readDevice, err)
			require.Nil(t, readDevice)
			require.Equal(t, sql.ErrNoRows, err)
		})
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalDevDAO.Delete(ctx, uuid.New().String(), createOrg.ID)
		t.Logf("err: %v", err)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16)}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDevDAO.Delete(ctx, createDev.Id, uuid.New().String())
		t.Logf("err: %v", err)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func TestListDevices(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	var lastDeviceID string
	var lastDeviceDisabled bool
	for i := 0; i < 3; i++ {
		dev := &api.Device{OrgId: createOrg.ID, UniqId: random.String(16),
			IsDisabled: []bool{true, false}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		lastDeviceID = createDev.Id
		lastDeviceDisabled = createDev.IsDisabled
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listDevs, err := globalDevDAO.List(ctx, createOrg.ID)
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 3)

		var found bool
		for _, dev := range listDevs {
			if dev.Id == lastDeviceID && dev.IsDisabled == lastDeviceDisabled {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listDevs, err := globalDevDAO.List(ctx, uuid.New().String())
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 0)
	})
}
