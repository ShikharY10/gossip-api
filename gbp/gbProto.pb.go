// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: gbProto.proto

package gbp

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

type Transport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg []byte `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Id  string `protobuf:"bytes,2,opt,name=Id,proto3" json:"Id,omitempty"`
	Tp  int32  `protobuf:"varint,3,opt,name=Tp,proto3" json:"Tp,omitempty"`
}

func (x *Transport) Reset() {
	*x = Transport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gbProto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transport) ProtoMessage() {}

func (x *Transport) ProtoReflect() protoreflect.Message {
	mi := &file_gbProto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transport.ProtoReflect.Descriptor instead.
func (*Transport) Descriptor() ([]byte, []int) {
	return file_gbProto_proto_rawDescGZIP(), []int{0}
}

func (x *Transport) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

func (x *Transport) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Transport) GetTp() int32 {
	if x != nil {
		return x.Tp
	}
	return 0
}

type HandshakeDeleteNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderMID string `protobuf:"bytes,1,opt,name=SenderMID,proto3" json:"SenderMID,omitempty"`
	TargetMID string `protobuf:"bytes,2,opt,name=TargetMID,proto3" json:"TargetMID,omitempty"`
	Number    string `protobuf:"bytes,3,opt,name=Number,proto3" json:"Number,omitempty"`
	Mloc      string `protobuf:"bytes,4,opt,name=Mloc,proto3" json:"Mloc,omitempty"`
}

func (x *HandshakeDeleteNotify) Reset() {
	*x = HandshakeDeleteNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gbProto_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeDeleteNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeDeleteNotify) ProtoMessage() {}

func (x *HandshakeDeleteNotify) ProtoReflect() protoreflect.Message {
	mi := &file_gbProto_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeDeleteNotify.ProtoReflect.Descriptor instead.
func (*HandshakeDeleteNotify) Descriptor() ([]byte, []int) {
	return file_gbProto_proto_rawDescGZIP(), []int{1}
}

func (x *HandshakeDeleteNotify) GetSenderMID() string {
	if x != nil {
		return x.SenderMID
	}
	return ""
}

func (x *HandshakeDeleteNotify) GetTargetMID() string {
	if x != nil {
		return x.TargetMID
	}
	return ""
}

func (x *HandshakeDeleteNotify) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *HandshakeDeleteNotify) GetMloc() string {
	if x != nil {
		return x.Mloc
	}
	return ""
}

type ChangeProfilePayloads struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	All       []string `protobuf:"bytes,1,rep,name=all,proto3" json:"all,omitempty"`
	PicData   string   `protobuf:"bytes,2,opt,name=PicData,proto3" json:"PicData,omitempty"`
	SenderMID string   `protobuf:"bytes,3,opt,name=SenderMID,proto3" json:"SenderMID,omitempty"`
}

func (x *ChangeProfilePayloads) Reset() {
	*x = ChangeProfilePayloads{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gbProto_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeProfilePayloads) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeProfilePayloads) ProtoMessage() {}

func (x *ChangeProfilePayloads) ProtoReflect() protoreflect.Message {
	mi := &file_gbProto_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeProfilePayloads.ProtoReflect.Descriptor instead.
func (*ChangeProfilePayloads) Descriptor() ([]byte, []int) {
	return file_gbProto_proto_rawDescGZIP(), []int{2}
}

func (x *ChangeProfilePayloads) GetAll() []string {
	if x != nil {
		return x.All
	}
	return nil
}

func (x *ChangeProfilePayloads) GetPicData() string {
	if x != nil {
		return x.PicData
	}
	return ""
}

func (x *ChangeProfilePayloads) GetSenderMID() string {
	if x != nil {
		return x.SenderMID
	}
	return ""
}

type NotifyChangeNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	All       []string `protobuf:"bytes,1,rep,name=all,proto3" json:"all,omitempty"`
	Number    string   `protobuf:"bytes,2,opt,name=Number,proto3" json:"Number,omitempty"`
	SenderMID string   `protobuf:"bytes,3,opt,name=SenderMID,proto3" json:"SenderMID,omitempty"`
}

func (x *NotifyChangeNumber) Reset() {
	*x = NotifyChangeNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gbProto_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotifyChangeNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotifyChangeNumber) ProtoMessage() {}

func (x *NotifyChangeNumber) ProtoReflect() protoreflect.Message {
	mi := &file_gbProto_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotifyChangeNumber.ProtoReflect.Descriptor instead.
func (*NotifyChangeNumber) Descriptor() ([]byte, []int) {
	return file_gbProto_proto_rawDescGZIP(), []int{3}
}

func (x *NotifyChangeNumber) GetAll() []string {
	if x != nil {
		return x.All
	}
	return nil
}

func (x *NotifyChangeNumber) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *NotifyChangeNumber) GetSenderMID() string {
	if x != nil {
		return x.SenderMID
	}
	return ""
}

var File_gbProto_proto protoreflect.FileDescriptor

var file_gbProto_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x67, 0x62, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x3d, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x4d, 0x73, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x54, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x54, 0x70, 0x22, 0x7f, 0x0a, 0x15, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b,
	0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x1c, 0x0a,
	0x09, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4d, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4d, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x4d, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4d, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x4d, 0x6c, 0x6f, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x4d, 0x6c, 0x6f, 0x63, 0x22, 0x61, 0x0a, 0x15, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x73, 0x12, 0x10,
	0x0a, 0x03, 0x61, 0x6c, 0x6c, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x61, 0x6c, 0x6c,
	0x12, 0x18, 0x0a, 0x07, 0x50, 0x69, 0x63, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x50, 0x69, 0x63, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x4d, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x4d, 0x49, 0x44, 0x22, 0x5c, 0x0a, 0x12, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x79, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x10,
	0x0a, 0x03, 0x61, 0x6c, 0x6c, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x61, 0x6c, 0x6c,
	0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x4d, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x4d, 0x49, 0x44, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x67, 0x62, 0x70, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gbProto_proto_rawDescOnce sync.Once
	file_gbProto_proto_rawDescData = file_gbProto_proto_rawDesc
)

func file_gbProto_proto_rawDescGZIP() []byte {
	file_gbProto_proto_rawDescOnce.Do(func() {
		file_gbProto_proto_rawDescData = protoimpl.X.CompressGZIP(file_gbProto_proto_rawDescData)
	})
	return file_gbProto_proto_rawDescData
}

var file_gbProto_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_gbProto_proto_goTypes = []interface{}{
	(*Transport)(nil),             // 0: main.Transport
	(*HandshakeDeleteNotify)(nil), // 1: main.HandshakeDeleteNotify
	(*ChangeProfilePayloads)(nil), // 2: main.ChangeProfilePayloads
	(*NotifyChangeNumber)(nil),    // 3: main.NotifyChangeNumber
}
var file_gbProto_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gbProto_proto_init() }
func file_gbProto_proto_init() {
	if File_gbProto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gbProto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transport); i {
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
		file_gbProto_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeDeleteNotify); i {
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
		file_gbProto_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeProfilePayloads); i {
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
		file_gbProto_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotifyChangeNumber); i {
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
			RawDescriptor: file_gbProto_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gbProto_proto_goTypes,
		DependencyIndexes: file_gbProto_proto_depIdxs,
		MessageInfos:      file_gbProto_proto_msgTypes,
	}.Build()
	File_gbProto_proto = out.File
	file_gbProto_proto_rawDesc = nil
	file_gbProto_proto_goTypes = nil
	file_gbProto_proto_depIdxs = nil
}
