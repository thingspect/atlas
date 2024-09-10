// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.24.4
// source: message/thingspect_eventer_out.proto

package message

import (
	api "github.com/thingspect/proto/go/api"
	common "github.com/thingspect/proto/go/common"
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

// EventerOut represents a data point and associated metadata as used in message queues.
type EventerOut struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Data point.
	Point *common.DataPoint `protobuf:"bytes,1,opt,name=point,proto3" json:"point,omitempty"`
	// Device.
	Device *api.Device `protobuf:"bytes,2,opt,name=device,proto3" json:"device,omitempty"`
	// Rule.
	Rule *api.Rule `protobuf:"bytes,3,opt,name=rule,proto3" json:"rule,omitempty"`
}

func (x *EventerOut) Reset() {
	*x = EventerOut{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_thingspect_eventer_out_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventerOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventerOut) ProtoMessage() {}

func (x *EventerOut) ProtoReflect() protoreflect.Message {
	mi := &file_message_thingspect_eventer_out_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventerOut.ProtoReflect.Descriptor instead.
func (*EventerOut) Descriptor() ([]byte, []int) {
	return file_message_thingspect_eventer_out_proto_rawDescGZIP(), []int{0}
}

func (x *EventerOut) GetPoint() *common.DataPoint {
	if x != nil {
		return x.Point
	}
	return nil
}

func (x *EventerOut) GetDevice() *api.Device {
	if x != nil {
		return x.Device
	}
	return nil
}

func (x *EventerOut) GetRule() *api.Rule {
	if x != nil {
		return x.Rule
	}
	return nil
}

var File_message_thingspect_eventer_out_proto protoreflect.FileDescriptor

var file_message_thingspect_eventer_out_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x5f, 0x6f, 0x75, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65,
	0x63, 0x74, 0x2e, 0x69, 0x6e, 0x74, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1b,
	0x61, 0x70, 0x69, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x5f, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x61, 0x70, 0x69,
	0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x5f, 0x72, 0x75, 0x6c, 0x65,
	0x5f, 0x61, 0x6c, 0x61, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x5f,
	0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x9a, 0x01, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x12, 0x32,
	0x0a, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x05, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x75, 0x6c, 0x65, 0x42, 0x2e, 0x5a, 0x2c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x74, 0x6c, 0x61, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_thingspect_eventer_out_proto_rawDescOnce sync.Once
	file_message_thingspect_eventer_out_proto_rawDescData = file_message_thingspect_eventer_out_proto_rawDesc
)

func file_message_thingspect_eventer_out_proto_rawDescGZIP() []byte {
	file_message_thingspect_eventer_out_proto_rawDescOnce.Do(func() {
		file_message_thingspect_eventer_out_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_thingspect_eventer_out_proto_rawDescData)
	})
	return file_message_thingspect_eventer_out_proto_rawDescData
}

var file_message_thingspect_eventer_out_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_message_thingspect_eventer_out_proto_goTypes = []any{
	(*EventerOut)(nil),       // 0: thingspect.int.message.EventerOut
	(*common.DataPoint)(nil), // 1: thingspect.common.DataPoint
	(*api.Device)(nil),       // 2: thingspect.api.Device
	(*api.Rule)(nil),         // 3: thingspect.api.Rule
}
var file_message_thingspect_eventer_out_proto_depIdxs = []int32{
	1, // 0: thingspect.int.message.EventerOut.point:type_name -> thingspect.common.DataPoint
	2, // 1: thingspect.int.message.EventerOut.device:type_name -> thingspect.api.Device
	3, // 2: thingspect.int.message.EventerOut.rule:type_name -> thingspect.api.Rule
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_message_thingspect_eventer_out_proto_init() }
func file_message_thingspect_eventer_out_proto_init() {
	if File_message_thingspect_eventer_out_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_thingspect_eventer_out_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*EventerOut); i {
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
			RawDescriptor: file_message_thingspect_eventer_out_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_thingspect_eventer_out_proto_goTypes,
		DependencyIndexes: file_message_thingspect_eventer_out_proto_depIdxs,
		MessageInfos:      file_message_thingspect_eventer_out_proto_msgTypes,
	}.Build()
	File_message_thingspect_eventer_out_proto = out.File
	file_message_thingspect_eventer_out_proto_rawDesc = nil
	file_message_thingspect_eventer_out_proto_goTypes = nil
	file_message_thingspect_eventer_out_proto_depIdxs = nil
}
