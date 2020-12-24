// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// DeviceServiceClient is the client API for DeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceServiceClient interface {
	// Create creates a device.
	Create(ctx context.Context, in *CreateDeviceRequest, opts ...grpc.CallOption) (*CreateDeviceResponse, error)
	// Read retrieves a device by ID.
	Read(ctx context.Context, in *ReadDeviceRequest, opts ...grpc.CallOption) (*ReadDeviceResponse, error)
	// Update updates a device.
	Update(ctx context.Context, in *UpdateDeviceRequest, opts ...grpc.CallOption) (*UpdateDeviceResponse, error)
	// Delete deletes a device by ID.
	Delete(ctx context.Context, in *DeleteDeviceRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// List retrieves all devices.
	List(ctx context.Context, in *ListDeviceRequest, opts ...grpc.CallOption) (*ListDeviceResponse, error)
}

type deviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceServiceClient(cc grpc.ClientConnInterface) DeviceServiceClient {
	return &deviceServiceClient{cc}
}

func (c *deviceServiceClient) Create(ctx context.Context, in *CreateDeviceRequest, opts ...grpc.CallOption) (*CreateDeviceResponse, error) {
	out := new(CreateDeviceResponse)
	err := c.cc.Invoke(ctx, "/api.DeviceService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) Read(ctx context.Context, in *ReadDeviceRequest, opts ...grpc.CallOption) (*ReadDeviceResponse, error) {
	out := new(ReadDeviceResponse)
	err := c.cc.Invoke(ctx, "/api.DeviceService/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) Update(ctx context.Context, in *UpdateDeviceRequest, opts ...grpc.CallOption) (*UpdateDeviceResponse, error) {
	out := new(UpdateDeviceResponse)
	err := c.cc.Invoke(ctx, "/api.DeviceService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) Delete(ctx context.Context, in *DeleteDeviceRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/api.DeviceService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceServiceClient) List(ctx context.Context, in *ListDeviceRequest, opts ...grpc.CallOption) (*ListDeviceResponse, error) {
	out := new(ListDeviceResponse)
	err := c.cc.Invoke(ctx, "/api.DeviceService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceServiceServer is the server API for DeviceService service.
// All implementations must embed UnimplementedDeviceServiceServer
// for forward compatibility
type DeviceServiceServer interface {
	// Create creates a device.
	Create(context.Context, *CreateDeviceRequest) (*CreateDeviceResponse, error)
	// Read retrieves a device by ID.
	Read(context.Context, *ReadDeviceRequest) (*ReadDeviceResponse, error)
	// Update updates a device.
	Update(context.Context, *UpdateDeviceRequest) (*UpdateDeviceResponse, error)
	// Delete deletes a device by ID.
	Delete(context.Context, *DeleteDeviceRequest) (*empty.Empty, error)
	// List retrieves all devices.
	List(context.Context, *ListDeviceRequest) (*ListDeviceResponse, error)
	mustEmbedUnimplementedDeviceServiceServer()
}

// UnimplementedDeviceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceServiceServer struct {
}

func (UnimplementedDeviceServiceServer) Create(context.Context, *CreateDeviceRequest) (*CreateDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedDeviceServiceServer) Read(context.Context, *ReadDeviceRequest) (*ReadDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedDeviceServiceServer) Update(context.Context, *UpdateDeviceRequest) (*UpdateDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDeviceServiceServer) Delete(context.Context, *DeleteDeviceRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDeviceServiceServer) List(context.Context, *ListDeviceRequest) (*ListDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDeviceServiceServer) mustEmbedUnimplementedDeviceServiceServer() {}

// UnsafeDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceServiceServer will
// result in compilation errors.
type UnsafeDeviceServiceServer interface {
	mustEmbedUnimplementedDeviceServiceServer()
}

func RegisterDeviceServiceServer(s grpc.ServiceRegistrar, srv DeviceServiceServer) {
	s.RegisterService(&_DeviceService_serviceDesc, srv)
}

func _DeviceService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.DeviceService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Create(ctx, req.(*CreateDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.DeviceService/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Read(ctx, req.(*ReadDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.DeviceService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Update(ctx, req.(*UpdateDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.DeviceService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).Delete(ctx, req.(*DeleteDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.DeviceService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceServiceServer).List(ctx, req.(*ListDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DeviceService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.DeviceService",
	HandlerType: (*DeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _DeviceService_Create_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _DeviceService_Read_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _DeviceService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DeviceService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _DeviceService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/device.proto",
}
