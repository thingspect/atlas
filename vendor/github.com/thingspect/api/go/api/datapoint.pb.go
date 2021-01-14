// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.2
// source: api/datapoint.proto

package api

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	common "github.com/thingspect/api/go/common"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// PublishDataPointRequest is sent to publish data points.
type PublishDataPointRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Data point array to publish.
	Points []*common.DataPoint `protobuf:"bytes,1,rep,name=points,proto3" json:"points,omitempty"`
}

func (x *PublishDataPointRequest) Reset() {
	*x = PublishDataPointRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_datapoint_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PublishDataPointRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishDataPointRequest) ProtoMessage() {}

func (x *PublishDataPointRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_datapoint_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishDataPointRequest.ProtoReflect.Descriptor instead.
func (*PublishDataPointRequest) Descriptor() ([]byte, []int) {
	return file_api_datapoint_proto_rawDescGZIP(), []int{0}
}

func (x *PublishDataPointRequest) GetPoints() []*common.DataPoint {
	if x != nil {
		return x.Points
	}
	return nil
}

// LatestDataPointRequest is sent to list latest device data points.
type LatestDataPointRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Identifier.
	//
	// Types that are assignable to IdOneof:
	//	*LatestDataPointRequest_UniqId
	//	*LatestDataPointRequest_DevId
	IdOneof isLatestDataPointRequest_IdOneof `protobuf_oneof:"id_oneof"`
}

func (x *LatestDataPointRequest) Reset() {
	*x = LatestDataPointRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_datapoint_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LatestDataPointRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatestDataPointRequest) ProtoMessage() {}

func (x *LatestDataPointRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_datapoint_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatestDataPointRequest.ProtoReflect.Descriptor instead.
func (*LatestDataPointRequest) Descriptor() ([]byte, []int) {
	return file_api_datapoint_proto_rawDescGZIP(), []int{1}
}

func (m *LatestDataPointRequest) GetIdOneof() isLatestDataPointRequest_IdOneof {
	if m != nil {
		return m.IdOneof
	}
	return nil
}

func (x *LatestDataPointRequest) GetUniqId() string {
	if x, ok := x.GetIdOneof().(*LatestDataPointRequest_UniqId); ok {
		return x.UniqId
	}
	return ""
}

func (x *LatestDataPointRequest) GetDevId() string {
	if x, ok := x.GetIdOneof().(*LatestDataPointRequest_DevId); ok {
		return x.DevId
	}
	return ""
}

type isLatestDataPointRequest_IdOneof interface {
	isLatestDataPointRequest_IdOneof()
}

type LatestDataPointRequest_UniqId struct {
	// Device unique ID.
	UniqId string `protobuf:"bytes,1,opt,name=uniq_id,json=uniqID,proto3,oneof"`
}

type LatestDataPointRequest_DevId struct {
	// Device ID (UUID).
	DevId string `protobuf:"bytes,2,opt,name=dev_id,json=devID,proto3,oneof"`
}

func (*LatestDataPointRequest_UniqId) isLatestDataPointRequest_IdOneof() {}

func (*LatestDataPointRequest_DevId) isLatestDataPointRequest_IdOneof() {}

// LatestDataPointResponse is sent in response to a device latest list.
type LatestDataPointResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Data point array.
	Points []*common.DataPoint `protobuf:"bytes,1,rep,name=points,proto3" json:"points,omitempty"`
}

func (x *LatestDataPointResponse) Reset() {
	*x = LatestDataPointResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_datapoint_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LatestDataPointResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatestDataPointResponse) ProtoMessage() {}

func (x *LatestDataPointResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_datapoint_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatestDataPointResponse.ProtoReflect.Descriptor instead.
func (*LatestDataPointResponse) Descriptor() ([]byte, []int) {
	return file_api_datapoint_proto_rawDescGZIP(), []int{2}
}

func (x *LatestDataPointResponse) GetPoints() []*common.DataPoint {
	if x != nil {
		return x.Points
	}
	return nil
}

var File_api_datapoint_proto protoreflect.FileDescriptor

var file_api_datapoint_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x16, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f,
	0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51, 0x0a, 0x17, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x36, 0x0a, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x92, 0x01, 0x02, 0x08, 0x01, 0xe0,
	0x41, 0x02, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x5d, 0x0a, 0x16, 0x4c, 0x61,
	0x74, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x71, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x49, 0x44, 0x12,
	0x17, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x05, 0x64, 0x65, 0x76, 0x49, 0x44, 0x42, 0x0f, 0x0a, 0x08, 0x69, 0x64, 0x5f, 0x6f,
	0x6e, 0x65, 0x6f, 0x66, 0x12, 0x03, 0xf8, 0x42, 0x01, 0x22, 0x44, 0x0a, 0x17, 0x4c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x32,
	0xd2, 0x01, 0x0a, 0x10, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x07, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x12,
	0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x44, 0x61, 0x74,
	0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x22, 0x0e, 0x2f,
	0x76, 0x31, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x3a, 0x01, 0x2a,
	0x12, 0x62, 0x0a, 0x06, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x61,
	0x74, 0x65, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f,
	0x76, 0x31, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2f, 0x6c, 0x61,
	0x74, 0x65, 0x73, 0x74, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_datapoint_proto_rawDescOnce sync.Once
	file_api_datapoint_proto_rawDescData = file_api_datapoint_proto_rawDesc
)

func file_api_datapoint_proto_rawDescGZIP() []byte {
	file_api_datapoint_proto_rawDescOnce.Do(func() {
		file_api_datapoint_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_datapoint_proto_rawDescData)
	})
	return file_api_datapoint_proto_rawDescData
}

var file_api_datapoint_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_datapoint_proto_goTypes = []interface{}{
	(*PublishDataPointRequest)(nil), // 0: api.PublishDataPointRequest
	(*LatestDataPointRequest)(nil),  // 1: api.LatestDataPointRequest
	(*LatestDataPointResponse)(nil), // 2: api.LatestDataPointResponse
	(*common.DataPoint)(nil),        // 3: common.DataPoint
	(*empty.Empty)(nil),             // 4: google.protobuf.Empty
}
var file_api_datapoint_proto_depIdxs = []int32{
	3, // 0: api.PublishDataPointRequest.points:type_name -> common.DataPoint
	3, // 1: api.LatestDataPointResponse.points:type_name -> common.DataPoint
	0, // 2: api.DataPointService.Publish:input_type -> api.PublishDataPointRequest
	1, // 3: api.DataPointService.Latest:input_type -> api.LatestDataPointRequest
	4, // 4: api.DataPointService.Publish:output_type -> google.protobuf.Empty
	2, // 5: api.DataPointService.Latest:output_type -> api.LatestDataPointResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_datapoint_proto_init() }
func file_api_datapoint_proto_init() {
	if File_api_datapoint_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_datapoint_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PublishDataPointRequest); i {
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
		file_api_datapoint_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LatestDataPointRequest); i {
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
		file_api_datapoint_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LatestDataPointResponse); i {
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
	file_api_datapoint_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*LatestDataPointRequest_UniqId)(nil),
		(*LatestDataPointRequest_DevId)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_datapoint_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_datapoint_proto_goTypes,
		DependencyIndexes: file_api_datapoint_proto_depIdxs,
		MessageInfos:      file_api_datapoint_proto_msgTypes,
	}.Build()
	File_api_datapoint_proto = out.File
	file_api_datapoint_proto_rawDesc = nil
	file_api_datapoint_proto_goTypes = nil
	file_api_datapoint_proto_depIdxs = nil
}
