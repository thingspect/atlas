package service

//go:generate mockgen -source device.go -destination mock_devicer_test.go -package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/thingspect/api/go/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Devicer defines the methods provided by a device.DAO.
type Devicer interface {
	Create(ctx context.Context, dev *api.Device) (*api.Device, error)
	Read(ctx context.Context, devID, orgID string) (*api.Device, error)
	Update(ctx context.Context, dev *api.Device) (*api.Device, error)
	Delete(ctx context.Context, devID, orgID string) error
	List(ctx context.Context, orgID string) ([]*api.Device, error)
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
	if req.Device == nil {
		return nil, status.Error(codes.InvalidArgument,
			"device must not be nil")
	}

	dev, err := d.devDAO.Create(ctx, req.Device)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.CreateDeviceResponse{Device: dev}, nil
}

// Read retrieves a device by ID.
func (d *Device) Read(ctx context.Context,
	req *api.ReadDeviceRequest) (*api.ReadDeviceResponse, error) {
	dev, err := d.devDAO.Read(ctx, req.Id, req.OrgId)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.ReadDeviceResponse{Device: dev}, nil
}

// Update updates a device.
func (d *Device) Update(ctx context.Context,
	req *api.UpdateDeviceRequest) (*api.UpdateDeviceResponse, error) {
	if req.Device == nil {
		return nil, status.Error(codes.InvalidArgument,
			"device must not be nil")
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
	if err := d.devDAO.Delete(ctx, req.Id, req.OrgId); err != nil {
		return nil, errToStatus(err)
	}

	return &empty.Empty{}, nil
}

// List retrieves all devices.
func (d *Device) List(ctx context.Context,
	req *api.ListDeviceRequest) (*api.ListDeviceResponse, error) {
	devs, err := d.devDAO.List(ctx, req.OrgId)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.ListDeviceResponse{Devices: devs}, nil
}
