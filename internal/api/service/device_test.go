// +build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{OrgId: uuid.New().String(),
			UniqId: random.String(16), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Create(gomock.Any(), matcher.NewProtoMatcher(dev)).
			Return(dev, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: dev.OrgId}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		createDev, err := devSvc.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.CreateDeviceResponse{Device: dev}, createDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.CreateDeviceResponse{Device: dev}, createDev)
		}
	})

	t.Run("Create device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		createDev, err := devSvc.Create(ctx, &api.CreateDeviceRequest{
			Device: nil})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"), err)
	})

	t.Run("Create invalid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{OrgId: uuid.New().String(),
			UniqId: random.String(41), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Create(gomock.Any(), matcher.NewProtoMatcher(dev)).
			Return(nil, dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: dev.OrgId}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		createDev, err := devSvc.Create(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.Nil(t, createDev)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	t.Run("Read device by valid ID", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{Id: uuid.New().String(), OrgId: uuid.New().String(),
			UniqId: random.String(16), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.Id, dev.OrgId).Return(dev, nil).
			Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: dev.OrgId}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		readDev, err := devSvc.Read(ctx, &api.ReadDeviceRequest{Id: dev.Id})
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ReadDeviceResponse{Device: dev}, readDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ReadDeviceResponse{Device: dev}, readDev)
		}
	})

	t.Run("Read device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		readDev, err := devSvc.Read(ctx, &api.ReadDeviceRequest{
			Id: uuid.New().String()})
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"), err)
	})

	t.Run("Read device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.New().String()}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		readDev, err := devSvc.Read(ctx, &api.ReadDeviceRequest{
			Id: uuid.New().String()})
		t.Logf("readDev, err: %+v, %v", readDev, err)
		require.Nil(t, readDev)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("Update device by valid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{Id: uuid.New().String(), OrgId: uuid.New().String(),
			UniqId: random.String(16), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(dev)).
			Return(dev, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: dev.OrgId}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: dev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.UpdateDeviceResponse{Device: dev}, updateDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.UpdateDeviceResponse{Device: dev}, updateDev)
		}
	})

	t.Run("Partial update device by valid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{Id: uuid.New().String(), OrgId: uuid.New().String(),
			UniqId: random.String(16)}
		part := &api.Device{Id: dev.Id, Status: common.Status_ACTIVE}
		merged := &api.Device{Id: dev.Id, OrgId: dev.OrgId, UniqId: dev.UniqId,
			Status: common.Status_ACTIVE}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), dev.Id, dev.OrgId).Return(dev, nil).
			Times(1)
		devicer.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(merged)).
			Return(merged, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: dev.OrgId}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.UpdateDeviceResponse{Device: merged}, updateDev) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.UpdateDeviceResponse{Device: merged}, updateDev)
		}
	})

	t.Run("Update device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: nil})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"), err)
	})

	t.Run("Update nil device", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.New().String()}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: nil})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateDeviceRequest.Device: value is required"), err)
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{Id: uuid.New().String(), OrgId: uuid.New().String(),
			UniqId: random.String(16), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		devicer.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.New().String()}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: dev, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid field mask"), err)
	})

	t.Run("Partial update device by unknown device", func(t *testing.T) {
		t.Parallel()

		OrgID := uuid.New().String()
		part := &api.Device{Id: uuid.New().String(),
			Status: common.Status_ACTIVE}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Read(gomock.Any(), part.Id, OrgID).
			Return(nil, dao.ErrNotFound).Times(1)
		devicer.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: OrgID}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"}}})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Update device validation failure", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{Id: uuid.New().String(), OrgId: uuid.New().String(),
			UniqId: random.String(41), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: dev.OrgId}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: dev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateDeviceRequest.Device: embedded message failed "+
				"validation | caused by: invalid Device.UniqId: value length "+
				"must be between 5 and 40 runes, inclusive"), err)
	})

	t.Run("Update device by invalid device", func(t *testing.T) {
		t.Parallel()

		dev := &api.Device{Id: uuid.New().String(), OrgId: uuid.New().String(),
			UniqId: random.String(16), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)],
			Token: random.String(10)}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(dev)).
			Return(nil, dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: dev.OrgId}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		updateDev, err := devSvc.Update(ctx, &api.UpdateDeviceRequest{
			Device: dev})
		t.Logf("updateDev, err: %+v, %v", updateDev, err)
		require.Nil(t, updateDev)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("Delete device by valid ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.New().String()}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		_, err := devSvc.Delete(ctx, &api.DeleteDeviceRequest{
			Id: uuid.New().String()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete device with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		_, err := devSvc.Delete(ctx, &api.DeleteDeviceRequest{
			Id: uuid.New().String()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"), err)
	})

	t.Run("Delete device by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.New().String()}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		_, err := devSvc.Delete(ctx, &api.DeleteDeviceRequest{
			Id: uuid.New().String()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("List devices by valid org ID", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.New().String()

		devs := []*api.Device{
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)]},
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)]},
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)]},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(101)).
			Return(devs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		listDevs, err := devSvc.List(ctx, &api.ListDeviceRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listDevs.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDeviceResponse{Devices: devs, TotalSize: 3},
			listDevs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListDeviceResponse{Devices: devs, TotalSize: 3}, listDevs)
		}
	})

	t.Run("List devices by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.New().String()

		devs := []*api.Device{
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)],
				CreatedAt: timestamppb.Now()},
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)],
				CreatedAt: timestamppb.Now()},
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)],
				CreatedAt: timestamppb.Now()},
		}

		next, err := session.GeneratePageToken(devs[1].CreatedAt.AsTime(),
			devs[1].Id)
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(devs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		listDevs, err := devSvc.List(ctx, &api.ListDeviceRequest{PageSize: 2})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listDevs.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDeviceResponse{Devices: devs[:2],
			NextPageToken: next, TotalSize: 3}, listDevs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListDeviceResponse{Devices: devs[:2], NextPageToken: next,
					TotalSize: 3}, listDevs)
		}
	})

	t.Run("List devices with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		listDevs, err := devSvc.List(ctx, &api.ListDeviceRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"),
			err)
	})

	t.Run("List devices by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.New().String()}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		listDevs, err := devSvc.List(ctx, &api.ListDeviceRequest{
			PageToken: "..."})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid page token"), err)
	})

	t.Run("List devices by invalid orgID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil, int32(0), dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: "aaa"}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		listDevs, err := devSvc.List(ctx, &api.ListDeviceRequest{})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.Nil(t, listDevs)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("List devices with generation failure", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.New().String()

		devs := []*api.Device{
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)],
				CreatedAt: timestamppb.Now()},
			{Id: "...", OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)],
				CreatedAt: timestamppb.Now()},
			{Id: uuid.New().String(), OrgId: orgID, UniqId: random.String(16),
				Status: []common.Status{common.Status_ACTIVE,
					common.Status_DISABLED}[random.Intn(2)],
				CreatedAt: timestamppb.Now()},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		devicer := NewMockDevicer(ctrl)
		devicer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(devs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID}),
			2*time.Second)
		defer cancel()

		devSvc := NewDevice(devicer)
		listDevs, err := devSvc.List(ctx, &api.ListDeviceRequest{PageSize: 2})
		t.Logf("listDevs, err: %+v, %v", listDevs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listDevs.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDeviceResponse{Devices: devs[:2],
			TotalSize: 3}, listDevs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListDeviceResponse{Devices: devs[:2], TotalSize: 3},
				listDevs)
		}
	})
}
