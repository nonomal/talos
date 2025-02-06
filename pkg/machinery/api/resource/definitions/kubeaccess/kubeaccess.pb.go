// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: resource/definitions/kubeaccess/kubeaccess.proto

package kubeaccess

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ConfigSpec describes KubeSpan configuration..
type ConfigSpec struct {
	state                       protoimpl.MessageState `protogen:"open.v1"`
	Enabled                     bool                   `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	AllowedApiRoles             []string               `protobuf:"bytes,2,rep,name=allowed_api_roles,json=allowedApiRoles,proto3" json:"allowed_api_roles,omitempty"`
	AllowedKubernetesNamespaces []string               `protobuf:"bytes,3,rep,name=allowed_kubernetes_namespaces,json=allowedKubernetesNamespaces,proto3" json:"allowed_kubernetes_namespaces,omitempty"`
	unknownFields               protoimpl.UnknownFields
	sizeCache                   protoimpl.SizeCache
}

func (x *ConfigSpec) Reset() {
	*x = ConfigSpec{}
	mi := &file_resource_definitions_kubeaccess_kubeaccess_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConfigSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigSpec) ProtoMessage() {}

func (x *ConfigSpec) ProtoReflect() protoreflect.Message {
	mi := &file_resource_definitions_kubeaccess_kubeaccess_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigSpec.ProtoReflect.Descriptor instead.
func (*ConfigSpec) Descriptor() ([]byte, []int) {
	return file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigSpec) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *ConfigSpec) GetAllowedApiRoles() []string {
	if x != nil {
		return x.AllowedApiRoles
	}
	return nil
}

func (x *ConfigSpec) GetAllowedKubernetesNamespaces() []string {
	if x != nil {
		return x.AllowedKubernetesNamespaces
	}
	return nil
}

var File_resource_definitions_kubeaccess_kubeaccess_proto protoreflect.FileDescriptor

var file_resource_definitions_kubeaccess_kubeaccess_proto_rawDesc = []byte{
	0x0a, 0x30, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x64, 0x65, 0x66, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x25, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x2e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x6b,
	0x75, 0x62, 0x65, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x96, 0x01, 0x0a, 0x0a, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x53, 0x70, 0x65, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c,
	0x65, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x61, 0x70,
	0x69, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x61,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x41, 0x70, 0x69, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x42,
	0x0a, 0x1d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e,
	0x65, 0x74, 0x65, 0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x1b, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x4b, 0x75,
	0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x73, 0x42, 0x7e, 0x0a, 0x2d, 0x64, 0x65, 0x76, 0x2e, 0x74, 0x61, 0x6c, 0x6f, 0x73, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5a, 0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x73, 0x69, 0x64, 0x65, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x74, 0x61, 0x6c, 0x6f, 0x73,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x72, 0x79, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x64, 0x65, 0x66, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescOnce sync.Once
	file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescData = file_resource_definitions_kubeaccess_kubeaccess_proto_rawDesc
)

func file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescGZIP() []byte {
	file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescOnce.Do(func() {
		file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescData = protoimpl.X.CompressGZIP(file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescData)
	})
	return file_resource_definitions_kubeaccess_kubeaccess_proto_rawDescData
}

var file_resource_definitions_kubeaccess_kubeaccess_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_resource_definitions_kubeaccess_kubeaccess_proto_goTypes = []any{
	(*ConfigSpec)(nil), // 0: talos.resource.definitions.kubeaccess.ConfigSpec
}
var file_resource_definitions_kubeaccess_kubeaccess_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_resource_definitions_kubeaccess_kubeaccess_proto_init() }
func file_resource_definitions_kubeaccess_kubeaccess_proto_init() {
	if File_resource_definitions_kubeaccess_kubeaccess_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_resource_definitions_kubeaccess_kubeaccess_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_resource_definitions_kubeaccess_kubeaccess_proto_goTypes,
		DependencyIndexes: file_resource_definitions_kubeaccess_kubeaccess_proto_depIdxs,
		MessageInfos:      file_resource_definitions_kubeaccess_kubeaccess_proto_msgTypes,
	}.Build()
	File_resource_definitions_kubeaccess_kubeaccess_proto = out.File
	file_resource_definitions_kubeaccess_kubeaccess_proto_rawDesc = nil
	file_resource_definitions_kubeaccess_kubeaccess_proto_goTypes = nil
	file_resource_definitions_kubeaccess_kubeaccess_proto_depIdxs = nil
}
