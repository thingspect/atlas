package service

//go:generate mockgen -source device.go -destination mock_devicer_test.go -package service

import (
	"context"
	"time"

	"github.com/mennanov/fmutils"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/lora"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Devicer defines the methods provided by a device.DAO.
type Devicer interface {
	Create(ctx context.Context, dev *common.Device) (*common.Device, error)
	Read(ctx context.Context, devID, orgID string) (*common.Device, error)
	Update(ctx context.Context, dev *common.Device) (*common.Device, error)
	Delete(ctx context.Context, devID, orgID string) error
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32, tag string) ([]*common.Device, int32, error)
}

// Device service contains functions to query and modify devices.
type Device struct {
	api.UnimplementedDeviceServiceServer

	devDAO Devicer
	lora   lora.Loraer
}

// NewDevice instantiates and returns a new Device service.
func NewDevice(devDAO Devicer, lora lora.Loraer) *Device {
	return &Device{
		devDAO: devDAO,
		lora:   lora,
	}
}

// CreateDevice creates a device.
func (d *Device) CreateDevice(ctx context.Context,
	req *api.CreateDeviceRequest) (*common.Device, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	req.Device.OrgId = sess.OrgID

	dev, err := d.devDAO.Create(ctx, req.Device)
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		"201")); err != nil {
		logger.Errorf("CreateDevice grpc.SetHeader: %v", err)
	}

	return dev, nil
}

// CreateDeviceLoRaWAN adds LoRaWAN configuration to a device.
func (d *Device) CreateDeviceLoRaWAN(ctx context.Context,
	req *api.CreateDeviceLoRaWANRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	dev, err := d.devDAO.Read(ctx, req.Id, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	switch v := req.TypeOneof.(type) {
	case *api.CreateDeviceLoRaWANRequest_GatewayLorawanType:
		err = d.lora.CreateGateway(ctx, dev.UniqId)
	case *api.CreateDeviceLoRaWANRequest_DeviceLorawanType:
		err = d.lora.CreateDevice(ctx, dev.UniqId, v.DeviceLorawanType.AppKey)
	}
	if err != nil {
		logger.Errorf("CreateDeviceLoRaWAN d.lora.CreateX: %v", err)

		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		"204")); err != nil {
		logger.Errorf("CreateDeviceLoRaWAN grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetDevice retrieves a device by ID.
func (d *Device) GetDevice(ctx context.Context,
	req *api.GetDeviceRequest) (*common.Device, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_VIEWER {
		return nil, errPerm(common.Role_VIEWER)
	}

	dev, err := d.devDAO.Read(ctx, req.Id, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return dev, nil
}

// UpdateDevice updates a device. Update actions validate after merge to support
// partial updates.
func (d *Device) UpdateDevice(ctx context.Context,
	req *api.UpdateDeviceRequest) (*common.Device, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	if req.Device == nil {
		return nil, status.Error(codes.InvalidArgument,
			req.Validate().Error())
	}
	req.Device.OrgId = sess.OrgID

	// Perform partial update if directed.
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		// Normalize and validate field mask.
		req.UpdateMask.Normalize()
		if !req.UpdateMask.IsValid(req.Device) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		dev, err := d.devDAO.Read(ctx, req.Device.Id, sess.OrgID)
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.Device, req.UpdateMask.Paths)
		if req.Device.Tags != nil {
			dev.Tags = nil
		}
		proto.Merge(dev, req.Device)
		req.Device = dev
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dev, err := d.devDAO.Update(ctx, req.Device)
	if err != nil {
		return nil, errToStatus(err)
	}

	return dev, nil
}

// DeleteDeviceLoRaWAN removes LoRaWAN configuration from a device.
func (d *Device) DeleteDeviceLoRaWAN(ctx context.Context,
	req *api.DeleteDeviceLoRaWANRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	dev, err := d.devDAO.Read(ctx, req.Id, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	// Delete any gateways and devices present. 'Unauthenticated' is currently
	// returned for gateways and devices that do not exist.
	err = d.lora.DeleteGateway(ctx, dev.UniqId)
	if code := status.Code(err); code != codes.OK &&
		code != codes.Unauthenticated {
		logger.Errorf("DeleteDeviceLoRaWAN d.lora.DeleteGateway: %v", err)

		return nil, errToStatus(err)
	}

	err = d.lora.DeleteDevice(ctx, dev.UniqId)
	if code := status.Code(err); code != codes.OK &&
		code != codes.Unauthenticated {
		logger.Errorf("DeleteDeviceLoRaWAN d.lora.DeleteDevice: %v", err)

		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		"204")); err != nil {
		logger.Errorf("DeleteDeviceLoRaWAN grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeleteDevice deletes a device by ID.
func (d *Device) DeleteDevice(ctx context.Context,
	req *api.DeleteDeviceRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	if err := d.devDAO.Delete(ctx, req.Id, sess.OrgID); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		"204")); err != nil {
		logger.Errorf("DeleteDevice grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListDevices retrieves all devices.
func (d *Device) ListDevices(ctx context.Context,
	req *api.ListDevicesRequest) (*api.ListDevicesResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_VIEWER {
		return nil, errPerm(common.Role_VIEWER)
	}

	if req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}

	lBoundTS, prevID, err := session.ParsePageToken(req.PageToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	devs, count, err := d.devDAO.List(ctx, sess.OrgID, lBoundTS, prevID,
		req.PageSize+1, req.Tag)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListDevicesResponse{Devices: devs, TotalSize: count}

	// Populate next page token.
	if len(devs) == int(req.PageSize+1) {
		resp.Devices = devs[:len(devs)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			devs[len(devs)-2].CreatedAt.AsTime(),
			devs[len(devs)-2].Id); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger.Errorf("ListDevices session.GeneratePageToken dev, err: "+
				"%+v, %v", devs[len(devs)-2], err)
		}
	}

	return resp, nil
}
