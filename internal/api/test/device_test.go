// +build !unit

package test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "api-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotNil(t, createDev.Device)
		require.Equal(t, dev.OrgId, createDev.Device.OrgId)
		require.Equal(t, dev.UniqId, createDev.Device.UniqId)
		require.Equal(t, dev.IsDisabled, createDev.Device.IsDisabled)
		require.WithinDuration(t, time.Now(),
			createDev.Device.CreatedAt.AsTime(), 2*time.Second)
		require.WithinDuration(t, time.Now(),
			createDev.Device.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Create valid device with uppercase UniqId", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{OrgId: createOrg.ID,
			UniqId:     strings.ToUpper("api-device-" + random.String(16)),
			IsDisabled: []bool{true, false}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotNil(t, createDev.Device)
		require.Equal(t, dev.OrgId, createDev.Device.OrgId)
		require.Equal(t, strings.ToLower(dev.UniqId), createDev.Device.UniqId)
		require.Equal(t, dev.IsDisabled, createDev.Device.IsDisabled)
		require.WithinDuration(t, time.Now(),
			createDev.Device.CreatedAt.AsTime(), 2*time.Second)
		require.WithinDuration(t, time.Now(),
			createDev.Device.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Create nil device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: nil})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"device must not be nil")
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(40), IsDisabled: []bool{true, false}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: value too long")
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	org := org.Org{Name: "api-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
		random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

	devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
	createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{Device: dev})
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Read device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		readDev, err := devCli.Read(ctx, &api.ReadDeviceRequest{
			Id: createDev.Device.Id, OrgId: createDev.Device.OrgId})
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ReadDeviceResponse{Device: createDev.Device},
			readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ReadDeviceResponse{Device: createDev.Device}, readDev)
		}
	})

	t.Run("Read device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		readDev, err := devCli.Read(ctx, &api.ReadDeviceRequest{
			Id: uuid.New().String(), OrgId: createOrg.ID})
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		readDev, err := devCli.Read(ctx, &api.ReadDeviceRequest{
			Id: createDev.Device.Id, OrgId: uuid.New().String()})
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Read device by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		readDev, err := devCli.Read(ctx, &api.ReadDeviceRequest{
			Id: random.String(10), OrgId: createOrg.ID})
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: UUID")
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "api-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.Device.UniqId = "api-device-" + random.String(16)
		createDev.Device.IsDisabled = []bool{true, false}[random.Intn(2)]

		updateDev, err := devCli.Update(ctx, &api.UpdateDeviceRequest{
			Device: createDev.Device})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev.Device.UniqId, updateDev.Device.UniqId)
		require.Equal(t, createDev.Device.IsDisabled,
			updateDev.Device.IsDisabled)
		require.Equal(t, createDev.Device.CreatedAt.AsTime(),
			updateDev.Device.CreatedAt.AsTime())
		require.True(t, updateDev.Device.UpdatedAt.AsTime().After(
			updateDev.Device.CreatedAt.AsTime()))
		require.WithinDuration(t, createDev.Device.CreatedAt.AsTime(),
			updateDev.Device.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Update nil device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		updateDev, err := devCli.Update(ctx, &api.UpdateDeviceRequest{
			Device: nil})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"device must not be nil")
	})

	t.Run("Update unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownDevice := &api.Device{Id: uuid.New().String(),
			OrgId: createOrg.ID, UniqId: "api-device-" + random.String(16),
			Token: uuid.New().String()}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		updateDev, err := devCli.Update(ctx, &api.UpdateDeviceRequest{
			Device: unknownDevice})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.Device.OrgId = uuid.New().String()
		createDev.Device.UniqId = "api-device-" + random.String(16)

		updateDev, err := devCli.Update(ctx, &api.UpdateDeviceRequest{
			Device: createDev.Device})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.Device.UniqId = "api-device-" + random.String(40)

		updateDev, err := devCli.Update(ctx, &api.UpdateDeviceRequest{
			Device: createDev.Device})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: value too long")
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "api-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.Delete(ctx, &api.DeleteDeviceRequest{
			Id: createDev.Device.Id, OrgId: createDev.Device.OrgId})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
			readDev, err := devCli.Read(ctx, &api.ReadDeviceRequest{
				Id: createDev.Device.Id, OrgId: createDev.Device.OrgId})
			t.Logf("readDev, err: %+v, %v", readDev, err)
			require.Nil(t, readDev)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		_, err = devCli.Delete(ctx, &api.DeleteDeviceRequest{
			Id: uuid.New().String(), OrgId: createOrg.ID})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.Delete(ctx, &api.DeleteDeviceRequest{
			Id: createDev.Device.Id, OrgId: uuid.New().String()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	org := org.Org{Name: "api-device-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	var lastDeviceID string
	var lastDeviceDisabled bool
	for i := 0; i < 3; i++ {
		dev := &api.Device{OrgId: createOrg.ID, UniqId: "api-device-" +
			random.String(16), IsDisabled: []bool{true, false}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		lastDeviceID = createDev.Device.Id
		lastDeviceDisabled = createDev.Device.IsDisabled
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		listDevs, err := devCli.List(ctx, &api.ListDeviceRequest{
			OrgId: createOrg.ID})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.Devices, 3)

		var found bool
		for _, dev := range listDevs.Devices {
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

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		listDevs, err := devCli.List(ctx, &api.ListDeviceRequest{
			OrgId: uuid.New().String()})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.Devices, 0)
	})

	t.Run("Lists device by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		listDevs, err := devCli.List(ctx, &api.ListDeviceRequest{
			OrgId: random.String(10)})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: UUID")
	})
}
