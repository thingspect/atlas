// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.18.1
// source: api/user.proto

package api

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
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

// User represents a user as stored in the database.
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// User ID (UUID).
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Organization ID (UUID).
	OrgId string `protobuf:"bytes,2,opt,name=org_id,json=orgID,proto3" json:"org_id,omitempty"`
	// User name.
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	// User email.
	Email string `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	// User phone number. If specified, must be in valid [E.164](https://www.twilio.com/docs/glossary/what-e164) format.
	Phone string `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	// User role.
	Role Role `protobuf:"varint,6,opt,name=role,proto3,enum=thingspect.api.Role" json:"role,omitempty"`
	// User status.
	Status Status `protobuf:"varint,7,opt,name=status,proto3,enum=thingspect.api.Status" json:"status,omitempty"`
	// User tags. Nested tags should be delimited by '/'.
	Tags []string `protobuf:"bytes,8,rep,name=tags,proto3" json:"tags,omitempty"`
	// Mobile application user key.
	AppKey string `protobuf:"bytes,9,opt,name=app_key,json=appKey,proto3" json:"app_key,omitempty"`
	// User creation timestamp.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// User modification timestamp.
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *User) GetRole() Role {
	if x != nil {
		return x.Role
	}
	return Role_ROLE_UNSPECIFIED
}

func (x *User) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_STATUS_UNSPECIFIED
}

func (x *User) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *User) GetAppKey() string {
	if x != nil {
		return x.AppKey
	}
	return ""
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// CreateUserRequest is sent to create a user.
type CreateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// User message to create.
	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserRequest.ProtoReflect.Descriptor instead.
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{1}
}

func (x *CreateUserRequest) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

// GetUserRequest is sent to get a user.
type GetUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// User ID (UUID) to get.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// UpdateUserRequest is sent to update a user.
type UpdateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// User message to update.
	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// Fields to update. Automatically populated by a PATCH request. If not present, a full resource update is performed.
	UpdateMask *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
}

func (x *UpdateUserRequest) Reset() {
	*x = UpdateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserRequest) ProtoMessage() {}

func (x *UpdateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserRequest) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateUserRequest) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *UpdateUserRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

// UpdateUserPasswordRequest is sent to update a user's password.
type UpdateUserPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// User ID (UUID) to update password.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// User password to update.
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *UpdateUserPasswordRequest) Reset() {
	*x = UpdateUserPasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserPasswordRequest) ProtoMessage() {}

func (x *UpdateUserPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserPasswordRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserPasswordRequest) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateUserPasswordRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateUserPasswordRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// DeleteUserRequest is sent to delete a user.
type DeleteUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// User ID (UUID) to delete.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteUserRequest) Reset() {
	*x = DeleteUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserRequest) ProtoMessage() {}

func (x *DeleteUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserRequest) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// ListUsersRequest is sent to list users.
type ListUsersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Number of users to retrieve in a single page. Defaults to 50 if not specified, with a maximum of 250.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Token of the page to retrieve. If not specified, the first page of results will be returned. To request the next page of results, use next_page_token from the previous response.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// User tag. If not specified, all tags are included.
	Tag string `protobuf:"bytes,3,opt,name=tag,proto3" json:"tag,omitempty"`
}

func (x *ListUsersRequest) Reset() {
	*x = ListUsersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUsersRequest) ProtoMessage() {}

func (x *ListUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUsersRequest.ProtoReflect.Descriptor instead.
func (*ListUsersRequest) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{6}
}

