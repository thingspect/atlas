// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: message/decoder_in.proto

package message

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// DecoderIn represents a data payload and associated metadata as used in message queues.
type DecoderIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Device unique ID.
	UniqId string `protobuf:"bytes,1,opt,name=uniq_id,json=uniqId,proto3" json:"uniq_id,omitempty"`
	// Data payload.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// Timestamp.
	Ts *timestamp.Timestamp `protobuf:"bytes,3,opt,name=ts,proto3" json:"ts,omitempty"`
	// Trace ID (UUID).
	TraceId string `protobuf:"bytes,4,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
}

func (x *DecoderIn) Reset() {
	*x = DecoderIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_decoder_in_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DecoderIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecoderIn) ProtoMessage() {}

func (x *DecoderIn) ProtoReflect() protoreflect.Message {
	mi := &file_message_decoder_in_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecoderIn.ProtoReflect.Descriptor instead.
func (*DecoderIn) Descriptor() ([]byte, []int) {
	return file_message_decoder_in_proto_rawDescGZIP(), []int{0}
}

func (x *DecoderIn) GetUniqId() string {
	if x != nil {
		return x.UniqId
	}
	return ""
}

func (x *DecoderIn) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DecoderIn) GetTs() *timestamp.Timestamp {
	if x != nil {
		return x.Ts
	}
	return nil
}

func (x *DecoderIn) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

var File_message_decoder_in_proto protoreflect.FileDescriptor

var file_message_decoder_in_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x64, 0x65, 0x63, 0x6f, 0x64, 0x65,
	0x72, 0x5f, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x09, 0x44, 0x65, 0x63, 0x6f, 0x64, 0x65, 0x72, 0x49,
	0x6e, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x71, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2a,
	0x0a, 0x02, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72,
	0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x72,
	0x61, 0x63, 0x65, 0x49, 0x64, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2f, 0x61,
	0x74, 0x6c, 0x61, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_decoder_in_proto_rawDescOnce sync.Once
	file_message_decoder_in_proto_rawDescData = file_message_decoder_in_proto_rawDesc
)

func file_message_decoder_in_proto_rawDescGZIP() []byte {
	file_message_decoder_in_proto_rawDescOnce.Do(func() {
		file_message_decoder_in_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_decoder_in_proto_rawDescData)
	})
	return file_message_decoder_in_proto_rawDescData
}

var file_message_decoder_in_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_message_decoder_in_proto_goTypes = []interface{}{
	(*DecoderIn)(nil),           // 0: message.DecoderIn
	(*timestamp.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_message_decoder_in_proto_depIdxs = []int32{
	1, // 0: message.DecoderIn.ts:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_message_decoder_in_proto_init() }
func file_message_decoder_in_proto_init() {
	if File_message_decoder_in_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_decoder_in_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DecoderIn); i {
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
			RawDescriptor: file_message_decoder_in_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_decoder_in_proto_goTypes,
		DependencyIndexes: file_message_decoder_in_proto_depIdxs,
		MessageInfos:      file_message_decoder_in_proto_msgTypes,
	}.Build()
	File_message_decoder_in_proto = out.File
	file_message_decoder_in_proto_rawDesc = nil
	file_message_decoder_in_proto_goTypes = nil
	file_message_decoder_in_proto_depIdxs = nil
}
