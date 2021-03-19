// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.13.0
// source: api/rule.proto

package api

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	common "github.com/thingspect/api/go/common"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Rule represents a rule as stored in the database.
type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rule ID (UUID).
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Organization ID (UUID).
	OrgId string `protobuf:"bytes,2,opt,name=org_id,json=orgID,proto3" json:"org_id,omitempty"`
	// Rule name.
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// Rule status.
	Status common.Status `protobuf:"varint,4,opt,name=status,proto3,enum=thingspect.common.Status" json:"status,omitempty"`
	// Device tag to which the rule applies.
	DeviceTag string `protobuf:"bytes,5,opt,name=device_tag,json=deviceTag,proto3" json:"device_tag,omitempty"`
	// Device and data point attribute to which the rule applies.
	Attr string `protobuf:"bytes,6,opt,name=attr,proto3" json:"attr,omitempty"`
	// Rule expression. The rules engine evaluates a boolean expression using the [Expr language](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md).
	Expr string `protobuf:"bytes,7,opt,name=expr,proto3" json:"expr,omitempty"`
	// Rule creation timestamp.
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Rule modification timestamp.
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Rule) Reset() {
	*x = Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{0}
}

func (x *Rule) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Rule) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

func (x *Rule) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Rule) GetStatus() common.Status {
	if x != nil {
		return x.Status
	}
	return common.Status_STATUS_UNSPECIFIED
}

func (x *Rule) GetDeviceTag() string {
	if x != nil {
		return x.DeviceTag
	}
	return ""
}

func (x *Rule) GetAttr() string {
	if x != nil {
		return x.Attr
	}
	return ""
}

func (x *Rule) GetExpr() string {
	if x != nil {
		return x.Expr
	}
	return ""
}

func (x *Rule) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Rule) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// CreateRuleRequest is sent to create a rule.
type CreateRuleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rule message to create.
	Rule *Rule `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
}

func (x *CreateRuleRequest) Reset() {
	*x = CreateRuleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRuleRequest) ProtoMessage() {}

func (x *CreateRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRuleRequest.ProtoReflect.Descriptor instead.
func (*CreateRuleRequest) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRuleRequest) GetRule() *Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

// GetRuleRequest is sent to get a rule.
type GetRuleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rule ID (UUID) to get.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetRuleRequest) Reset() {
	*x = GetRuleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRuleRequest) ProtoMessage() {}

func (x *GetRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRuleRequest.ProtoReflect.Descriptor instead.
func (*GetRuleRequest) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{2}
}

func (x *GetRuleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// UpdateRuleRequest is sent to update a rule.
type UpdateRuleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rule message to update.
	Rule *Rule `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
	// Fields to update. Automatically populated by a PATCH request. If not present, a full resource update is performed.
	UpdateMask *field_mask.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
}

func (x *UpdateRuleRequest) Reset() {
	*x = UpdateRuleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRuleRequest) ProtoMessage() {}

func (x *UpdateRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRuleRequest.ProtoReflect.Descriptor instead.
func (*UpdateRuleRequest) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateRuleRequest) GetRule() *Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

func (x *UpdateRuleRequest) GetUpdateMask() *field_mask.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

// DeleteRuleRequest is sent to delete a rule.
type DeleteRuleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rule ID (UUID) to delete.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteRuleRequest) Reset() {
	*x = DeleteRuleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRuleRequest) ProtoMessage() {}

func (x *DeleteRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRuleRequest.ProtoReflect.Descriptor instead.
func (*DeleteRuleRequest) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRuleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// ListRulesRequest is sent to list rules.
type ListRulesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Number of rules to retrieve in a single page. Defaults to 50 if not specified, with a maximum of 250.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Token of the page to retrieve. If not specified, the first page of results will be returned. To request the next page of results, use next_page_token from the previous response.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
}

func (x *ListRulesRequest) Reset() {
	*x = ListRulesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRulesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRulesRequest) ProtoMessage() {}

