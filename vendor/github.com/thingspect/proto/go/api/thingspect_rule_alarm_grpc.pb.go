// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v4.24.4
// source: api/thingspect_rule_alarm.proto

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
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	RuleAlarmService_CreateRule_FullMethodName  = "/thingspect.api.RuleAlarmService/CreateRule"
	RuleAlarmService_CreateAlarm_FullMethodName = "/thingspect.api.RuleAlarmService/CreateAlarm"
	RuleAlarmService_GetRule_FullMethodName     = "/thingspect.api.RuleAlarmService/GetRule"
	RuleAlarmService_GetAlarm_FullMethodName    = "/thingspect.api.RuleAlarmService/GetAlarm"
	RuleAlarmService_UpdateRule_FullMethodName  = "/thingspect.api.RuleAlarmService/UpdateRule"
	RuleAlarmService_UpdateAlarm_FullMethodName = "/thingspect.api.RuleAlarmService/UpdateAlarm"
	RuleAlarmService_DeleteRule_FullMethodName  = "/thingspect.api.RuleAlarmService/DeleteRule"
	RuleAlarmService_DeleteAlarm_FullMethodName = "/thingspect.api.RuleAlarmService/DeleteAlarm"
	RuleAlarmService_ListRules_FullMethodName   = "/thingspect.api.RuleAlarmService/ListRules"
	RuleAlarmService_ListAlarms_FullMethodName  = "/thingspect.api.RuleAlarmService/ListAlarms"
	RuleAlarmService_TestRule_FullMethodName    = "/thingspect.api.RuleAlarmService/TestRule"
	RuleAlarmService_TestAlarm_FullMethodName   = "/thingspect.api.RuleAlarmService/TestAlarm"
)

