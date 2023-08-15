// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: api/thingspect_alert.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AlertService_ListAlerts_FullMethodName = "/thingspect.api.AlertService/ListAlerts"
)

// AlertServiceClient is the client API for AlertService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AlertServiceClient interface {
	// List all alerts for a device, alarm, and/or user in an [end, start) time range, in descending timestamp order. Alerts are generated by alarms via parent rules.
	ListAlerts(ctx context.Context, in *ListAlertsRequest, opts ...grpc.CallOption) (*ListAlertsResponse, error)
}

type alertServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAlertServiceClient(cc grpc.ClientConnInterface) AlertServiceClient {
	return &alertServiceClient{cc}
}

func (c *alertServiceClient) ListAlerts(ctx context.Context, in *ListAlertsRequest, opts ...grpc.CallOption) (*ListAlertsResponse, error) {
	out := new(ListAlertsResponse)
	err := c.cc.Invoke(ctx, AlertService_ListAlerts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlertServiceServer is the server API for AlertService service.
// All implementations must embed UnimplementedAlertServiceServer
// for forward compatibility
type AlertServiceServer interface {
	// List all alerts for a device, alarm, and/or user in an [end, start) time range, in descending timestamp order. Alerts are generated by alarms via parent rules.
	ListAlerts(context.Context, *ListAlertsRequest) (*ListAlertsResponse, error)
	mustEmbedUnimplementedAlertServiceServer()
}

// UnimplementedAlertServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAlertServiceServer struct {
}

func (UnimplementedAlertServiceServer) ListAlerts(context.Context, *ListAlertsRequest) (*ListAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAlerts not implemented")
}
func (UnimplementedAlertServiceServer) mustEmbedUnimplementedAlertServiceServer() {}

// UnsafeAlertServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AlertServiceServer will
// result in compilation errors.
type UnsafeAlertServiceServer interface {
	mustEmbedUnimplementedAlertServiceServer()
}

func RegisterAlertServiceServer(s grpc.ServiceRegistrar, srv AlertServiceServer) {
	s.RegisterService(&AlertService_ServiceDesc, srv)
}

func _AlertService_ListAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertServiceServer).ListAlerts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AlertService_ListAlerts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertServiceServer).ListAlerts(ctx, req.(*ListAlertsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AlertService_ServiceDesc is the grpc.ServiceDesc for AlertService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AlertService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "thingspect.api.AlertService",
	HandlerType: (*AlertServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAlerts",
			Handler:    _AlertService_ListAlerts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/thingspect_alert.proto",
}
