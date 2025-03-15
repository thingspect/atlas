//go:build !unit

package test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateDevice(t *testing.T) {
	t.Parallel()

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx,
			&api.CreateDeviceRequest{Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotEqual(t, dev.GetId(), createDev.GetId())
		require.Equal(t, dev.GetUniqId(), createDev.GetUniqId())
		require.NotEqual(t, dev.GetToken(), createDev.GetToken())
		require.WithinDuration(t, time.Now(), createDev.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with uppercase UniqId", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		dev.UniqId = strings.ToUpper(dev.GetUniqId())

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminKeyGRPCConn)
		createDev, err := devCli.CreateDevice(ctx,
			&api.CreateDeviceRequest{Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)
		require.NotEqual(t, dev.GetId(), createDev.GetId())
		require.Equal(t, strings.ToLower(dev.GetUniqId()), createDev.GetUniqId())
		require.NotEqual(t, dev.GetToken(), createDev.GetToken())
		require.WithinDuration(t, time.Now(), createDev.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createDev.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		dev.UniqId = "api-device-" + random.String(40)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx,
			&api.CreateDeviceRequest{Device: dev})
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id:        createDev.GetId(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create valid device configuration", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminKeyGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id: createDev.GetId(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_DeviceLorawanType{
					DeviceLorawanType: &api.CreateDeviceLoRaWANRequest_DeviceLoRaWANType{
						AppKey: random.String(32),
					},
				},
			})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create configuration with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		_, err := devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id:        uuid.NewString(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create configuration by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		_, err := devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id:        uuid.NewString(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Configurations are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id:        createDev.GetId(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestGetDevice(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
	createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
		Device: random.Device("api-device", uuid.NewString()),
	})
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	t.Run("Get device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		getDev, err := devCli.GetDevice(ctx,
			&api.GetDeviceRequest{Id: createDev.GetId()})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, createDev, getDev)
	})

	t.Run("Get device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		getDev, err := devCli.GetDevice(ctx,
			&api.GetDeviceRequest{Id: uuid.NewString()})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.Nil(t, getDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Gets are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		getDev, err := secCli.GetDevice(ctx,
			&api.GetDeviceRequest{Id: createDev.GetId()})
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "api-device-" + random.String(16)
		createDev.Name = "api-device-" + random.String(10)
		createDev.Status = api.Status_DISABLED
		createDev.Decoder = api.Decoder_GATEWAY
		createDev.Tags = random.Tags("api-device", 2)

		updateDev, err := devCli.UpdateDevice(ctx,
			&api.UpdateDeviceRequest{Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
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

		getDev, err := devCli.GetDevice(ctx,
			&api.GetDeviceRequest{Id: createDev.GetId()})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateDev, getDev)
	})

	t.Run("Partial update device by valid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminKeyGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		part := &api.Device{
			Id: createDev.GetId(), UniqId: "api-device-" + random.String(16),
			Name:   "api-device-" + random.String(10),
			Status: api.Status_DISABLED, Decoder: api.Decoder_GATEWAY,
			Tags: random.Tags("api-device", 2),
		}

		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: part, UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{
				"uniq_id", "name", "status", "decoder", "tags",
			}},
		})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)
		require.Equal(t, part.GetUniqId(), updateDev.GetUniqId())
		require.Equal(t, part.GetName(), updateDev.GetName())
		require.Equal(t, part.GetStatus(), updateDev.GetStatus())
		require.Equal(t, part.GetDecoder(), updateDev.GetDecoder())
		require.Equal(t, part.GetTags(), updateDev.GetTags())
		require.True(t, updateDev.GetUpdatedAt().AsTime().After(
			updateDev.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createDev.GetCreatedAt().AsTime(),
			updateDev.GetUpdatedAt().AsTime(), 2*time.Second)

		getDev, err := devCli.GetDevice(ctx,
			&api.GetDeviceRequest{Id: createDev.GetId()})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateDev, getDev)
	})

	t.Run("Update device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx,
			&api.UpdateDeviceRequest{Device: nil})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateDeviceRequest.Device: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device:     random.Device("api-device", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"aaa"}},
		})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update device by unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"uniq_id", "token"},
			},
		})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update device by unknown device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		updateDev, err := devCli.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.OrgId = uuid.NewString()
		createDev.UniqId = "api-device-" + random.String(16)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		updateDev, err := secCli.UpdateDevice(ctx,
			&api.UpdateDeviceRequest{Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update device validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.UniqId = "api-device-" + random.String(40)

		updateDev, err := devCli.UpdateDevice(ctx,
			&api.UpdateDeviceRequest{Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateDeviceRequest.Device: embedded message failed "+
			"validation | caused by: invalid Device.UniqId: value length must "+
			"be between 5 and 40 runes, inclusive")
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Update device fields.
		createDev.Token = random.String(10)

		updateDev, err := devCli.UpdateDevice(ctx,
			&api.UpdateDeviceRequest{Device: createDev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateDeviceRequest.Device: embedded message failed "+
			"validation | caused by: invalid Device.Token: value must be a "+
			"valid UUID | caused by: invalid uuid format")
	})
}

func TestDeleteDeviceLoRaWAN(t *testing.T) {
	t.Parallel()

	t.Run("Delete valid configurations by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id:        createDev.GetId(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		_, err = devCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: createDev.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete unknown configurations by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: createDev.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete configurations with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: createDev.GetId()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		_, err = devCli.DeleteDevice(ctx,
			&api.DeleteDeviceRequest{Id: createDev.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read device by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()

			devCli := api.NewDeviceServiceClient(globalAdminKeyGRPCConn)
			getDev, err := devCli.GetDevice(ctx,
				&api.GetDeviceRequest{Id: createDev.GetId()})
			t.Logf("getDev, err: %+v, %v", getDev, err)
			require.Nil(t, getDev)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(secondaryViewerGRPCConn)
		_, err := devCli.DeleteDevice(ctx,
			&api.DeleteDeviceRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		_, err := devCli.DeleteDevice(ctx,
			&api.DeleteDeviceRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteDevice(ctx,
			&api.DeleteDeviceRequest{Id: createDev.GetId()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListDevices(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	devIDs := []string{}
	devNames := []string{}
	devStatuses := []api.Status{}
	devDecoders := []api.Decoder{}
	devTags := [][]string{}
	for range 3 {
		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-device", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		devIDs = append(devIDs, createDev.GetId())
		devNames = append(devNames, createDev.GetName())
		devStatuses = append(devStatuses, createDev.GetStatus())
		devDecoders = append(devDecoders, createDev.GetDecoder())
		devTags = append(devTags, createDev.GetTags())
	}

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		listDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listDevs.GetDevices()), 3)
		require.GreaterOrEqual(t, listDevs.GetTotalSize(), int32(3))

		var found bool
		for _, dev := range listDevs.GetDevices() {
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

	t.Run("List devices by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminKeyGRPCConn)
		listDevs, err := devCli.ListDevices(ctx,
			&api.ListDevicesRequest{PageSize: 2})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.GetDevices(), 2)
		require.NotEmpty(t, listDevs.GetNextPageToken())
		require.GreaterOrEqual(t, listDevs.GetTotalSize(), int32(3))

		nextDevs, err := devCli.ListDevices(ctx, &api.ListDevicesRequest{
			PageSize: 2, PageToken: listDevs.GetNextPageToken(),
		})
		t.Logf("nextDevs, err: %+v, %v", nextDevs, err)
		require.NoError(t, err)
		require.NotEmpty(t, nextDevs.GetDevices())
		require.GreaterOrEqual(t, nextDevs.GetTotalSize(), int32(3))
	})

	t.Run("List devices with tag filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		listDevs, err := devCli.ListDevices(ctx,
			&api.ListDevicesRequest{Tag: devTags[2][0]})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Len(t, listDevs.GetDevices(), 1)
		require.Equal(t, int32(1), listDevs.GetTotalSize())

		require.Equal(t, devIDs[len(devIDs)-1], listDevs.GetDevices()[0].GetId())
		require.Equal(t, devNames[len(devNames)-1], listDevs.GetDevices()[0].GetName())
		require.Equal(t, devStatuses[len(devStatuses)-1],
			listDevs.GetDevices()[0].GetStatus())
		require.Equal(t, devDecoders[len(devDecoders)-1],
			listDevs.GetDevices()[0].GetDecoder())
		require.Equal(t, devTags[len(devTags)-1], listDevs.GetDevices()[0].GetTags())
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewDeviceServiceClient(secondaryAdminGRPCConn)
		listDevs, err := secCli.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Empty(t, listDevs.GetDevices())
		require.Equal(t, int32(0), listDevs.GetTotalSize())
	})

	t.Run("List devices by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		listDevs, err := devCli.ListDevices(ctx,
			&api.ListDevicesRequest{PageToken: badUUID})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}
