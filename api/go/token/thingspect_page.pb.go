// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: token/thingspect_page.proto

package token

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Page represents a pagination token.
type Page struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Lower or upper bound timestamp, depending on ordering. Can represent any timestamp, but primarily used for created_at and representing the start of a page.
	BoundTs *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=bound_ts,json=boundTs,proto3" json:"bound_ts,omitempty"`
	// Previous ID (UUID). Can represent any UUID-based identifier.
	PrevId []byte `protobuf:"bytes,2,opt,name=prev_id,json=prevId,proto3" json:"prev_id,omitempty"`
}

func (x *Page) Reset() {
	*x = Page{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_thingspect_page_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Page) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Page) ProtoMessage() {}

func (x *Page) ProtoReflect() protoreflect.Message {
	mi := &file_token_thingspect_page_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Page.ProtoReflect.Descriptor instead.
func (*Page) Descriptor() ([]byte, []int) {
	return file_token_thingspect_page_proto_rawDescGZIP(), []int{0}
}

func (x *Page) GetBoundTs() *timestamppb.Timestamp {
	if x != nil {
		return x.BoundTs
	}
	return nil
}

func (x *Page) GetPrevId() []byte {
	if x != nil {
		return x.PrevId
	}
	return nil
}

var File_token_thingspect_page_proto protoreflect.FileDescriptor

var file_token_thingspect_page_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65,
	0x63, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x74,
	0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x69, 0x6e, 0x74, 0x2e, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x56, 0x0a, 0x04, 0x50, 0x61, 0x67, 0x65, 0x12, 0x35, 0x0a, 0x08,
	0x62, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x62, 0x6f, 0x75, 0x6e,
	0x64, 0x54, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x72, 0x65, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x72, 0x65, 0x76, 0x49, 0x64, 0x42, 0x2a, 0x5a, 0x28,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x74, 0x6c, 0x61, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x67, 0x6f, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_thingspect_page_proto_rawDescOnce sync.Once
	file_token_thingspect_page_proto_rawDescData = file_token_thingspect_page_proto_rawDesc
)

func file_token_thingspect_page_proto_rawDescGZIP() []byte {
	file_token_thingspect_page_proto_rawDescOnce.Do(func() {
		file_token_thingspect_page_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_thingspect_page_proto_rawDescData)
	})
	return file_token_thingspect_page_proto_rawDescData
}

var file_token_thingspect_page_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_token_thingspect_page_proto_goTypes = []interface{}{
	(*Page)(nil),                  // 0: thingspect.int.token.Page
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_token_thingspect_page_proto_depIdxs = []int32{
	1, // 0: thingspect.int.token.Page.bound_ts:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_token_thingspect_page_proto_init() }
func file_token_thingspect_page_proto_init() {
	if File_token_thingspect_page_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_token_thingspect_page_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Page); i {
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
			RawDescriptor: file_token_thingspect_page_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_token_thingspect_page_proto_goTypes,
		DependencyIndexes: file_token_thingspect_page_proto_depIdxs,
		MessageInfos:      file_token_thingspect_page_proto_msgTypes,
	}.Build()
	File_token_thingspect_page_proto = out.File
	file_token_thingspect_page_proto_rawDesc = nil
	file_token_thingspect_page_proto_goTypes = nil
	file_token_thingspect_page_proto_depIdxs = nil
}
