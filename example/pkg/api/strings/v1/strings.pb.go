// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: example/api/strings/v1/strings.proto

package stringsV1

import (
	v1 "github.com/loghole/tron/example/pkg/api/types/v1"
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

type ToUpperReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str string `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *ToUpperReq) Reset() {
	*x = ToUpperReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_api_strings_v1_strings_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToUpperReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToUpperReq) ProtoMessage() {}

func (x *ToUpperReq) ProtoReflect() protoreflect.Message {
	mi := &file_example_api_strings_v1_strings_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToUpperReq.ProtoReflect.Descriptor instead.
func (*ToUpperReq) Descriptor() ([]byte, []int) {
	return file_example_api_strings_v1_strings_proto_rawDescGZIP(), []int{0}
}

func (x *ToUpperReq) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

type ToUpperResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str *v1.String `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *ToUpperResp) Reset() {
	*x = ToUpperResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_api_strings_v1_strings_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToUpperResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToUpperResp) ProtoMessage() {}

func (x *ToUpperResp) ProtoReflect() protoreflect.Message {
	mi := &file_example_api_strings_v1_strings_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToUpperResp.ProtoReflect.Descriptor instead.
func (*ToUpperResp) Descriptor() ([]byte, []int) {
	return file_example_api_strings_v1_strings_proto_rawDescGZIP(), []int{1}
}

func (x *ToUpperResp) GetStr() *v1.String {
	if x != nil {
		return x.Str
	}
	return nil
}

type GetInfoReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str string `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *GetInfoReq) Reset() {
	*x = GetInfoReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_api_strings_v1_strings_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInfoReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInfoReq) ProtoMessage() {}

func (x *GetInfoReq) ProtoReflect() protoreflect.Message {
	mi := &file_example_api_strings_v1_strings_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInfoReq.ProtoReflect.Descriptor instead.
func (*GetInfoReq) Descriptor() ([]byte, []int) {
	return file_example_api_strings_v1_strings_proto_rawDescGZIP(), []int{2}
}

func (x *GetInfoReq) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

type GetInfoResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str *v1.String `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *GetInfoResp) Reset() {
	*x = GetInfoResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_api_strings_v1_strings_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInfoResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInfoResp) ProtoMessage() {}

func (x *GetInfoResp) ProtoReflect() protoreflect.Message {
	mi := &file_example_api_strings_v1_strings_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInfoResp.ProtoReflect.Descriptor instead.
func (*GetInfoResp) Descriptor() ([]byte, []int) {
	return file_example_api_strings_v1_strings_proto_rawDescGZIP(), []int{3}
}

func (x *GetInfoResp) GetStr() *v1.String {
	if x != nil {
		return x.Str
	}
	return nil
}

var File_example_api_strings_v1_strings_proto protoreflect.FileDescriptor

var file_example_api_strings_v1_strings_proto_rawDesc = []byte{
	0x0a, 0x24, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f,
	0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1e,
	0x0a, 0x0a, 0x54, 0x6f, 0x55, 0x70, 0x70, 0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03,
	0x73, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x22, 0x3d,
	0x0a, 0x0b, 0x54, 0x6f, 0x55, 0x70, 0x70, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2e, 0x0a,
	0x03, 0x73, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x03, 0x73, 0x74, 0x72, 0x22, 0x1e, 0x0a,
	0x0a, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x73,
	0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x22, 0x3d, 0x0a,
	0x0b, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x12, 0x2e, 0x0a, 0x03,
	0x73, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x03, 0x73, 0x74, 0x72, 0x32, 0xf7, 0x01, 0x0a,
	0x0a, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x41, 0x50, 0x49, 0x12, 0x77, 0x0a, 0x07, 0x54,
	0x6f, 0x55, 0x70, 0x70, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x6f, 0x55, 0x70, 0x70, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x23, 0x2e, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x55, 0x70, 0x70, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x22,
	0x23, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x12, 0x1b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x2f, 0x75, 0x70, 0x70, 0x65, 0x72, 0x2f, 0x7b,
	0x73, 0x74, 0x72, 0x7d, 0x12, 0x70, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x22, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x71, 0x1a, 0x23, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16,
	0x12, 0x14, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x73, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6f, 0x67, 0x68, 0x6f, 0x6c, 0x65, 0x2f, 0x74, 0x72, 0x6f,
	0x6e, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x73, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_example_api_strings_v1_strings_proto_rawDescOnce sync.Once
	file_example_api_strings_v1_strings_proto_rawDescData = file_example_api_strings_v1_strings_proto_rawDesc
)

func file_example_api_strings_v1_strings_proto_rawDescGZIP() []byte {
	file_example_api_strings_v1_strings_proto_rawDescOnce.Do(func() {
		file_example_api_strings_v1_strings_proto_rawDescData = protoimpl.X.CompressGZIP(file_example_api_strings_v1_strings_proto_rawDescData)
	})
	return file_example_api_strings_v1_strings_proto_rawDescData
}

var file_example_api_strings_v1_strings_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_example_api_strings_v1_strings_proto_goTypes = []interface{}{
	(*ToUpperReq)(nil),  // 0: example.api.strings.v1.ToUpperReq
	(*ToUpperResp)(nil), // 1: example.api.strings.v1.ToUpperResp
	(*GetInfoReq)(nil),  // 2: example.api.strings.v1.GetInfoReq
	(*GetInfoResp)(nil), // 3: example.api.strings.v1.GetInfoResp
	(*v1.String)(nil),   // 4: example.api.types.v1.String
}
var file_example_api_strings_v1_strings_proto_depIdxs = []int32{
	4, // 0: example.api.strings.v1.ToUpperResp.str:type_name -> example.api.types.v1.String
	4, // 1: example.api.strings.v1.GetInfoResp.str:type_name -> example.api.types.v1.String
	0, // 2: example.api.strings.v1.StringsAPI.ToUpper:input_type -> example.api.strings.v1.ToUpperReq
	2, // 3: example.api.strings.v1.StringsAPI.GetInfo:input_type -> example.api.strings.v1.GetInfoReq
	1, // 4: example.api.strings.v1.StringsAPI.ToUpper:output_type -> example.api.strings.v1.ToUpperResp
	3, // 5: example.api.strings.v1.StringsAPI.GetInfo:output_type -> example.api.strings.v1.GetInfoResp
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_example_api_strings_v1_strings_proto_init() }
func file_example_api_strings_v1_strings_proto_init() {
	if File_example_api_strings_v1_strings_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_example_api_strings_v1_strings_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToUpperReq); i {
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
		file_example_api_strings_v1_strings_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToUpperResp); i {
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
		file_example_api_strings_v1_strings_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInfoReq); i {
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
		file_example_api_strings_v1_strings_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInfoResp); i {
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
			RawDescriptor: file_example_api_strings_v1_strings_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_example_api_strings_v1_strings_proto_goTypes,
		DependencyIndexes: file_example_api_strings_v1_strings_proto_depIdxs,
		MessageInfos:      file_example_api_strings_v1_strings_proto_msgTypes,
	}.Build()
	File_example_api_strings_v1_strings_proto = out.File
	file_example_api_strings_v1_strings_proto_rawDesc = nil
	file_example_api_strings_v1_strings_proto_goTypes = nil
	file_example_api_strings_v1_strings_proto_depIdxs = nil
}
