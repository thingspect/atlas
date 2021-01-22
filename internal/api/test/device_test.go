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
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateDevice(t *testing.T) {
	t.Parallel()

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotNil(t, createDev)
		require.Equal(t, globalAuthOrgID, createDev.OrgId)
		require.Equal(t, dev.UniqId, createDev.UniqId)
		require.Equal(t, dev.Status, createDev.Status)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with uppercase UniqId", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{UniqId: strings.ToUpper("api-device-" +
			random.String(16)), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotNil(t, createDev)
		require.Equal(t, globalAuthOrgID, createDev.OrgId)
		require.Equal(t, strings.ToLower(dev.UniqId), createDev.UniqId)
		require.Equal(t, dev.Status, createDev.Status)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{UniqId: "api-device-" + random.String(40),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateDeviceRequest.Device: embedded message failed "+
			"validation | caused by: invalid Device.UniqId: value length must "+
			"be between 5 and 40 runes, inclusive")
	})
}

func TestGetDevice(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dev := &api.Device{UniqId: "api-device-" + random.String(16),
		Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}

	devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
	createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
		Device: dev})
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Get device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		getDev, err := devCli.GetDevice(ctx, &api.GetDeviceRequest{
			Id: createDev.Id})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createDev, getDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createDev, getDev)
		}
	})

	t.Run("Get device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		getDev, err := devCli.GetDevice(ctx, &api.GetDeviceRequest{
			Id: uuid.NewString()})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.Nil(t, getDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Get are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		secCli := api.NewDeviceServiceClient(secondaryAuthGRPCConn)
		getDev, err := secCli.GetDevice(ctx, &api.GetDeviceRequest{
			Id: createDev.Id})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.Nil(t, getDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestUpdateDevice(t *testing.T) {
	t.Parallel()

	t.Run("Update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "api-device-" + random.String(16)
		createDev.Status = []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]

		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev.UniqId, updateDev.UniqId)
		require.Equal(t, createDev.Status, updateDev.Status)
		require.Equal(t, createDev.CreatedAt.AsTime(),
			updateDev.CreatedAt.AsTime())
		require.True(t, updateDev.UpdatedAt.AsTime().After(
			updateDev.CreatedAt.AsTime()))
		require.WithinDuration(t, createDev.CreatedAt.AsTime(),
			updateDev.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Partial update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		part := &api.Device{Id: createDev.Id, UniqId: "api-device-" +
			random.String(16), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}

		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"uniq_id", "status"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, part.UniqId, updateDev.UniqId)
		require.Equal(t, part.Status, updateDev.Status)
		require.Equal(t, createDev.CreatedAt.AsTime(),
			updateDev.CreatedAt.AsTime())
		require.True(t, updateDev.UpdatedAt.AsTime().After(
			updateDev.CreatedAt.AsTime()))
		require.WithinDuration(t, createDev.CreatedAt.AsTime(),
			updateDev.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Update nil device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: nil})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateDeviceRequest.Device: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownDevice := &api.Device{Id: uuid.NewString(),
			UniqId: "api-device-" + random.String(16),
			Token:  uuid.NewString()}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: unknownDevice, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update device by unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownDevice := &api.Device{Id: uuid.NewString(),
			UniqId: "api-device-" + random.String(16),
			Token:  uuid.NewString()}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: unknownDevice, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"uniq_id", "token"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update device by unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownDevice := &api.Device{Id: uuid.NewString(),
			UniqId: "api-device-" + random.String(16), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)],
			Token: uuid.NewString()}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
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

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.OrgId = uuid.NewString()
		createDev.UniqId = "api-device-" + random.String(16)

		secCli := api.NewDeviceServiceClient(secondaryAuthGRPCConn)
		updateDev, err := secCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update device validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "api-device-" + random.String(40)

		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateDeviceRequest.Device: embedded message failed "+
			"validation | caused by: invalid Device.UniqId: value length must "+
			"be between 5 and 40 runes, inclusive")
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.Token = random.String(10)

		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: UUID")
	})
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: createDev.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
			getDev, err := devCli.GetDevice(ctx, &api.GetDeviceRequest{
				Id: createDev.Id})
			t.Logf("getDev, err: %+v, %v", getDev, err)
			require.Nil(t, getDev)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		_, err := devCli.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		secCli := api.NewDeviceServiceClient(secondaryAuthGRPCConn)
		_, err = secCli.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: createDev.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListDevices(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	devIDs := []string{}
	devStatuses := []common.Status{}
	for i := 0; i < 3; i++ {
		dev := &api.Device{UniqId: "api-device-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		devIDs = append(devIDs, createDev.Id)
		devStatuses = append(devStatuses, createDev.Status)
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listDevs.Devices), 3)
		require.GreaterOrEqual(t, listDevs.TotalSize, int32(3))

		var found bool
		for _, dev := range listDevs.Devices {
			if dev.Id == devIDs[len(devIDs)-1] &&
				dev.Status == devStatuses[len(devIDs)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{PageSize: 2})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.Devices, 2)
		require.Empty(t, listDevs.PrevPageToken)
		require.NotEmpty(t, listDevs.NextPageToken)
		require.GreaterOrEqual(t, listDevs.TotalSize, int32(3))

		nextDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{PageSize: 2,
			PageToken: listDevs.NextPageToken})
		t.Logf("nextDevs, err: %+v, %v", nextDevs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(nextDevs.Devices), 1)
		require.NotEmpty(t, nextDevs.PrevPageToken)
		require.GreaterOrEqual(t, nextDevs.TotalSize, int32(3))
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		secCli := api.NewDeviceServiceClient(secondaryAuthGRPCConn)
		listDevs, err := secCli.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.Devices, 0)
		require.Equal(t, int32(0), listDevs.TotalSize)
	})

	t.Run("List devices by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{
			PageToken: "..."})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}