func (x *ListRulesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRulesRequest.ProtoReflect.Descriptor instead.
func (*ListRulesRequest) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{5}
}

func (x *ListRulesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListRulesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

// ListRulesResponse is sent in response to a rule list.
type ListRulesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Rule array, ordered by ascending created_at timestamp.
	Rules []*Rule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
	// Pagination token used to retrieve the next page of results. Not returned for the last page.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	// Total number of rules available.
	TotalSize int32 `protobuf:"varint,3,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
}

func (x *ListRulesResponse) Reset() {
	*x = ListRulesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRulesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRulesResponse) ProtoMessage() {}

func (x *ListRulesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRulesResponse.ProtoReflect.Descriptor instead.
func (*ListRulesResponse) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{6}
}

func (x *ListRulesResponse) GetRules() []*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

func (x *ListRulesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

func (x *ListRulesResponse) GetTotalSize() int32 {
	if x != nil {
		return x.TotalSize
	}
	return 0
}

// TestRuleRequest is sent to test a rule.
type TestRuleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Data point to test against a rule.
	Point *common.DataPoint `protobuf:"bytes,1,opt,name=point,proto3" json:"point,omitempty"`
	// Rule message to test.
	Rule *Rule `protobuf:"bytes,2,opt,name=rule,proto3" json:"rule,omitempty"`
}

func (x *TestRuleRequest) Reset() {
	*x = TestRuleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRuleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRuleRequest) ProtoMessage() {}

func (x *TestRuleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRuleRequest.ProtoReflect.Descriptor instead.
func (*TestRuleRequest) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{7}
}

func (x *TestRuleRequest) GetPoint() *common.DataPoint {
	if x != nil {
		return x.Point
	}
	return nil
}

func (x *TestRuleRequest) GetRule() *Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

// TestRuleResponse is sent in response to a rule test.
type TestRuleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Result of the rule evaluation.
	Result bool `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *TestRuleResponse) Reset() {
	*x = TestRuleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_rule_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRuleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRuleResponse) ProtoMessage() {}

