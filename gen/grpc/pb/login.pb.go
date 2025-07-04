// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: pb/login.proto

package pb

import (
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

type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_pb_login_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_login_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_pb_login_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_pb_login_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_login_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_pb_login_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

var File_pb_login_proto protoreflect.FileDescriptor

const file_pb_login_proto_rawDesc = "" +
	"\n" +
	"\x0epb/login.proto\x12\x02pb\"@\n" +
	"\fLoginRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"2\n" +
	"\rLoginResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken2;\n" +
	"\vAuthService\x12,\n" +
	"\x05Login\x12\x10.pb.LoginRequest\x1a\x11.pb.LoginResponseBc\n" +
	"\x06com.pbB\n" +
	"LoginProtoP\x01Z%github.com/DucTran999/auth-service/pb\xa2\x02\x03PXX\xaa\x02\x02Pb\xca\x02\x02Pb\xe2\x02\x0ePb\\GPBMetadata\xea\x02\x02Pbb\x06proto3"

var (
	file_pb_login_proto_rawDescOnce sync.Once
	file_pb_login_proto_rawDescData []byte
)

func file_pb_login_proto_rawDescGZIP() []byte {
	file_pb_login_proto_rawDescOnce.Do(func() {
		file_pb_login_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pb_login_proto_rawDesc), len(file_pb_login_proto_rawDesc)))
	})
	return file_pb_login_proto_rawDescData
}

var file_pb_login_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_login_proto_goTypes = []any{
	(*LoginRequest)(nil),  // 0: pb.LoginRequest
	(*LoginResponse)(nil), // 1: pb.LoginResponse
}
var file_pb_login_proto_depIdxs = []int32{
	0, // 0: pb.AuthService.Login:input_type -> pb.LoginRequest
	1, // 1: pb.AuthService.Login:output_type -> pb.LoginResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_login_proto_init() }
func file_pb_login_proto_init() {
	if File_pb_login_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pb_login_proto_rawDesc), len(file_pb_login_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_login_proto_goTypes,
		DependencyIndexes: file_pb_login_proto_depIdxs,
		MessageInfos:      file_pb_login_proto_msgTypes,
	}.Build()
	File_pb_login_proto = out.File
	file_pb_login_proto_goTypes = nil
	file_pb_login_proto_depIdxs = nil
}