func (x *ListUsersRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListUsersRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListUsersRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

// ListUsersResponse is sent in response to a user list.
type ListUsersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// User array, ordered by ascending created_at timestamp. The completeness of the array will depend on whether the request was made by an admin user, non-admin user, or API key.
	Users []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	// Pagination token used to retrieve the next page of results. Not returned for the last page.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	// Total number of users available.
	TotalSize int32 `protobuf:"varint,3,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
}

func (x *ListUsersResponse) Reset() {
	*x = ListUsersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUsersResponse) ProtoMessage() {}

func (x *ListUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUsersResponse.ProtoReflect.Descriptor instead.
func (*ListUsersResponse) Descriptor() ([]byte, []int) {
	return file_api_user_proto_rawDescGZIP(), []int{7}
}

func (x *ListUsersResponse) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *ListUsersResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

func (x *ListUsersResponse) GetTotalSize() int32 {
	if x != nil {
		return x.TotalSize
	}
	return 0
}

var File_api_user_proto protoreflect.FileDescriptor

var file_api_user_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x1a, 0x0e, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70,
	0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x03, 0x0a, 0x04, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x03, 0xe0, 0x41, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x05, 0x6f,
	0x72, 0x67, 0x49, 0x44, 0x12, 0x1d, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x05, 0x18, 0x50, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x60, 0x01, 0x52, 0x05, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x1d, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x10, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e,
	0x65, 0x12, 0x3c, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x42, 0x12, 0xfa, 0x42, 0x0f, 0x82, 0x01, 0x0c, 0x18, 0x03, 0x18,
	0x06, 0x18, 0x07, 0x18, 0x09, 0x18, 0x0c, 0x18, 0x0f, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12,
	0x3a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x82, 0x01, 0x04, 0x18,
	0x03, 0x18, 0x06, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x29, 0x0a, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x42, 0x15, 0xfa, 0x42, 0x05, 0x92, 0x01,
	0x02, 0x18, 0x01, 0xfa, 0x42, 0x0a, 0x92, 0x01, 0x07, 0x22, 0x05, 0x72, 0x03, 0x18, 0xff, 0x01,
	0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x20, 0x0a, 0x07, 0x61, 0x70, 0x70, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x18, 0x2d,
	0x52, 0x06, 0x61, 0x70, 0x70, 0x4b, 0x65, 0x79, 0x12, 0x3e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x03, 0xe0, 0x41, 0x03, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x4a, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x68,
	0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0xe0, 0x41, 0x02, 0x52, 0x04,
	0x75, 0x73, 0x65, 0x72, 0x22, 0x2d, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0xe0, 0x41, 0x02, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x87, 0x01, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x42, 0x0b, 0xfa,
	0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72,
	0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73,
	0x6b, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x73, 0x6b, 0x22, 0x62, 0x0a,
	0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01,
	0xe0, 0x41, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x28, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0c, 0xfa, 0x42, 0x06, 0x72, 0x04,
	0x10, 0x0a, 0x18, 0x64, 0xe0, 0x41, 0x02, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x30, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0xe0, 0x41, 0x02, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x74, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x1a,
	0x03, 0x18, 0xfa, 0x01, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a,
	0x03, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72,
	0x03, 0x18, 0xff, 0x01, 0x52, 0x03, 0x74, 0x61, 0x67, 0x22, 0x86, 0x01, 0x0a, 0x11, 0x4c, 0x69,
	0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2a, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e,
	0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x53, 0x69,
	0x7a, 0x65, 0x32, 0xf7, 0x05, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x9c, 0x01, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x21, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63,
	0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x22, 0x55, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x11, 0x22, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x3a, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x92, 0x41, 0x3b, 0x4a, 0x39, 0x0a, 0x03, 0x32, 0x30, 0x31, 0x12, 0x32, 0x0a,
	0x16, 0x41, 0x20, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x20, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x12, 0x18, 0x0a, 0x16, 0x1a, 0x14, 0x2e, 0x74, 0x68,
	0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x57, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1e, 0x2e, 0x74,
	0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x74,
	0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x76, 0x31, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x85, 0x01, 0x0a, 0x0a, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x74, 0x68, 0x69, 0x6e,
	0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x74,
	0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x22, 0x3e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x38, 0x1a, 0x13, 0x2f, 0x76, 0x31, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x69, 0x64, 0x7d, 0x3a,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x5a, 0x1b, 0x32, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2f, 0x7b, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x69, 0x64, 0x7d, 0x3a, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x12, 0x7b, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x29, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x22, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x1c, 0x1a, 0x17, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x3a, 0x01, 0x2a, 0x12,
	0x85, 0x01, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x21,
	0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x3c, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x10, 0x2a, 0x0e, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x69, 0x64,
	0x7d, 0x92, 0x41, 0x23, 0x4a, 0x21, 0x0a, 0x03, 0x32, 0x30, 0x34, 0x12, 0x1a, 0x0a, 0x16, 0x41,
	0x20, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x66, 0x75, 0x6c, 0x20, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x12, 0x00, 0x12, 0x63, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x12, 0x20, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70, 0x65, 0x63,
	0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x73, 0x70,
	0x65, 0x63, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0b, 0x12, 0x09, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x42, 0x22, 0x5a, 0x20,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x69, 0x6e, 0x67,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_user_proto_rawDescOnce sync.Once
	file_api_user_proto_rawDescData = file_api_user_proto_rawDesc
)

func file_api_user_proto_rawDescGZIP() []byte {
	file_api_user_proto_rawDescOnce.Do(func() {
		file_api_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_user_proto_rawDescData)
	})
	return file_api_user_proto_rawDescData
}

var file_api_user_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_user_proto_goTypes = []interface{}{
	(*User)(nil),                      // 0: thingspect.api.User
	(*CreateUserRequest)(nil),         // 1: thingspect.api.CreateUserRequest
	(*GetUserRequest)(nil),            // 2: thingspect.api.GetUserRequest
	(*UpdateUserRequest)(nil),         // 3: thingspect.api.UpdateUserRequest
	(*UpdateUserPasswordRequest)(nil), // 4: thingspect.api.UpdateUserPasswordRequest
	(*DeleteUserRequest)(nil),         // 5: thingspect.api.DeleteUserRequest
	(*ListUsersRequest)(nil),          // 6: thingspect.api.ListUsersRequest
	(*ListUsersResponse)(nil),         // 7: thingspect.api.ListUsersResponse
	(Role)(0),                         // 8: thingspect.api.Role
	(Status)(0),                       // 9: thingspect.api.Status
	(*timestamppb.Timestamp)(nil),     // 10: google.protobuf.Timestamp
	(*fieldmaskpb.FieldMask)(nil),     // 11: google.protobuf.FieldMask
	(*emptypb.Empty)(nil),             // 12: google.protobuf.Empty
}
var file_api_user_proto_depIdxs = []int32{
	8,  // 0: thingspect.api.User.role:type_name -> thingspect.api.Role
	9,  // 1: thingspect.api.User.status:type_name -> thingspect.api.Status
	10, // 2: thingspect.api.User.created_at:type_name -> google.protobuf.Timestamp
	10, // 3: thingspect.api.User.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 4: thingspect.api.CreateUserRequest.user:type_name -> thingspect.api.User
	0,  // 5: thingspect.api.UpdateUserRequest.user:type_name -> thingspect.api.User
	11, // 6: thingspect.api.UpdateUserRequest.update_mask:type_name -> google.protobuf.FieldMask
	0,  // 7: thingspect.api.ListUsersResponse.users:type_name -> thingspect.api.User
	1,  // 8: thingspect.api.UserService.CreateUser:input_type -> thingspect.api.CreateUserRequest
	2,  // 9: thingspect.api.UserService.GetUser:input_type -> thingspect.api.GetUserRequest
	3,  // 10: thingspect.api.UserService.UpdateUser:input_type -> thingspect.api.UpdateUserRequest
	4,  // 11: thingspect.api.UserService.UpdateUserPassword:input_type -> thingspect.api.UpdateUserPasswordRequest
	5,  // 12: thingspect.api.UserService.DeleteUser:input_type -> thingspect.api.DeleteUserRequest
	6,  // 13: thingspect.api.UserService.ListUsers:input_type -> thingspect.api.ListUsersRequest
	0,  // 14: thingspect.api.UserService.CreateUser:output_type -> thingspect.api.User
	0,  // 15: thingspect.api.UserService.GetUser:output_type -> thingspect.api.User
	0,  // 16: thingspect.api.UserService.UpdateUser:output_type -> thingspect.api.User
	12, // 17: thingspect.api.UserService.UpdateUserPassword:output_type -> google.protobuf.Empty
	12, // 18: thingspect.api.UserService.DeleteUser:output_type -> google.protobuf.Empty
	7,  // 19: thingspect.api.UserService.ListUsers:output_type -> thingspect.api.ListUsersResponse
	14, // [14:20] is the sub-list for method output_type
	8,  // [8:14] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_api_user_proto_init() }
func file_api_user_proto_init() {
	if File_api_user_proto != nil {
		return
	}
	file_api_role_proto_init()
	file_api_status_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_api_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateUserRequest); i {
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
		file_api_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserRequest); i {
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
		file_api_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserRequest); i {
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
		file_api_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserPasswordRequest); i {
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
		file_api_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserRequest); i {
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
		file_api_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUsersRequest); i {
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
		file_api_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUsersResponse); i {
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
			RawDescriptor: file_api_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_user_proto_goTypes,
		DependencyIndexes: file_api_user_proto_depIdxs,
		MessageInfos:      file_api_user_proto_msgTypes,
	}.Build()
	File_api_user_proto = out.File
	file_api_user_proto_rawDesc = nil
	file_api_user_proto_goTypes = nil
	file_api_user_proto_depIdxs = nil
}