func (x *TestRuleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_rule_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRuleResponse.ProtoReflect.Descriptor instead.
func (*TestRuleResponse) Descriptor() ([]byte, []int) {
	return file_api_rule_proto_rawDescGZIP(), []int{8}
}

func (x *TestRuleResponse) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

var File_api_rule_proto protoreflect.FileDescriptor

var file_api_rule_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65,
	0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69,
	0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x88, 0x03, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x13,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x44, 0x12,
	0x20, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0c, 0xfa,
	0x42, 0x06, 0x72, 0x04, 0x10, 0x05, 0x18, 0x50, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x40, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x19, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x0d, 0xfa, 0x42,
	0x07, 0x82, 0x01, 0x04, 0x18, 0x03, 0x18, 0x06, 0xe0, 0x41, 0x02, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x2a, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x74, 0x61,
	0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xff,
	0x01, 0xe0, 0x41, 0x02, 0x52, 0x09, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x54, 0x61, 0x67, 0x12,
	0x1e, 0x0a, 0x04, 0x61, 0x74, 0x74, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa,
	0x42, 0x04, 0x72, 0x02, 0x18, 0x28, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x61, 0x74, 0x74, 0x72, 0x12,
	0x1f, 0x0a, 0x04, 0x65, 0x78, 0x70, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xfa,
	0x42, 0x05, 0x72, 0x03, 0x18, 0x80, 0x08, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x65, 0x78, 0x70, 0x72,
	0x12, 0x3e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x3e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x22, 0x4a, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x8a, 0x01,
	0x02, 0x10, 0x01, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x22, 0x2d, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x72,
	0x03, 0xb0, 0x01, 0x01, 0xe0, 0x41, 0x02, 0x52, 0x02, 0x69, 0x64, 0x22, 0x87, 0x01, 0x0a, 0x11,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x35, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x52, 0x75, 0x6c, 0x65, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0xe0,
	0x41, 0x02, 0x52, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x4d, 0x61, 0x73, 0x6b, 0x22, 0x30, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01,
	0xe0, 0x41, 0x02, 0x52, 0x02, 0x69, 0x64, 0x22, 0x5a, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x09, 0x70,
	0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x0a,
	0xfa, 0x42, 0x07, 0x1a, 0x05, 0x18, 0xfa, 0x01, 0x28, 0x00, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x86, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x72, 0x75, 0x6c,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a,
	0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x89, 0x01, 0x0a,
	0x0f, 0x54, 0x65, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x3f, 0x0a, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x42, 0x0b, 0xfa,
	0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0xe0, 0x41, 0x02, 0x52, 0x05, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x12, 0x35, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x52, 0x75, 0x6c, 0x65, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0xe0,
	0x41, 0x02, 0x52, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x22, 0x2a, 0x0a, 0x10, 0x54, 0x65, 0x73, 0x74,
	0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x32, 0xe4, 0x05, 0x0a, 0x0b, 0x52, 0x75, 0x6c, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x9c, 0x01, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52,
	0x75, 0x6c, 0x65, 0x12, 0x21, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70,
	0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x22, 0x55, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x11, 0x22, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x3a,
	0x04, 0x72, 0x75, 0x6c, 0x65, 0x92, 0x41, 0x3b, 0x4a, 0x39, 0x0a, 0x03, 0x32, 0x30, 0x31, 0x12,
	0x32, 0x0a, 0x16, 0x41, 0x20, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x20,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x12, 0x18, 0x0a, 0x16, 0x1a, 0x14, 0x2e,
	0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52,
	0x75, 0x6c, 0x65, 0x12, 0x57, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x1e,
	0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x52, 0x75, 0x6c, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x76,
	0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x85, 0x01, 0x0a,
	0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x21, 0x2e, 0x74, 0x68,
	0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x52, 0x75, 0x6c, 0x65, 0x22, 0x3e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x38, 0x1a, 0x13, 0x2f, 0x76,
	0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x72, 0x75, 0x6c, 0x65, 0x2e, 0x69, 0x64,
	0x7d, 0x3a, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x5a, 0x1b, 0x32, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x72,
	0x75, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x72, 0x75, 0x6c, 0x65, 0x2e, 0x69, 0x64, 0x7d, 0x3a, 0x04,
	0x72, 0x75, 0x6c, 0x65, 0x12, 0x85, 0x01, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x75, 0x6c, 0x65, 0x12, 0x21, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x75, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x3c,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x2a, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65,
	0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x92, 0x41, 0x23, 0x4a, 0x21, 0x0a, 0x03, 0x32, 0x30, 0x34,
	0x12, 0x1a, 0x0a, 0x16, 0x41, 0x20, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c,
	0x20, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x12, 0x00, 0x12, 0x63, 0x0a, 0x09,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x74, 0x68, 0x69, 0x6e,
	0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x74, 0x68,
	0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65,
	0x73, 0x12, 0x68, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x1f, 0x2e,
	0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x54,
	0x65, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x22, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75,
	0x6c, 0x65, 0x73, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x3a, 0x01, 0x2a, 0x42, 0x22, 0x5a, 0x20, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_rule_proto_rawDescOnce sync.Once
	file_api_rule_proto_rawDescData = file_api_rule_proto_rawDesc
)

func file_api_rule_proto_rawDescGZIP() []byte {
	file_api_rule_proto_rawDescOnce.Do(func() {
		file_api_rule_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_rule_proto_rawDescData)
	})
	return file_api_rule_proto_rawDescData
}

