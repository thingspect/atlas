// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.18.1
// source: api/role.proto

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

// Role represents the role of a user.
type Role int32

const (
	// Role is not specified.
	Role_ROLE_UNSPECIFIED Role = 0
	// Contacts are not allowed to log in to the system, but can receive and respond to alerts from their organization.
	Role_CONTACT Role = 3
	// Viewers can only read resources in their organization, but can update their own user.
	Role_VIEWER Role = 6
	// Publishers can publish data points, but otherwise can only read resources in their organization.
	Role_PUBLISHER Role = 7
	// Builders can read and modify resources in their organization, but can only update their own user.
	Role_BUILDER Role = 9
	// Admins can read and modify anything in their organization, including creating users of an equal or lesser role.
	Role_ADMIN Role = 12
	// System admins can create organizations and modify anything in their organization.
	Role_SYS_ADMIN Role = 15
)

// Enum value maps for Role.
var (
	Role_name = map[int32]string{
		0:  "ROLE_UNSPECIFIED",
		3:  "CONTACT",
		6:  "VIEWER",
		7:  "PUBLISHER",
		9:  "BUILDER",
		12: "ADMIN",
		15: "SYS_ADMIN",
	}
	Role_value = map[string]int32{
		"ROLE_UNSPECIFIED": 0,
		"CONTACT":          3,
		"VIEWER":           6,
		"PUBLISHER":        7,
		"BUILDER":          9,
		"ADMIN":            12,
		"SYS_ADMIN":        15,
	}
)

func (x Role) Enum() *Role {
	p := new(Role)
	*p = x
	return p
}

func (x Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Role) Descriptor() protoreflect.EnumDescriptor {
	return file_api_role_proto_enumTypes[0].Descriptor()
}

func (Role) Type() protoreflect.EnumType {
	return &file_api_role_proto_enumTypes[0]
}

func (x Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Role.Descriptor instead.
func (Role) EnumDescriptor() ([]byte, []int) {
	return file_api_role_proto_rawDescGZIP(), []int{0}
}

var File_api_role_proto protoreflect.FileDescriptor

var file_api_role_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x2a, 0x6b, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x52, 0x4f, 0x4c, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b,
	0x0a, 0x07, 0x43, 0x4f, 0x4e, 0x54, 0x41, 0x43, 0x54, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x56,
	0x49, 0x45, 0x57, 0x45, 0x52, 0x10, 0x06, 0x12, 0x0d, 0x0a, 0x09, 0x50, 0x55, 0x42, 0x4c, 0x49,
	0x53, 0x48, 0x45, 0x52, 0x10, 0x07, 0x12, 0x0b, 0x0a, 0x07, 0x42, 0x55, 0x49, 0x4c, 0x44, 0x45,
	0x52, 0x10, 0x09, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x10, 0x0c, 0x12, 0x0d,
	0x0a, 0x09, 0x53, 0x59, 0x53, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x10, 0x0f, 0x42, 0x22, 0x5a,
	0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e,
	0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_role_proto_rawDescOnce sync.Once
	file_api_role_proto_rawDescData = file_api_role_proto_rawDesc
)

func file_api_role_proto_rawDescGZIP() []byte {
	file_api_role_proto_rawDescOnce.Do(func() {
		file_api_role_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_role_proto_rawDescData)
	})
	return file_api_role_proto_rawDescData
}

var file_api_role_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_role_proto_goTypes = []interface{}{
	(Role)(0), // 0: thingspect.api.Role
}
var file_api_role_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_role_proto_init() }
func file_api_role_proto_init() {
	if File_api_role_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_role_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_role_proto_goTypes,
		DependencyIndexes: file_api_role_proto_depIdxs,
		EnumInfos:         file_api_role_proto_enumTypes,
	}.Build()
	File_api_role_proto = out.File
	file_api_role_proto_rawDesc = nil
	file_api_role_proto_goTypes = nil
	file_api_role_proto_depIdxs = nil
}
