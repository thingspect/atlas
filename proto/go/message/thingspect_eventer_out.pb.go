// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
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
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// EventerOut represents a data point and associated metadata as used in message queues.
type EventerOut struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Data point.
	Point *common.DataPoint `protobuf:"bytes,1,opt,name=point,proto3" json:"point,omitempty"`
	// Device.
	Device *api.Device `protobuf:"bytes,2,opt,name=device,proto3" json:"device,omitempty"`
	// Rule.
	Rule          *api.Rule `protobuf:"bytes,3,opt,name=rule,proto3" json:"rule,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventerOut) Reset() {
	*x = EventerOut{}
	mi := &file_message_thingspect_eventer_out_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventerOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventerOut) ProtoMessage() {}

func (x *EventerOut) ProtoReflect() protoreflect.Message {
	mi := &file_message_thingspect_eventer_out_proto_msgTypes[0]
	if x != nil {
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

const file_message_thingspect_eventer_out_proto_rawDesc = "" +
	"\n" +
	"$message/thingspect_eventer_out.proto\x12\x16thingspect.int.message\x1a\x1bapi/thingspect_device.proto\x1a\x1fapi/thingspect_rule_alarm.proto\x1a!common/thingspect_datapoint.proto\"\x9a\x01\n" +
	"\n" +
	"EventerOut\x122\n" +
	"\x05point\x18\x01 \x01(\v2\x1c.thingspect.common.DataPointR\x05point\x12.\n" +
	"\x06device\x18\x02 \x01(\v2\x16.thingspect.api.DeviceR\x06device\x12(\n" +
	"\x04rule\x18\x03 \x01(\v2\x14.thingspect.api.RuleR\x04ruleB.Z,github.com/thingspect/atlas/proto/go/messageb\x06proto3"

var (
	file_message_thingspect_eventer_out_proto_rawDescOnce sync.Once
	file_message_thingspect_eventer_out_proto_rawDescData []byte
)

func file_message_thingspect_eventer_out_proto_rawDescGZIP() []byte {
	file_message_thingspect_eventer_out_proto_rawDescOnce.Do(func() {
		file_message_thingspect_eventer_out_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_message_thingspect_eventer_out_proto_rawDesc), len(file_message_thingspect_eventer_out_proto_rawDesc)))
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_message_thingspect_eventer_out_proto_rawDesc), len(file_message_thingspect_eventer_out_proto_rawDesc)),
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
	file_message_thingspect_eventer_out_proto_goTypes = nil
	file_message_thingspect_eventer_out_proto_depIdxs = nil
}