var file_api_rule_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_rule_proto_goTypes = []interface{}{
	(*Rule)(nil),                 // 0: thingspect.api.Rule
	(*CreateRuleRequest)(nil),    // 1: thingspect.api.CreateRuleRequest
	(*GetRuleRequest)(nil),       // 2: thingspect.api.GetRuleRequest
	(*UpdateRuleRequest)(nil),    // 3: thingspect.api.UpdateRuleRequest
	(*DeleteRuleRequest)(nil),    // 4: thingspect.api.DeleteRuleRequest
	(*ListRulesRequest)(nil),     // 5: thingspect.api.ListRulesRequest
	(*ListRulesResponse)(nil),    // 6: thingspect.api.ListRulesResponse
	(*TestRuleRequest)(nil),      // 7: thingspect.api.TestRuleRequest
	(*TestRuleResponse)(nil),     // 8: thingspect.api.TestRuleResponse
	(common.Status)(0),           // 9: thingspect.common.Status
	(*timestamp.Timestamp)(nil),  // 10: google.protobuf.Timestamp
	(*field_mask.FieldMask)(nil), // 11: google.protobuf.FieldMask
	(*common.DataPoint)(nil),     // 12: thingspect.common.DataPoint
	(*empty.Empty)(nil),          // 13: google.protobuf.Empty
}
var file_api_rule_proto_depIdxs = []int32{
	9,  // 0: thingspect.api.Rule.status:type_name -> thingspect.common.Status
	10, // 1: thingspect.api.Rule.created_at:type_name -> google.protobuf.Timestamp
	10, // 2: thingspect.api.Rule.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 3: thingspect.api.CreateRuleRequest.rule:type_name -> thingspect.api.Rule
	0,  // 4: thingspect.api.UpdateRuleRequest.rule:type_name -> thingspect.api.Rule
	11, // 5: thingspect.api.UpdateRuleRequest.update_mask:type_name -> google.protobuf.FieldMask
	0,  // 6: thingspect.api.ListRulesResponse.rules:type_name -> thingspect.api.Rule
	12, // 7: thingspect.api.TestRuleRequest.point:type_name -> thingspect.common.DataPoint
	0,  // 8: thingspect.api.TestRuleRequest.rule:type_name -> thingspect.api.Rule
	1,  // 9: thingspect.api.RuleService.CreateRule:input_type -> thingspect.api.CreateRuleRequest
	2,  // 10: thingspect.api.RuleService.GetRule:input_type -> thingspect.api.GetRuleRequest
	3,  // 11: thingspect.api.RuleService.UpdateRule:input_type -> thingspect.api.UpdateRuleRequest
	4,  // 12: thingspect.api.RuleService.DeleteRule:input_type -> thingspect.api.DeleteRuleRequest
	5,  // 13: thingspect.api.RuleService.ListRules:input_type -> thingspect.api.ListRulesRequest
	7,  // 14: thingspect.api.RuleService.TestRule:input_type -> thingspect.api.TestRuleRequest
	0,  // 15: thingspect.api.RuleService.CreateRule:output_type -> thingspect.api.Rule
	0,  // 16: thingspect.api.RuleService.GetRule:output_type -> thingspect.api.Rule
	0,  // 17: thingspect.api.RuleService.UpdateRule:output_type -> thingspect.api.Rule
	13, // 18: thingspect.api.RuleService.DeleteRule:output_type -> google.protobuf.Empty
	6,  // 19: thingspect.api.RuleService.ListRules:output_type -> thingspect.api.ListRulesResponse
	8,  // 20: thingspect.api.RuleService.TestRule:output_type -> thingspect.api.TestRuleResponse
	15, // [15:21] is the sub-list for method output_type
	9,  // [9:15] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_api_rule_proto_init() }
func file_api_rule_proto_init() {
	if File_api_rule_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_rule_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRuleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRuleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRuleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRuleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRulesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRulesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRuleRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_rule_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRuleResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_rule_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_rule_proto_goTypes,
		DependencyIndexes: file_api_rule_proto_depIdxs,
		MessageInfos:      file_api_rule_proto_msgTypes,
	}.Build()
	File_api_rule_proto = out.File
	file_api_rule_proto_rawDesc = nil
	file_api_rule_proto_goTypes = nil
	file_api_rule_proto_depIdxs = nil
}
