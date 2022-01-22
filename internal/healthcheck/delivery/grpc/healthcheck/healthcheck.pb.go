// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: internal/healthcheck/delivery/grpc/healthcheck/healthcheck.proto

package healthcheck

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

type HealthcheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthcheckRequest) Reset() {
	*x = HealthcheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthcheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthcheckRequest) ProtoMessage() {}

func (x *HealthcheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthcheckRequest.ProtoReflect.Descriptor instead.
func (*HealthcheckRequest) Descriptor() ([]byte, []int) {
	return file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescGZIP(), []int{0}
}

type HealthcheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Environment string `protobuf:"bytes,1,opt,name=environment,proto3" json:"environment,omitempty"`
	Status      string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Version     string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *HealthcheckResponse) Reset() {
	*x = HealthcheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthcheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthcheckResponse) ProtoMessage() {}

func (x *HealthcheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthcheckResponse.ProtoReflect.Descriptor instead.
func (*HealthcheckResponse) Descriptor() ([]byte, []int) {
	return file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescGZIP(), []int{1}
}

func (x *HealthcheckResponse) GetEnvironment() string {
	if x != nil {
		return x.Environment
	}
	return ""
}

func (x *HealthcheckResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *HealthcheckResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto protoreflect.FileDescriptor

var file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDesc = []byte{
	0x0a, 0x40, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2f, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b,
	0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x22,
	0x14, 0x0a, 0x12, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x69, 0x0a, 0x13, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x32, 0x61, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x12,
	0x52, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x1f,
	0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x48, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x20, 0x2e, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x74, 0x70, 0x6f, 0x2f, 0x74, 0x6f, 0x64,
	0x6f, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x68, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x68, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescOnce sync.Once
	file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescData = file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDesc
)

func file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescGZIP() []byte {
	file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescOnce.Do(func() {
		file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescData)
	})
	return file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDescData
}

var file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_goTypes = []interface{}{
	(*HealthcheckRequest)(nil),  // 0: healthcheck.HealthcheckRequest
	(*HealthcheckResponse)(nil), // 1: healthcheck.HealthcheckResponse
}
var file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_depIdxs = []int32{
	0, // 0: healthcheck.Healthcheck.Healthcheck:input_type -> healthcheck.HealthcheckRequest
	1, // 1: healthcheck.Healthcheck.Healthcheck:output_type -> healthcheck.HealthcheckResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_init() }
func file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_init() {
	if File_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthcheckRequest); i {
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
		file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthcheckResponse); i {
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
			RawDescriptor: file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_goTypes,
		DependencyIndexes: file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_depIdxs,
		MessageInfos:      file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_msgTypes,
	}.Build()
	File_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto = out.File
	file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_rawDesc = nil
	file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_goTypes = nil
	file_internal_healthcheck_delivery_grpc_healthcheck_healthcheck_proto_depIdxs = nil
}
