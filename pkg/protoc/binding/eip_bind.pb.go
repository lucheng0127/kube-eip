// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: pkg/protoc/binding/eip_bind.proto

package binding

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

type EipOpReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action  string `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
	EipAddr string `protobuf:"bytes,2,opt,name=eipAddr,proto3" json:"eipAddr,omitempty"`
	VmiAddr string `protobuf:"bytes,3,opt,name=vmiAddr,proto3" json:"vmiAddr,omitempty"`
}

func (x *EipOpReq) Reset() {
	*x = EipOpReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protoc_binding_eip_bind_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EipOpReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EipOpReq) ProtoMessage() {}

func (x *EipOpReq) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protoc_binding_eip_bind_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EipOpReq.ProtoReflect.Descriptor instead.
func (*EipOpReq) Descriptor() ([]byte, []int) {
	return file_pkg_protoc_binding_eip_bind_proto_rawDescGZIP(), []int{0}
}

func (x *EipOpReq) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *EipOpReq) GetEipAddr() string {
	if x != nil {
		return x.EipAddr
	}
	return ""
}

func (x *EipOpReq) GetVmiAddr() string {
	if x != nil {
		return x.VmiAddr
	}
	return ""
}

type EipOpRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result   string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	ErrPhase int32  `protobuf:"varint,2,opt,name=errPhase,proto3" json:"errPhase,omitempty"`
}

func (x *EipOpRsp) Reset() {
	*x = EipOpRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_protoc_binding_eip_bind_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EipOpRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EipOpRsp) ProtoMessage() {}

func (x *EipOpRsp) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protoc_binding_eip_bind_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EipOpRsp.ProtoReflect.Descriptor instead.
func (*EipOpRsp) Descriptor() ([]byte, []int) {
	return file_pkg_protoc_binding_eip_bind_proto_rawDescGZIP(), []int{1}
}

func (x *EipOpRsp) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *EipOpRsp) GetErrPhase() int32 {
	if x != nil {
		return x.ErrPhase
	}
	return 0
}

var File_pkg_protoc_binding_eip_bind_proto protoreflect.FileDescriptor

var file_pkg_protoc_binding_eip_bind_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2f, 0x62, 0x69, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x2f, 0x65, 0x69, 0x70, 0x5f, 0x62, 0x69, 0x6e, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x22, 0x56, 0x0a, 0x08, 0x45,
	0x69, 0x70, 0x4f, 0x70, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x65, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x6d, 0x69,
	0x41, 0x64, 0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x6d, 0x69, 0x41,
	0x64, 0x64, 0x72, 0x22, 0x3e, 0x0a, 0x08, 0x45, 0x69, 0x70, 0x4f, 0x70, 0x52, 0x73, 0x70, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x72, 0x72, 0x50, 0x68,
	0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x72, 0x72, 0x50, 0x68,
	0x61, 0x73, 0x65, 0x32, 0x3c, 0x0a, 0x08, 0x45, 0x69, 0x70, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12,
	0x30, 0x0a, 0x0a, 0x45, 0x69, 0x70, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x65, 0x12, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x45, 0x69, 0x70, 0x4f, 0x70, 0x52, 0x65, 0x71, 0x1a,
	0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2e, 0x45, 0x69, 0x70, 0x4f, 0x70, 0x52, 0x73,
	0x70, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6c, 0x75, 0x63, 0x68, 0x65, 0x6e, 0x67, 0x30, 0x31, 0x32, 0x37, 0x2f, 0x6b, 0x75, 0x62, 0x65,
	0x2d, 0x65, 0x69, 0x70, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2f,
	0x62, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_protoc_binding_eip_bind_proto_rawDescOnce sync.Once
	file_pkg_protoc_binding_eip_bind_proto_rawDescData = file_pkg_protoc_binding_eip_bind_proto_rawDesc
)

func file_pkg_protoc_binding_eip_bind_proto_rawDescGZIP() []byte {
	file_pkg_protoc_binding_eip_bind_proto_rawDescOnce.Do(func() {
		file_pkg_protoc_binding_eip_bind_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_protoc_binding_eip_bind_proto_rawDescData)
	})
	return file_pkg_protoc_binding_eip_bind_proto_rawDescData
}

var file_pkg_protoc_binding_eip_bind_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_protoc_binding_eip_bind_proto_goTypes = []interface{}{
	(*EipOpReq)(nil), // 0: protoc.EipOpReq
	(*EipOpRsp)(nil), // 1: protoc.EipOpRsp
}
var file_pkg_protoc_binding_eip_bind_proto_depIdxs = []int32{
	0, // 0: protoc.EipAgent.EipOperate:input_type -> protoc.EipOpReq
	1, // 1: protoc.EipAgent.EipOperate:output_type -> protoc.EipOpRsp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_protoc_binding_eip_bind_proto_init() }
func file_pkg_protoc_binding_eip_bind_proto_init() {
	if File_pkg_protoc_binding_eip_bind_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_protoc_binding_eip_bind_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EipOpReq); i {
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
		file_pkg_protoc_binding_eip_bind_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EipOpRsp); i {
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
			RawDescriptor: file_pkg_protoc_binding_eip_bind_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_protoc_binding_eip_bind_proto_goTypes,
		DependencyIndexes: file_pkg_protoc_binding_eip_bind_proto_depIdxs,
		MessageInfos:      file_pkg_protoc_binding_eip_bind_proto_msgTypes,
	}.Build()
	File_pkg_protoc_binding_eip_bind_proto = out.File
	file_pkg_protoc_binding_eip_bind_proto_rawDesc = nil
	file_pkg_protoc_binding_eip_bind_proto_goTypes = nil
	file_pkg_protoc_binding_eip_bind_proto_depIdxs = nil
}
