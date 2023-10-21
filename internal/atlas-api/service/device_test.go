//go:build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/atlas-api/lora"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateDevice(t *testing.T) {
	t.Parallel()

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		retDev, _ := proto.Clone(dev).(*api.Device)

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Create(gomock.Any(), dev).Return(retDev, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		createDev, err := devSvc.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev,
		})
		t.Logf("dev, createDev, err: %+v, %+v, %v", dev, createDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(dev, createDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", dev, createDev)
		}
	})

	t.Run("Create device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		createDev, err := devSvc.CreateDevice(ctx, &api.CreateDeviceRequest{})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		createDev, err := devSvc.CreateDevice(ctx, &api.CreateDeviceRequest{})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		dev.UniqId = random.String(41)

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Create(gomock.Any(), dev).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		createDev, err := devSvc.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev,
		})
		t.Logf("dev, createDev, err: %+v, %+v, %v", dev, createDev, err)
		require.Nil(t, createDev)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestCreateDeviceLoRaWAN(t *testing.T) {
	t.Parallel()

	t.Run("Create valid gateway configuration", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		ctrl := gomock.NewController(t)
		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(dev, nil).
			Times(1)
		loraer := lora.NewMockLoraer(ctrl)
		loraer.EXPECT().CreateGateway(gomock.Any(), dev.GetUniqId()).Return(nil).
			Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, loraer)
		_, err := devSvc.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id:        dev.GetId(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_GatewayLorawanType{},
			})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create valid device configuration", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		appKey := random.String(32)

		ctrl := gomock.NewController(t)
		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(dev, nil).
			Times(1)
		loraer := lora.NewMockLoraer(ctrl)
		loraer.EXPECT().CreateDevice(gomock.Any(), dev.GetUniqId(), appKey).
			Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, loraer)
		_, err := devSvc.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id: dev.GetId(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_DeviceLorawanType{
					DeviceLorawanType: &api.CreateDeviceLoRaWANRequest_DeviceLoRaWANType{
						AppKey: appKey,
					},
				},
			})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create configuration with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		_, err := devSvc.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create configuration with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		_, err := devSvc.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create configuration by unknown ID", func(t *testing.T) {
		t.Parallel()

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		_, err := devSvc.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Create invalid configuration", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		appKey := random.String(33)

		ctrl := gomock.NewController(t)
		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(dev, nil).
			Times(1)
		loraer := lora.NewMockLoraer(ctrl)
		loraer.EXPECT().CreateDevice(gomock.Any(), dev.GetUniqId(), appKey).
			Return(dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, loraer)
		_, err := devSvc.CreateDeviceLoRaWAN(ctx,
			&api.CreateDeviceLoRaWANRequest{
				Id: dev.GetId(),
				TypeOneof: &api.CreateDeviceLoRaWANRequest_DeviceLorawanType{
					DeviceLorawanType: &api.CreateDeviceLoRaWANRequest_DeviceLoRaWANType{
						AppKey: appKey,
					},
				},
			})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestGetDevice(t *testing.T) {
	t.Parallel()

	t.Run("Get device by valid ID", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		retDev, _ := proto.Clone(dev).(*api.Device)

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(retDev,
			nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		getDev, err := devSvc.GetDevice(ctx, &api.GetDeviceRequest{Id: dev.GetId()})
		t.Logf("dev, getDev, err: %+v, %+v, %v", dev, getDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(dev, getDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", dev, getDev)
		}
	})

	t.Run("Get device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		getDev, err := devSvc.GetDevice(ctx, &api.GetDeviceRequest{})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.Nil(t, getDev)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Get device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		getDev, err := devSvc.GetDevice(ctx, &api.GetDeviceRequest{})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.Nil(t, getDev)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Get device by unknown ID", func(t *testing.T) {
		t.Parallel()

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		getDev, err := devSvc.GetDevice(ctx, &api.GetDeviceRequest{
			Id: uuid.NewString(),
		})
		t.Logf("getDev, err: %+v, %v", getDev, err)
		require.Nil(t, getDev)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestUpdateDevice(t *testing.T) {
	t.Parallel()

	t.Run("Update device by valid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		retDev, _ := proto.Clone(dev).(*api.Device)

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Update(gomock.Any(), dev).Return(retDev, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: dev,
		})
		t.Logf("dev, updateDev, err: %+v, %+v, %v", dev, updateDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(dev, updateDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", dev, updateDev)
		}
	})

	t.Run("Partial update device by valid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		retDev, _ := proto.Clone(dev).(*api.Device)
		part := &api.Device{
			Id: dev.GetId(), Status: api.Status_ACTIVE,
			Decoder: api.Decoder_GATEWAY, Tags: random.Tags("api-device", 2),
		}
		merged := &api.Device{
			Id: dev.GetId(), OrgId: dev.GetOrgId(), UniqId: dev.GetUniqId(), Name: dev.GetName(),
			Status: part.GetStatus(), Token: dev.GetToken(), Decoder: part.GetDecoder(),
			Tags: part.GetTags(),
		}
		retMerged, _ := proto.Clone(merged).(*api.Device)

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(retDev,
			nil).Times(1)
		devicer.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(merged)).
			Return(retMerged, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status", "decoder", "tags"},
			},
		})
		t.Logf("merged, updateDev, err: %+v, %+v, %v", merged, updateDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(merged, updateDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", merged, updateDev)
		}
	})

	t.Run("Update device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Update device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Update nil device", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: nil,
		})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateDeviceRequest.Device: value is required"), err)
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: dev, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"},
			},
		})
		t.Logf("dev, updateDev, err: %+v, %+v, %v", dev, updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid field mask"), err)
	})

	t.Run("Partial update device by unknown device", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		part := &api.Device{Id: uuid.NewString(), Status: api.Status_ACTIVE}

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Read(gomock.Any(), part.GetId(), orgID).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"},
			},
		})
		t.Logf("part, updateDev, err: %+v, %+v, %v", part, updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Update device validation failure", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())
		dev.UniqId = random.String(41)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: dev,
		})
		t.Logf("dev, updateDev, err: %+v, %+v, %v", dev, updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateDeviceRequest.Device: embedded message failed "+
				"validation | caused by: invalid Device.UniqId: value length "+
				"must be between 5 and 40 runes, inclusive"), err)
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Update(gomock.Any(), dev).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		updateDev, err := devSvc.UpdateDevice(ctx, &api.UpdateDeviceRequest{
			Device: dev,
		})
		t.Logf("dev, updateDev, err: %+v, %+v, %v", dev, updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestDeleteDeviceLoRaWAN(t *testing.T) {
	t.Parallel()

	t.Run("Delete valid configurations by valid ID", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		ctrl := gomock.NewController(t)
		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(dev, nil).
			Times(1)
		loraer := lora.NewMockLoraer(ctrl)
		loraer.EXPECT().DeleteGateway(gomock.Any(), dev.GetUniqId()).Return(nil).
			Times(1)
		loraer.EXPECT().DeleteDevice(gomock.Any(), dev.GetUniqId()).Return(nil).
			Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, loraer)
		_, err := devSvc.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: dev.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete unknown configurations by valid ID", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		ctrl := gomock.NewController(t)
		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(dev, nil).
			Times(1)
		loraer := lora.NewMockLoraer(ctrl)
		loraer.EXPECT().DeleteGateway(gomock.Any(), dev.GetUniqId()).
			Return(status.Error(codes.Unauthenticated,
				"authentication failed: not authorized")).Times(1)
		loraer.EXPECT().DeleteDevice(gomock.Any(), dev.GetUniqId()).
			Return(status.Error(codes.Unauthenticated,
				"authentication failed: not authorized")).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, loraer)
		_, err := devSvc.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: dev.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete configurations with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		_, err := devSvc.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete configurations with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		_, err := devSvc.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete configurations by unknown ID", func(t *testing.T) {
		t.Parallel()

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		_, err := devSvc.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Delete invalid gateway configuration", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		ctrl := gomock.NewController(t)
		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(dev, nil).
			Times(1)
		loraer := lora.NewMockLoraer(ctrl)
		loraer.EXPECT().DeleteGateway(gomock.Any(), dev.GetUniqId()).
			Return(dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, loraer)
		_, err := devSvc.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: dev.GetId()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("Delete invalid device configuration", func(t *testing.T) {
		t.Parallel()

		dev := random.Device("api-device", uuid.NewString())

		ctrl := gomock.NewController(t)
		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.GetId(), dev.GetOrgId()).Return(dev, nil).
			Times(1)
		loraer := lora.NewMockLoraer(ctrl)
		loraer.EXPECT().DeleteGateway(gomock.Any(), dev.GetUniqId()).Return(nil).
			Times(1)
		loraer.EXPECT().DeleteDevice(gomock.Any(), dev.GetUniqId()).
			Return(dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: dev.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, loraer)
		_, err := devSvc.DeleteDeviceLoRaWAN(ctx,
			&api.DeleteDeviceLoRaWANRequest{Id: dev.GetId()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		_, err := devSvc.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		_, err := devSvc.DeleteDevice(ctx, &api.DeleteDeviceRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete device with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		_, err := devSvc.DeleteDevice(ctx, &api.DeleteDeviceRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		_, err := devSvc.DeleteDevice(ctx, &api.DeleteDeviceRequest{
			Id: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestListDevices(t *testing.T) {
	t.Parallel()

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		devs := []*api.Device{
			random.Device("api-device", uuid.NewString()),
			random.Device("api-device", uuid.NewString()),
			random.Device("api-device", uuid.NewString()),
		}

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(51),
			"").Return(devs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		listDevs, err := devSvc.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listDevs.GetTotalSize())

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDevicesResponse{Devices: devs, TotalSize: 3},
			listDevs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListDevicesResponse{Devices: devs, TotalSize: 3}, listDevs)
		}
	})

	t.Run("List devices by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		devs := []*api.Device{
			random.Device("api-device", uuid.NewString()),
			random.Device("api-device", uuid.NewString()),
			random.Device("api-device", uuid.NewString()),
		}

		next, err := session.GeneratePageToken(devs[1].GetCreatedAt().AsTime(),
			devs[1].GetId())
		require.NoError(t, err)

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3),
			"").Return(devs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		listDevs, err := devSvc.ListDevices(ctx, &api.ListDevicesRequest{
			PageSize: 2,
		})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listDevs.GetTotalSize())

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDevicesResponse{
			Devices: devs[:2], NextPageToken: next, TotalSize: 3,
		}, listDevs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListDevicesResponse{
				Devices: devs[:2], NextPageToken: next, TotalSize: 3,
			}, listDevs)
		}
	})

	t.Run("List devices with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		listDevs, err := devSvc.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List devices with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		listDevs, err := devSvc.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List devices by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(nil, nil)
		listDevs, err := devSvc.ListDevices(ctx, &api.ListDevicesRequest{
			PageToken: badUUID,
		})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid page token"), err)
	})

	t.Run("List devices by invalid org ID", func(t *testing.T) {
		t.Parallel()

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any()).Return(nil, int32(0),
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		listDevs, err := devSvc.ListDevices(ctx, &api.ListDevicesRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("List devices with generation failure", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		devs := []*api.Device{
			random.Device("api-device", uuid.NewString()),
			random.Device("api-device", uuid.NewString()),
			random.Device("api-device", uuid.NewString()),
		}
		devs[1].Id = badUUID

		devicer := NewMockDevicer(gomock.NewController(t))
		devicer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3),
			"").Return(devs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		devSvc := NewDevice(devicer, nil)
		listDevs, err := devSvc.ListDevices(ctx, &api.ListDevicesRequest{
			PageSize: 2,
		})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listDevs.GetTotalSize())

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDevicesResponse{
			Devices: devs[:2], TotalSize: 3,
		}, listDevs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListDevicesResponse{
				Devices: devs[:2], TotalSize: 3,
			}, listDevs)
		}
	})
}
