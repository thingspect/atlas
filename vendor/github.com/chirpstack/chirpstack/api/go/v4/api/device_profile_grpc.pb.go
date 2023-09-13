// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.9
// source: api/device_profile.proto

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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	DeviceProfileService_Create_FullMethodName            = "/api.DeviceProfileService/Create"
	DeviceProfileService_Get_FullMethodName               = "/api.DeviceProfileService/Get"
	DeviceProfileService_Update_FullMethodName            = "/api.DeviceProfileService/Update"
	DeviceProfileService_Delete_FullMethodName            = "/api.DeviceProfileService/Delete"
	DeviceProfileService_List_FullMethodName              = "/api.DeviceProfileService/List"
	DeviceProfileService_ListAdrAlgorithms_FullMethodName = "/api.DeviceProfileService/ListAdrAlgorithms"
)

// DeviceProfileServiceClient is the client API for DeviceProfileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceProfileServiceClient interface {
	// Create the given device-profile.
	Create(ctx context.Context, in *CreateDeviceProfileRequest, opts ...grpc.CallOption) (*CreateDeviceProfileResponse, error)
	// Get the device-profile for the given ID.
	Get(ctx context.Context, in *GetDeviceProfileRequest, opts ...grpc.CallOption) (*GetDeviceProfileResponse, error)
	// Update the given device-profile.
	Update(ctx context.Context, in *UpdateDeviceProfileRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete the device-profile with the given ID.
	Delete(ctx context.Context, in *DeleteDeviceProfileRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List the available device-profiles.
	List(ctx context.Context, in *ListDeviceProfilesRequest, opts ...grpc.CallOption) (*ListDeviceProfilesResponse, error)
	// List available ADR algorithms.
	ListAdrAlgorithms(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListDeviceProfileAdrAlgorithmsResponse, error)
}

type deviceProfileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceProfileServiceClient(cc grpc.ClientConnInterface) DeviceProfileServiceClient {
	return &deviceProfileServiceClient{cc}
}

func (c *deviceProfileServiceClient) Create(ctx context.Context, in *CreateDeviceProfileRequest, opts ...grpc.CallOption) (*CreateDeviceProfileResponse, error) {
	out := new(CreateDeviceProfileResponse)
	err := c.cc.Invoke(ctx, DeviceProfileService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceProfileServiceClient) Get(ctx context.Context, in *GetDeviceProfileRequest, opts ...grpc.CallOption) (*GetDeviceProfileResponse, error) {
	out := new(GetDeviceProfileResponse)
	err := c.cc.Invoke(ctx, DeviceProfileService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceProfileServiceClient) Update(ctx context.Context, in *UpdateDeviceProfileRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DeviceProfileService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceProfileServiceClient) Delete(ctx context.Context, in *DeleteDeviceProfileRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DeviceProfileService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceProfileServiceClient) List(ctx context.Context, in *ListDeviceProfilesRequest, opts ...grpc.CallOption) (*ListDeviceProfilesResponse, error) {
	out := new(ListDeviceProfilesResponse)
	err := c.cc.Invoke(ctx, DeviceProfileService_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceProfileServiceClient) ListAdrAlgorithms(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListDeviceProfileAdrAlgorithmsResponse, error) {
	out := new(ListDeviceProfileAdrAlgorithmsResponse)
	err := c.cc.Invoke(ctx, DeviceProfileService_ListAdrAlgorithms_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceProfileServiceServer is the server API for DeviceProfileService service.
// All implementations must embed UnimplementedDeviceProfileServiceServer
// for forward compatibility
type DeviceProfileServiceServer interface {
	// Create the given device-profile.
	Create(context.Context, *CreateDeviceProfileRequest) (*CreateDeviceProfileResponse, error)
	// Get the device-profile for the given ID.
	Get(context.Context, *GetDeviceProfileRequest) (*GetDeviceProfileResponse, error)
	// Update the given device-profile.
	Update(context.Context, *UpdateDeviceProfileRequest) (*emptypb.Empty, error)
	// Delete the device-profile with the given ID.
	Delete(context.Context, *DeleteDeviceProfileRequest) (*emptypb.Empty, error)
	// List the available device-profiles.
	List(context.Context, *ListDeviceProfilesRequest) (*ListDeviceProfilesResponse, error)
	// List available ADR algorithms.
	ListAdrAlgorithms(context.Context, *emptypb.Empty) (*ListDeviceProfileAdrAlgorithmsResponse, error)
	mustEmbedUnimplementedDeviceProfileServiceServer()
}

// UnimplementedDeviceProfileServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceProfileServiceServer struct {
}

func (UnimplementedDeviceProfileServiceServer) Create(context.Context, *CreateDeviceProfileRequest) (*CreateDeviceProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedDeviceProfileServiceServer) Get(context.Context, *GetDeviceProfileRequest) (*GetDeviceProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDeviceProfileServiceServer) Update(context.Context, *UpdateDeviceProfileRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDeviceProfileServiceServer) Delete(context.Context, *DeleteDeviceProfileRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDeviceProfileServiceServer) List(context.Context, *ListDeviceProfilesRequest) (*ListDeviceProfilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDeviceProfileServiceServer) ListAdrAlgorithms(context.Context, *emptypb.Empty) (*ListDeviceProfileAdrAlgorithmsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAdrAlgorithms not implemented")
}
func (UnimplementedDeviceProfileServiceServer) mustEmbedUnimplementedDeviceProfileServiceServer() {}

// UnsafeDeviceProfileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceProfileServiceServer will
// result in compilation errors.
type UnsafeDeviceProfileServiceServer interface {
	mustEmbedUnimplementedDeviceProfileServiceServer()
}

func RegisterDeviceProfileServiceServer(s grpc.ServiceRegistrar, srv DeviceProfileServiceServer) {
	s.RegisterService(&DeviceProfileService_ServiceDesc, srv)
}

func _DeviceProfileService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDeviceProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceProfileServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceProfileService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceProfileServiceServer).Create(ctx, req.(*CreateDeviceProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceProfileService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceProfileServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceProfileService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceProfileServiceServer).Get(ctx, req.(*GetDeviceProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceProfileService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDeviceProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceProfileServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceProfileService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceProfileServiceServer).Update(ctx, req.(*UpdateDeviceProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceProfileService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDeviceProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceProfileServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceProfileService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceProfileServiceServer).Delete(ctx, req.(*DeleteDeviceProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceProfileService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDeviceProfilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceProfileServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceProfileService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceProfileServiceServer).List(ctx, req.(*ListDeviceProfilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceProfileService_ListAdrAlgorithms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceProfileServiceServer).ListAdrAlgorithms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceProfileService_ListAdrAlgorithms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceProfileServiceServer).ListAdrAlgorithms(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// DeviceProfileService_ServiceDesc is the grpc.ServiceDesc for DeviceProfileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.DeviceProfileService",
	HandlerType: (*DeviceProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _DeviceProfileService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DeviceProfileService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _DeviceProfileService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DeviceProfileService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _DeviceProfileService_List_Handler,
		},
		{
			MethodName: "ListAdrAlgorithms",
			Handler:    _DeviceProfileService_ListAdrAlgorithms_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/device_profile.proto",
}
