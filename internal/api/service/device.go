package service

//go:generate mockgen -source device.go -destination mock_devicer_test.go -package service

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mennanov/fmutils"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Devicer defines the methods provided by a device.DAO.
type Devicer interface {
	Create(ctx context.Context, dev *api.Device) (*api.Device, error)
	Read(ctx context.Context, devID, orgID string) (*api.Device, error)
	Update(ctx context.Context, dev *api.Device) (*api.Device, error)
	Delete(ctx context.Context, devID, orgID string) error
	List(ctx context.Context, orgID string, lboundTS time.Time, prevID string,
		limit int32) ([]*api.Device, int32, error)
}

// Device service contains functions to query and modify devices.
type Device struct {
	api.UnimplementedDeviceServiceServer

	devDAO Devicer
}

// NewDevice instantiates and returns a new Device service.
func NewDevice(devDAO Devicer) *Device {
	return &Device{
		devDAO: devDAO,
	}
}

// Create creates a device.
func (d *Device) Create(ctx context.Context,
	req *api.CreateDeviceRequest) (*api.CreateDeviceResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	req.Device.OrgId = sess.OrgID

	dev, err := d.devDAO.Create(ctx, req.Device)
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"201")); err != nil {
		logger.Errorf("Device.Create grpc.SetHeader: %v", err)
	}
	return &api.CreateDeviceResponse{Device: dev}, nil
}

// Read retrieves a device by ID.
func (d *Device) Read(ctx context.Context,
	req *api.ReadDeviceRequest) (*api.ReadDeviceResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	dev, err := d.devDAO.Read(ctx, req.Id, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.ReadDeviceResponse{Device: dev}, nil
}

// Update updates a device.
func (d *Device) Update(ctx context.Context,
	req *api.UpdateDeviceRequest) (*api.UpdateDeviceResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.Device == nil {
		return nil, status.Error(codes.InvalidArgument, req.Validate().Error())
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

	return &api.UpdateDeviceResponse{Device: dev}, nil
}

// Delete deletes a device by ID.
func (d *Device) Delete(ctx context.Context,
	req *api.DeleteDeviceRequest) (*empty.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if err := d.devDAO.Delete(ctx, req.Id, sess.OrgID); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"204")); err != nil {
		logger.Errorf("Device.Delete grpc.SetHeader: %v", err)
	}
	return &empty.Empty{}, nil
}

// List retrieves all devices.
func (d *Device) List(ctx context.Context,
	req *api.ListDeviceRequest) (*api.ListDeviceResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}

	lboundTS, prevID, err := session.ParsePageToken(req.PageToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	devs, count, err := d.devDAO.List(ctx, sess.OrgID, lboundTS, prevID,
		req.PageSize+1)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListDeviceResponse{
		Devices:       devs,
		PrevPageToken: req.PageToken,
		TotalSize:     count,
	}

	// Populate next page token.
	if len(devs) == int(req.PageSize+1) {
		resp.Devices = devs[:len(devs)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			devs[len(devs)-2].CreatedAt.AsTime(),
			devs[len(devs)-2].Id); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger.Errorf("Device.List session.GeneratePageToken dev, err: "+
				"%+v, %v", devs[len(devs)-2], err)
		}
	}

	return resp, nil
}