// RuleAlarmServiceClient is the client API for RuleAlarmService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// RuleAlarmService contains functions to query and modify rules and alarms.
type RuleAlarmServiceClient interface {
	// Create a rule. Rules generate events when conditions are met.
	CreateRule(ctx context.Context, in *CreateRuleRequest, opts ...grpc.CallOption) (*Rule, error)
	// Create an alarm. Alarms generate alerts via parent rules.
	CreateAlarm(ctx context.Context, in *CreateAlarmRequest, opts ...grpc.CallOption) (*Alarm, error)
	// Get a rule by ID. Rules generate events when conditions are met.
	GetRule(ctx context.Context, in *GetRuleRequest, opts ...grpc.CallOption) (*Rule, error)
	// Get an alarm by ID. Alarms generate alerts via parent rules.
	GetAlarm(ctx context.Context, in *GetAlarmRequest, opts ...grpc.CallOption) (*Alarm, error)
	// Update a rule. Rules generate events when conditions are met.
	UpdateRule(ctx context.Context, in *UpdateRuleRequest, opts ...grpc.CallOption) (*Rule, error)
	// Update an alarm. Alarms generate alerts via parent rules.
	UpdateAlarm(ctx context.Context, in *UpdateAlarmRequest, opts ...grpc.CallOption) (*Alarm, error)
	// Delete a rule by ID. Rules generate events when conditions are met.
	DeleteRule(ctx context.Context, in *DeleteRuleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete an alarm by ID. Alarms generate alerts via parent rules.
	DeleteAlarm(ctx context.Context, in *DeleteAlarmRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// List all rules. Rules generate events when conditions are met.
	ListRules(ctx context.Context, in *ListRulesRequest, opts ...grpc.CallOption) (*ListRulesResponse, error)
	// List alarms. Alarms generate alerts via parent rules.
	ListAlarms(ctx context.Context, in *ListAlarmsRequest, opts ...grpc.CallOption) (*ListAlarmsResponse, error)
	// Test a rule. Rules generate events when conditions are met.
	TestRule(ctx context.Context, in *TestRuleRequest, opts ...grpc.CallOption) (*TestRuleResponse, error)
	// Test an alarm. Alarms generate alerts via parent rules.
	TestAlarm(ctx context.Context, in *TestAlarmRequest, opts ...grpc.CallOption) (*TestAlarmResponse, error)
}

type ruleAlarmServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRuleAlarmServiceClient(cc grpc.ClientConnInterface) RuleAlarmServiceClient {
	return &ruleAlarmServiceClient{cc}
}

func (c *ruleAlarmServiceClient) CreateRule(ctx context.Context, in *CreateRuleRequest, opts ...grpc.CallOption) (*Rule, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Rule)
	err := c.cc.Invoke(ctx, RuleAlarmService_CreateRule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) CreateAlarm(ctx context.Context, in *CreateAlarmRequest, opts ...grpc.CallOption) (*Alarm, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Alarm)
	err := c.cc.Invoke(ctx, RuleAlarmService_CreateAlarm_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) GetRule(ctx context.Context, in *GetRuleRequest, opts ...grpc.CallOption) (*Rule, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Rule)
	err := c.cc.Invoke(ctx, RuleAlarmService_GetRule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) GetAlarm(ctx context.Context, in *GetAlarmRequest, opts ...grpc.CallOption) (*Alarm, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Alarm)
	err := c.cc.Invoke(ctx, RuleAlarmService_GetAlarm_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) UpdateRule(ctx context.Context, in *UpdateRuleRequest, opts ...grpc.CallOption) (*Rule, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Rule)
	err := c.cc.Invoke(ctx, RuleAlarmService_UpdateRule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) UpdateAlarm(ctx context.Context, in *UpdateAlarmRequest, opts ...grpc.CallOption) (*Alarm, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Alarm)
	err := c.cc.Invoke(ctx, RuleAlarmService_UpdateAlarm_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) DeleteRule(ctx context.Context, in *DeleteRuleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RuleAlarmService_DeleteRule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) DeleteAlarm(ctx context.Context, in *DeleteAlarmRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RuleAlarmService_DeleteAlarm_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) ListRules(ctx context.Context, in *ListRulesRequest, opts ...grpc.CallOption) (*ListRulesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRulesResponse)
	err := c.cc.Invoke(ctx, RuleAlarmService_ListRules_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) ListAlarms(ctx context.Context, in *ListAlarmsRequest, opts ...grpc.CallOption) (*ListAlarmsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAlarmsResponse)
	err := c.cc.Invoke(ctx, RuleAlarmService_ListAlarms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) TestRule(ctx context.Context, in *TestRuleRequest, opts ...grpc.CallOption) (*TestRuleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestRuleResponse)
	err := c.cc.Invoke(ctx, RuleAlarmService_TestRule_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleAlarmServiceClient) TestAlarm(ctx context.Context, in *TestAlarmRequest, opts ...grpc.CallOption) (*TestAlarmResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestAlarmResponse)
	err := c.cc.Invoke(ctx, RuleAlarmService_TestAlarm_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuleAlarmServiceServer is the server API for RuleAlarmService service.
// All implementations must embed UnimplementedRuleAlarmServiceServer
// for forward compatibility
//
// RuleAlarmService contains functions to query and modify rules and alarms.
type RuleAlarmServiceServer interface {
	// Create a rule. Rules generate events when conditions are met.
	CreateRule(context.Context, *CreateRuleRequest) (*Rule, error)
	// Create an alarm. Alarms generate alerts via parent rules.
	CreateAlarm(context.Context, *CreateAlarmRequest) (*Alarm, error)
	// Get a rule by ID. Rules generate events when conditions are met.
	GetRule(context.Context, *GetRuleRequest) (*Rule, error)
	// Get an alarm by ID. Alarms generate alerts via parent rules.
	GetAlarm(context.Context, *GetAlarmRequest) (*Alarm, error)
	// Update a rule. Rules generate events when conditions are met.
	UpdateRule(context.Context, *UpdateRuleRequest) (*Rule, error)
	// Update an alarm. Alarms generate alerts via parent rules.
	UpdateAlarm(context.Context, *UpdateAlarmRequest) (*Alarm, error)
	// Delete a rule by ID. Rules generate events when conditions are met.
	DeleteRule(context.Context, *DeleteRuleRequest) (*emptypb.Empty, error)
	// Delete an alarm by ID. Alarms generate alerts via parent rules.
	DeleteAlarm(context.Context, *DeleteAlarmRequest) (*emptypb.Empty, error)
	// List all rules. Rules generate events when conditions are met.
	ListRules(context.Context, *ListRulesRequest) (*ListRulesResponse, error)
	// List alarms. Alarms generate alerts via parent rules.
	ListAlarms(context.Context, *ListAlarmsRequest) (*ListAlarmsResponse, error)
	// Test a rule. Rules generate events when conditions are met.
	TestRule(context.Context, *TestRuleRequest) (*TestRuleResponse, error)
	// Test an alarm. Alarms generate alerts via parent rules.
	TestAlarm(context.Context, *TestAlarmRequest) (*TestAlarmResponse, error)
	mustEmbedUnimplementedRuleAlarmServiceServer()
}

// UnimplementedRuleAlarmServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRuleAlarmServiceServer struct {
}

func (UnimplementedRuleAlarmServiceServer) CreateRule(context.Context, *CreateRuleRequest) (*Rule, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRule not implemented")
}
func (UnimplementedRuleAlarmServiceServer) CreateAlarm(context.Context, *CreateAlarmRequest) (*Alarm, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAlarm not implemented")
}
func (UnimplementedRuleAlarmServiceServer) GetRule(context.Context, *GetRuleRequest) (*Rule, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRule not implemented")
}
func (UnimplementedRuleAlarmServiceServer) GetAlarm(context.Context, *GetAlarmRequest) (*Alarm, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAlarm not implemented")
}
func (UnimplementedRuleAlarmServiceServer) UpdateRule(context.Context, *UpdateRuleRequest) (*Rule, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRule not implemented")
}
func (UnimplementedRuleAlarmServiceServer) UpdateAlarm(context.Context, *UpdateAlarmRequest) (*Alarm, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAlarm not implemented")
}
func (UnimplementedRuleAlarmServiceServer) DeleteRule(context.Context, *DeleteRuleRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRule not implemented")
}
func (UnimplementedRuleAlarmServiceServer) DeleteAlarm(context.Context, *DeleteAlarmRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAlarm not implemented")
}
func (UnimplementedRuleAlarmServiceServer) ListRules(context.Context, *ListRulesRequest) (*ListRulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRules not implemented")
}
func (UnimplementedRuleAlarmServiceServer) ListAlarms(context.Context, *ListAlarmsRequest) (*ListAlarmsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAlarms not implemented")
}
func (UnimplementedRuleAlarmServiceServer) TestRule(context.Context, *TestRuleRequest) (*TestRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestRule not implemented")
}
func (UnimplementedRuleAlarmServiceServer) TestAlarm(context.Context, *TestAlarmRequest) (*TestAlarmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestAlarm not implemented")
}
func (UnimplementedRuleAlarmServiceServer) mustEmbedUnimplementedRuleAlarmServiceServer() {}

// UnsafeRuleAlarmServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RuleAlarmServiceServer will
// result in compilation errors.
type UnsafeRuleAlarmServiceServer interface {
	mustEmbedUnimplementedRuleAlarmServiceServer()
}

func RegisterRuleAlarmServiceServer(s grpc.ServiceRegistrar, srv RuleAlarmServiceServer) {
	s.RegisterService(&RuleAlarmService_ServiceDesc, srv)
}

func _RuleAlarmService_CreateRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).CreateRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_CreateRule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).CreateRule(ctx, req.(*CreateRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_CreateAlarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAlarmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).CreateAlarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_CreateAlarm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).CreateAlarm(ctx, req.(*CreateAlarmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_GetRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).GetRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_GetRule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).GetRule(ctx, req.(*GetRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_GetAlarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlarmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).GetAlarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_GetAlarm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).GetAlarm(ctx, req.(*GetAlarmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_UpdateRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).UpdateRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_UpdateRule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).UpdateRule(ctx, req.(*UpdateRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_UpdateAlarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAlarmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).UpdateAlarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_UpdateAlarm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).UpdateAlarm(ctx, req.(*UpdateAlarmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_DeleteRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).DeleteRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_DeleteRule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).DeleteRule(ctx, req.(*DeleteRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_DeleteAlarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAlarmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).DeleteAlarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_DeleteAlarm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).DeleteAlarm(ctx, req.(*DeleteAlarmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_ListRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).ListRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_ListRules_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).ListRules(ctx, req.(*ListRulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_ListAlarms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAlarmsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).ListAlarms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_ListAlarms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).ListAlarms(ctx, req.(*ListAlarmsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_TestRule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestRuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).TestRule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_TestRule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).TestRule(ctx, req.(*TestRuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuleAlarmService_TestAlarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestAlarmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleAlarmServiceServer).TestAlarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuleAlarmService_TestAlarm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleAlarmServiceServer).TestAlarm(ctx, req.(*TestAlarmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RuleAlarmService_ServiceDesc is the grpc.ServiceDesc for RuleAlarmService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RuleAlarmService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "thingspect.api.RuleAlarmService",
	HandlerType: (*RuleAlarmServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRule",
			Handler:    _RuleAlarmService_CreateRule_Handler,
		},
		{
			MethodName: "CreateAlarm",
			Handler:    _RuleAlarmService_CreateAlarm_Handler,
		},
		{
			MethodName: "GetRule",
			Handler:    _RuleAlarmService_GetRule_Handler,
		},
		{
			MethodName: "GetAlarm",
			Handler:    _RuleAlarmService_GetAlarm_Handler,
		},
		{
			MethodName: "UpdateRule",
			Handler:    _RuleAlarmService_UpdateRule_Handler,
		},
		{
			MethodName: "UpdateAlarm",
			Handler:    _RuleAlarmService_UpdateAlarm_Handler,
		},
		{
			MethodName: "DeleteRule",
			Handler:    _RuleAlarmService_DeleteRule_Handler,
		},
		{
			MethodName: "DeleteAlarm",
			Handler:    _RuleAlarmService_DeleteAlarm_Handler,
		},
		{
			MethodName: "ListRules",
			Handler:    _RuleAlarmService_ListRules_Handler,
		},
		{
			MethodName: "ListAlarms",
			Handler:    _RuleAlarmService_ListAlarms_Handler,
		},
		{
			MethodName: "TestRule",
			Handler:    _RuleAlarmService_TestRule_Handler,
		},
		{
			MethodName: "TestAlarm",
			Handler:    _RuleAlarmService_TestAlarm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/thingspect_rule_alarm.proto",
}
