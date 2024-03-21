// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.24.4
// source: api/thingspect_status.proto

package api

import (
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

// Status represents the status of a message.
type Status int32

const (
	// Status is not specified.
	Status_STATUS_UNSPECIFIED Status = 0
	// Message subject is active.
	Status_ACTIVE Status = 3
	// Message subject is disabled.
	Status_DISABLED Status = 6
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		3: "ACTIVE",
		6: "DISABLED",
	}
	Status_value = map[string]int32{
		"STATUS_UNSPECIFIED": 0,
		"ACTIVE":             3,
		"DISABLED":           6,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_api_thingspect_status_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_api_thingspect_status_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_api_thingspect_status_proto_rawDescGZIP(), []int{0}
}

var File_api_thingspect_status_proto protoreflect.FileDescriptor

var file_api_thingspect_status_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x74,
	0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2a, 0x3a, 0x0a,
	0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x44,
	0x49, 0x53, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x06, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65,
	0x63, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_thingspect_status_proto_rawDescOnce sync.Once
	file_api_thingspect_status_proto_rawDescData = file_api_thingspect_status_proto_rawDesc
)

func file_api_thingspect_status_proto_rawDescGZIP() []byte {
	file_api_thingspect_status_proto_rawDescOnce.Do(func() {
		file_api_thingspect_status_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_thingspect_status_proto_rawDescData)
	})
	return file_api_thingspect_status_proto_rawDescData
}

var file_api_thingspect_status_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_thingspect_status_proto_goTypes = []interface{}{
	(Status)(0), // 0: thingspect.api.Status
}
var file_api_thingspect_status_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_thingspect_status_proto_init() }
func file_api_thingspect_status_proto_init() {
	if File_api_thingspect_status_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_thingspect_status_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_thingspect_status_proto_goTypes,
		DependencyIndexes: file_api_thingspect_status_proto_depIdxs,
		EnumInfos:         file_api_thingspect_status_proto_enumTypes,
	}.Build()
	File_api_thingspect_status_proto = out.File
	file_api_thingspect_status_proto_rawDesc = nil
	file_api_thingspect_status_proto_goTypes = nil
	file_api_thingspect_status_proto_depIdxs = nil
}
