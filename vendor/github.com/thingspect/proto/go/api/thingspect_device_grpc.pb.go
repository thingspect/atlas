// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.24.4
// source: api/thingspect_device.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DeviceService_CreateDevice_FullMethodName        = "/thingspect.api.DeviceService/CreateDevice"
	DeviceService_CreateDeviceLoRaWAN_FullMethodName = "/thingspect.api.DeviceService/CreateDeviceLoRaWAN"
	DeviceService_GetDevice_FullMethodName           = "/thingspect.api.DeviceService/GetDevice"
	DeviceService_UpdateDevice_FullMethodName        = "/thingspect.api.DeviceService/UpdateDevice"
	DeviceService_DeleteDeviceLoRaWAN_FullMethodName = "/thingspect.api.DeviceService/DeleteDeviceLoRaWAN"
	DeviceService_DeleteDevice_FullMethodName        = "/thingspect.api.DeviceService/DeleteDevice"
	DeviceService_ListDevices_FullMethodName         = "/thingspect.api.DeviceService/ListDevices"
)

// DeviceServiceClient is the client API for DeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// DeviceService contains functions to query and modify devices.
type DeviceServiceClient interface {
	// Create a device. Devices generate data points.
	CreateDevice(ctx context.Context, in *CreateDeviceRequest, opts ...grpc.CallOption) (*Device, error)
	// Add LoRaWAN configuration to a device.
	CreateDeviceLoRaWAN(ctx context.Context, in *CreateDeviceLoRaWANRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Get a device by ID. Devices generate data points.
	GetDevice(ctx context.Context, in *GetDeviceRequest, opts ...grpc.CallOption) (*Device, error)
	// Update a device. Devices generate data points.
	UpdateDevice(ctx context.Context, in *UpdateDeviceRequest, opts ...grpc.CallOption) (*Device, error)
	// Remove LoRaWAN configuration from a device.
	DeleteDeviceLoRaWAN(ctx context.Context, in *DeleteDeviceLoRaWANRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete a device by ID. Devices generate data points.
	DeleteDevice(ctx context.Context, in *DeleteDeviceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List all devices. Devices generate data points.
	ListDevices(ctx context.Context, in *ListDevicesRequest, opts ...grpc.CallOption) (*ListDevicesResponse, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) CreateDevice(ctx context.Context, in *CreateDeviceRequest, opts ...grpc.CallOption) (*Device, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Device)
	err := c.cc.Invoke(ctx, DeviceService_CreateDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) CreateDeviceLoRaWAN(ctx context.Context, in *CreateDeviceLoRaWANRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DeviceService_CreateDeviceLoRaWAN_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) GetDevice(ctx context.Context, in *GetDeviceRequest, opts ...grpc.CallOption) (*Device, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Device)
	err := c.cc.Invoke(ctx, DeviceService_GetDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) UpdateDevice(ctx context.Context, in *UpdateDeviceRequest, opts ...grpc.CallOption) (*Device, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Device)
	err := c.cc.Invoke(ctx, DeviceService_UpdateDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) DeleteDeviceLoRaWAN(ctx context.Context, in *DeleteDeviceLoRaWANRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DeviceService_DeleteDeviceLoRaWAN_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) DeleteDevice(ctx context.Context, in *DeleteDeviceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DeviceService_DeleteDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) ListDevices(ctx context.Context, in *ListDevicesRequest, opts ...grpc.CallOption) (*ListDevicesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDevicesResponse)
	err := c.cc.Invoke(ctx, DeviceService_ListDevices_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServiceServer is the server API for DeviceService service.
// All implementations must embed UnimplementedDeviceServiceServer
// for forward compatibility.
//
// DeviceService contains functions to query and modify devices.
type DeviceServiceServer interface {
	// Create a device. Devices generate data points.
	CreateDevice(context.Context, *CreateDeviceRequest) (*Device, error)
	// Add LoRaWAN configuration to a device.
	CreateDeviceLoRaWAN(context.Context, *CreateDeviceLoRaWANRequest) (*emptypb.Empty, error)
	// Get a device by ID. Devices generate data points.
	GetDevice(context.Context, *GetDeviceRequest) (*Device, error)
	// Update a device. Devices generate data points.
	UpdateDevice(context.Context, *UpdateDeviceRequest) (*Device, error)
	// Remove LoRaWAN configuration from a device.
	DeleteDeviceLoRaWAN(context.Context, *DeleteDeviceLoRaWANRequest) (*emptypb.Empty, error)
	// Delete a device by ID. Devices generate data points.
	DeleteDevice(context.Context, *DeleteDeviceRequest) (*emptypb.Empty, error)
	// List all devices. Devices generate data points.
	ListDevices(context.Context, *ListDevicesRequest) (*ListDevicesResponse, error)
	mustEmbedUnimplementedDeviceServiceServer()
}

// UnimplementedDeviceServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDeviceServiceServer struct{}

func (UnimplementedDeviceServiceServer) CreateDevice(context.Context, *CreateDeviceRequest) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDevice not implemented")
}
func (UnimplementedDeviceServiceServer) CreateDeviceLoRaWAN(context.Context, *CreateDeviceLoRaWANRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDeviceLoRaWAN not implemented")
}
func (UnimplementedDeviceServiceServer) GetDevice(context.Context, *GetDeviceRequest) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDevice not implemented")
}
func (UnimplementedDeviceServiceServer) UpdateDevice(context.Context, *UpdateDeviceRequest) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDevice not implemented")
}
func (UnimplementedDeviceServiceServer) DeleteDeviceLoRaWAN(context.Context, *DeleteDeviceLoRaWANRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDeviceLoRaWAN not implemented")
}
func (UnimplementedDeviceServiceServer) DeleteDevice(context.Context, *DeleteDeviceRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDevice not implemented")
}
func (UnimplementedDeviceServiceServer) ListDevices(context.Context, *ListDevicesRequest) (*ListDevicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDevices not implemented")
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}
func (UnimplementedDeviceServiceServer) testEmbeddedByValue()                       {}

// UnsafeDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServiceServer will
// result in compilation errors.
type UnsafeDeviceServiceServer interface {
	mustEmbedUnimplementedDeviceServiceServer()
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	// If the following call pancis, it indicates UnimplementedDeviceServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DeviceService_ServiceDesc, srv)
}

func _DeviceService_CreateDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).CreateDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_CreateDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).CreateDevice(ctx, req.(*CreateDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_CreateDeviceLoRaWAN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDeviceLoRaWANRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).CreateDeviceLoRaWAN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_CreateDeviceLoRaWAN_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).CreateDeviceLoRaWAN(ctx, req.(*CreateDeviceLoRaWANRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_GetDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).GetDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_GetDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).GetDevice(ctx, req.(*GetDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_UpdateDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).UpdateDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_UpdateDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).UpdateDevice(ctx, req.(*UpdateDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_DeleteDeviceLoRaWAN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDeviceLoRaWANRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).DeleteDeviceLoRaWAN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_DeleteDeviceLoRaWAN_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).DeleteDeviceLoRaWAN(ctx, req.(*DeleteDeviceLoRaWANRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_DeleteDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).DeleteDevice(ctx, req.(*DeleteDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_ListDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDevicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).ListDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceService_ListDevices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).ListDevices(ctx, req.(*ListDevicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DeviceService_ServiceDesc is the grpc.ServiceDesc for DeviceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "thingspect.api.DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDevice",
			Handler:    _DeviceService_CreateDevice_Handler,
		},
		{
			MethodName: "CreateDeviceLoRaWAN",
			Handler:    _DeviceService_CreateDeviceLoRaWAN_Handler,
		},
		{
			MethodName: "GetDevice",
			Handler:    _DeviceService_GetDevice_Handler,
		},
		{
			MethodName: "UpdateDevice",
			Handler:    _DeviceService_UpdateDevice_Handler,
		},
		{
			MethodName: "DeleteDeviceLoRaWAN",
			Handler:    _DeviceService_DeleteDeviceLoRaWAN_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _DeviceService_DeleteDevice_Handler,
		},
		{
			MethodName: "ListDevices",
			Handler:    _DeviceService_ListDevices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/thingspect_device.proto",
}
