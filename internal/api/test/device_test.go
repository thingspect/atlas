// +build !unit

package test

import (
	"context"
	"reflect"
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

		dev := random.Device("api-device", uuid.NewString())

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotEqual(t, dev.Id, createDev.Id)
		require.Equal(t, dev.UniqId, createDev.UniqId)
		require.NotEqual(t, dev.Token, createDev.Token)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with uppercase UniqId", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		dev.UniqId = strings.ToUpper(dev.UniqId)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotEqual(t, dev.Id, createDev.Id)
		require.Equal(t, strings.ToLower(dev.UniqId), createDev.UniqId)
		require.NotEqual(t, dev.Token, createDev.Token)
		require.WithinDuration(t, time.Now(), createDev.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		dev.UniqId = "api-device-" + random.String(40)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
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

func TestCreateDeviceLoRaWAN(t *testing.T) {
	t.Parallel()

	t.Run("Create valid gateway configuration", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{Id: createDev.Id,
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create valid device configuration", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{Id: createDev.Id,
				TypeOneof: &api.CreateDeviceLoRaWANRequest_DeviceLorawanType{
					DeviceLorawanType: &api.CreateDeviceLoRaWANRequest_DeviceLoRaWANType{
						AppKey: random.String(32)}}})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create configuration with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		_, err := devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{Id: uuid.NewString(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create configuration by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		_, err := devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{Id: uuid.NewString(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Configurations are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{Id: createDev.Id,
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestGetDevice(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
	createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
		Device: random.Device("api-device", uuid.NewString())})
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Get device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		getDev, err := devCli.GetDevice(ctx, &api.GetDeviceRequest{
			Id: uuid.NewString()})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.Nil(t, getDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Get are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "api-device-" + random.String(16)
		createDev.Name = "api-device-" + random.String(10)
		createDev.Status = common.Status_DISABLED
		createDev.Decoder = common.Decoder_GATEWAY
		createDev.Tags = random.Tags("api-device-", 2)

		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, createDev.UniqId, updateDev.UniqId)
		require.Equal(t, createDev.CreatedAt.AsTime(),
			updateDev.CreatedAt.AsTime())
		require.True(t, updateDev.UpdatedAt.AsTime().After(
			updateDev.CreatedAt.AsTime()))
		require.WithinDuration(t, createDev.CreatedAt.AsTime(),
			updateDev.UpdatedAt.AsTime(), 2*time.Second)

		getDev, err := devCli.GetDevice(ctx, &api.GetDeviceRequest{
			Id: createDev.Id})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateDev, getDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateDev, getDev)
		}
	})

	t.Run("Partial update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		part := &common.Device{Id: createDev.Id, UniqId: "api-device-" +
			random.String(16), Name: "api-device-" + random.String(10),
			Status: common.Status_DISABLED, Decoder: common.Decoder_GATEWAY}

		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"uniq_id", "name", "status", "decoder"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, part.UniqId, updateDev.UniqId)
		require.Equal(t, createDev.CreatedAt.AsTime(),
			updateDev.CreatedAt.AsTime())
		require.True(t, updateDev.UpdatedAt.AsTime().After(
			updateDev.CreatedAt.AsTime()))
		require.WithinDuration(t, createDev.CreatedAt.AsTime(),
			updateDev.UpdatedAt.AsTime(), 2*time.Second)

		getDev, err := devCli.GetDevice(ctx, &api.GetDeviceRequest{
			Id: createDev.Id})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateDev, getDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateDev, getDev)
		}
	})

	t.Run("Update device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Update nil device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: nil})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateDeviceRequest.Device: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device:     random.Device("api-device", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"aaa"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update device by unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"uniq_id", "token"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update device by unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.OrgId = uuid.NewString()
		createDev.UniqId = "api-device-" + random.String(16)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		updateDev, err := secCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update device validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
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

func TestDeleteDeviceLoRaWAN(t *testing.T) {
	t.Parallel()

	t.Run("Delete valid configurations by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{Id: createDev.Id,
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		_, err = devCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: createDev.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete unknown configurations by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: createDev.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete configurations with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		_, err := devCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete configurations by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		_, err := devCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Configurations are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: createDev.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: createDev.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
			getDev, err := devCli.GetDevice(ctx, &api.GetDeviceRequest{
				Id: createDev.Id})
			t.Logf("getDev, err: %+v, %v", getDev, err)
			require.Nil(t, getDev)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		_, err := devCli.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		_, err := devCli.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: createDev.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListDevices(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	devIDs := []string{}
	devNames := []string{}
	devStatuses := []common.Status{}
	devDecoders := []common.Decoder{}
	devTags := [][]string{}
	for i := 0; i < 3; i++ {
		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		devIDs = append(devIDs, createDev.Id)
		devNames = append(devNames, createDev.Name)
		devStatuses = append(devStatuses, createDev.Status)
		devDecoders = append(devDecoders, createDev.Decoder)
		devTags = append(devTags, createDev.Tags)
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listDevs.Devices), 3)
		require.GreaterOrEqual(t, listDevs.TotalSize, int32(3))

		var found bool
		for _, dev := range listDevs.Devices {
			if dev.Id == devIDs[len(devIDs)-1] &&
				dev.Name == devNames[len(devNames)-1] &&
				dev.Status == devStatuses[len(devStatuses)-1] &&
				dev.Decoder == devDecoders[len(devDecoders)-1] &&
				reflect.DeepEqual(dev.Tags, devTags[len(devTags)-1]) {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List devices by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{
			PageSize: 2})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.Devices, 2)
		require.NotEmpty(t, listDevs.NextPageToken)
		require.GreaterOrEqual(t, listDevs.TotalSize, int32(3))

		nextDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{
			PageSize: 2, PageToken: listDevs.NextPageToken})
		t.Logf("nextDevs, err: %+v, %v", nextDevs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(nextDevs.Devices), 1)
		require.GreaterOrEqual(t, nextDevs.TotalSize, int32(3))
	})

	t.Run("List devices with tag filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{
			Tag: devTags[2][0]})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.Devices, 1)
		require.Equal(t, int32(1), listDevs.TotalSize)

		require.Equal(t, devIDs[len(devIDs)-1], listDevs.Devices[0].Id)
		require.Equal(t, devNames[len(devNames)-1], listDevs.Devices[0].Name)
		require.Equal(t, devStatuses[len(devStatuses)-1],
			listDevs.Devices[0].Status)
		require.Equal(t, devDecoders[len(devDecoders)-1],
			listDevs.Devices[0].Decoder)
		require.Equal(t, devTags[len(devTags)-1], listDevs.Devices[0].Tags)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		listDevs, err := secCli.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.Devices, 0)
		require.Equal(t, int32(0), listDevs.TotalSize)
	})

	t.Run("List devices by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{
			PageToken: badUUID})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}
