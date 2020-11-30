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

		dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgID, createDev.OrgID)
		require.Equal(t, dev.UniqID, createDev.UniqID)
	})

	t.Run("Create valid device with uppercase UniqID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := Device{OrgID: createOrg.ID,
			UniqID: strings.ToUpper(random.String(10))}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.Equal(t, dev.OrgID, createDev.OrgID)
		require.Equal(t, strings.ToLower(dev.UniqID), createDev.UniqID)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := Device{OrgID: createOrg.ID, UniqID: random.String(41)}
		createDev, err := globalDAO.Create(ctx, dev)
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

	dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
	createDev, err := globalDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDAO.Read(ctx, createDev.ID, createDev.OrgID)
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.NoError(t, err)
		require.Equal(t, createDev, readDevice)
	})

	t.Run("Read device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDAO.Read(ctx, uuid.New().String(),
			uuid.New().String())
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.Nil(t, readDevice)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDAO.Read(ctx, createDev.ID,
			uuid.New().String())
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.Nil(t, readDevice)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Read device by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readDevice, err := globalDAO.Read(ctx, random.String(10),
			createDev.OrgID)
		t.Logf("readDevice, err: %+v, %v", readDevice, err)
		require.Nil(t, readDevice)
		require.Contains(t, err.Error(),
			"ERROR: invalid input syntax for type uuid")
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

		dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqID = random.String(10)

		updateDev, err := globalDAO.Update(ctx, *createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev.UniqID, updateDev.UniqID)
		require.Equal(t, createDev.CreatedAt, updateDev.CreatedAt)
		require.True(t, updateDev.UpdatedAt.After(updateDev.CreatedAt))
		require.WithinDuration(t, createDev.CreatedAt, updateDev.UpdatedAt,
			2*time.Second)
	})

	t.Run("Update unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownDevice := Device{ID: uuid.New().String(), OrgID: createOrg.ID,
			UniqID: random.String(10), Token: uuid.New().String()}
		updateDev, err := globalDAO.Update(ctx, unknownDevice)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.OrgID = uuid.New().String()
		createDev.UniqID = random.String(10)

		updateDev, err := globalDAO.Update(ctx, *createDev)
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqID = random.String(41)

		updateDev, err := globalDAO.Update(ctx, *createDev)
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

		dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDAO.Delete(ctx, createDev.ID, createOrg.ID)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			readDevice, err := globalDAO.Read(ctx, createDev.ID, createOrg.ID)
			t.Logf("readDevice, err: %+v, %v", readDevice, err)
			require.Equal(t, sql.ErrNoRows, err)
		})
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalDAO.Delete(ctx, uuid.New().String(), createOrg.ID)
		t.Logf("err: %v", err)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		err = globalDAO.Delete(ctx, createDev.ID, uuid.New().String())
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
	for i := 0; i < 3; i++ {
		dev := Device{OrgID: createOrg.ID, UniqID: random.String(10)}
		createDev, err := globalDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		lastDeviceID = createDev.ID
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listDevs, err := globalDAO.List(ctx, createOrg.ID)
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 3)

		var found bool
		for _, dev := range listDevs {
			if dev.ID == lastDeviceID {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listDevs, err := globalDAO.List(ctx, uuid.New().String())
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs, 0)
	})
}
