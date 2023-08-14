// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0-devel
// 	protoc        v3.21.5
// source: syncer/proto/syncer.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// GetBlocksRequest is a request for GetBlocks
type GetBlocksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The height of beginning block to sync
	From uint64 `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
}

func (x *GetBlocksRequest) Reset() {
	*x = GetBlocksRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_syncer_proto_syncer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBlocksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBlocksRequest) ProtoMessage() {}

func (x *GetBlocksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_syncer_proto_syncer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBlocksRequest.ProtoReflect.Descriptor instead.
func (*GetBlocksRequest) Descriptor() ([]byte, []int) {
	return file_syncer_proto_syncer_proto_rawDescGZIP(), []int{0}
}

func (x *GetBlocksRequest) GetFrom() uint64 {
	if x != nil {
		return x.From
	}
	return 0
}

// Block contains a block data
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// RLP Encoded Block Data
	Block []byte `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_syncer_proto_syncer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_syncer_proto_syncer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_syncer_proto_syncer_proto_rawDescGZIP(), []int{1}
}

func (x *Block) GetBlock() []byte {
	if x != nil {
		return x.Block
	}
	return nil
}

// SyncPeerStatus contains peer status
type SyncPeerStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Latest block height
	Number uint64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *SyncPeerStatus) Reset() {
	*x = SyncPeerStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_syncer_proto_syncer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncPeerStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncPeerStatus) ProtoMessage() {}

func (x *SyncPeerStatus) ProtoReflect() protoreflect.Message {
	mi := &file_syncer_proto_syncer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncPeerStatus.ProtoReflect.Descriptor instead.
func (*SyncPeerStatus) Descriptor() ([]byte, []int) {
	return file_syncer_proto_syncer_proto_rawDescGZIP(), []int{2}
}

func (x *SyncPeerStatus) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

type GetInitConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddrInfo []byte `protobuf:"bytes,1,opt,name=addr_info,json=addrInfo,proto3" json:"addr_info,omitempty"`
}

func (x *GetInitConfigRequest) Reset() {
	*x = GetInitConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_syncer_proto_syncer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetInitConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInitConfigRequest) ProtoMessage() {}

func (x *GetInitConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_syncer_proto_syncer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInitConfigRequest.ProtoReflect.Descriptor instead.
func (*GetInitConfigRequest) Descriptor() ([]byte, []int) {
	return file_syncer_proto_syncer_proto_rawDescGZIP(), []int{3}
}

func (x *GetInitConfigRequest) GetAddrInfo() []byte {
	if x != nil {
		return x.AddrInfo
	}
	return nil
}

type InitConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *InitConfig) Reset() {
	*x = InitConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_syncer_proto_syncer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitConfig) ProtoMessage() {}

func (x *InitConfig) ProtoReflect() protoreflect.Message {
	mi := &file_syncer_proto_syncer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitConfig.ProtoReflect.Descriptor instead.
func (*InitConfig) Descriptor() ([]byte, []int) {
	return file_syncer_proto_syncer_proto_rawDescGZIP(), []int{4}
}

func (x *InitConfig) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

// SprintSnapshot contains a sprint snapshot data
type SprintSnapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SprintHeight uint64   `protobuf:"varint,1,opt,name=sprint_height,json=sprintHeight,proto3" json:"sprint_height,omitempty"`
	Addresses    [][]byte `protobuf:"bytes,2,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *SprintSnapshot) Reset() {
	*x = SprintSnapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_syncer_proto_syncer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SprintSnapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SprintSnapshot) ProtoMessage() {}

func (x *SprintSnapshot) ProtoReflect() protoreflect.Message {
	mi := &file_syncer_proto_syncer_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SprintSnapshot.ProtoReflect.Descriptor instead.
func (*SprintSnapshot) Descriptor() ([]byte, []int) {
	return file_syncer_proto_syncer_proto_rawDescGZIP(), []int{5}
}

func (x *SprintSnapshot) GetSprintHeight() uint64 {
	if x != nil {
		return x.SprintHeight
	}
	return 0
}

func (x *SprintSnapshot) GetAddresses() [][]byte {
	if x != nil {
		return x.Addresses
	}
	return nil
}

var File_syncer_proto_syncer_proto protoreflect.FileDescriptor

var file_syncer_proto_syncer_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x79, 0x6e, 0x63, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73,
	0x79, 0x6e, 0x63, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x26, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x66, 0x72, 0x6f, 0x6d, 0x22, 0x1d, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x14, 0x0a,
	0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x22, 0x28, 0x0a, 0x0e, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x65, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x33, 0x0a,
	0x14, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x61, 0x64, 0x64, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x22, 0x26, 0x0a, 0x0a, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x53, 0x0a, 0x0e, 0x53, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x23, 0x0a, 0x0d,
	0x73, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0c, 0x73, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0c, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x32,
	0xef, 0x01, 0x0a, 0x08, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x65, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x09,
	0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x12, 0x14, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x09, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x30, 0x01, 0x12, 0x37, 0x0a, 0x09,
	0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x65, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x69, 0x74,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x49,
	0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0e, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x3f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x53, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f,
	0x74, 0x42, 0x0f, 0x5a, 0x0d, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_syncer_proto_syncer_proto_rawDescOnce sync.Once
	file_syncer_proto_syncer_proto_rawDescData = file_syncer_proto_syncer_proto_rawDesc
)

func file_syncer_proto_syncer_proto_rawDescGZIP() []byte {
	file_syncer_proto_syncer_proto_rawDescOnce.Do(func() {
		file_syncer_proto_syncer_proto_rawDescData = protoimpl.X.CompressGZIP(file_syncer_proto_syncer_proto_rawDescData)
	})
	return file_syncer_proto_syncer_proto_rawDescData
}

var file_syncer_proto_syncer_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_syncer_proto_syncer_proto_goTypes = []interface{}{
	(*GetBlocksRequest)(nil),     // 0: v1.GetBlocksRequest
	(*Block)(nil),                // 1: v1.Block
	(*SyncPeerStatus)(nil),       // 2: v1.SyncPeerStatus
	(*GetInitConfigRequest)(nil), // 3: v1.GetInitConfigRequest
	(*InitConfig)(nil),           // 4: v1.InitConfig
	(*SprintSnapshot)(nil),       // 5: v1.SprintSnapshot
	(*emptypb.Empty)(nil),        // 6: google.protobuf.Empty
}
var file_syncer_proto_syncer_proto_depIdxs = []int32{
	0, // 0: v1.SyncPeer.GetBlocks:input_type -> v1.GetBlocksRequest
	6, // 1: v1.SyncPeer.GetStatus:input_type -> google.protobuf.Empty
	3, // 2: v1.SyncPeer.GetInitConfig:input_type -> v1.GetInitConfigRequest
	6, // 3: v1.SyncPeer.GetSprintSnapshot:input_type -> google.protobuf.Empty
	1, // 4: v1.SyncPeer.GetBlocks:output_type -> v1.Block
	2, // 5: v1.SyncPeer.GetStatus:output_type -> v1.SyncPeerStatus
	4, // 6: v1.SyncPeer.GetInitConfig:output_type -> v1.InitConfig
	5, // 7: v1.SyncPeer.GetSprintSnapshot:output_type -> v1.SprintSnapshot
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_syncer_proto_syncer_proto_init() }
func file_syncer_proto_syncer_proto_init() {
	if File_syncer_proto_syncer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_syncer_proto_syncer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBlocksRequest); i {
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
		file_syncer_proto_syncer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_syncer_proto_syncer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncPeerStatus); i {
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
		file_syncer_proto_syncer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetInitConfigRequest); i {
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
		file_syncer_proto_syncer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitConfig); i {
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
		file_syncer_proto_syncer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SprintSnapshot); i {
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
			RawDescriptor: file_syncer_proto_syncer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_syncer_proto_syncer_proto_goTypes,
		DependencyIndexes: file_syncer_proto_syncer_proto_depIdxs,
		MessageInfos:      file_syncer_proto_syncer_proto_msgTypes,
	}.Build()
	File_syncer_proto_syncer_proto = out.File
	file_syncer_proto_syncer_proto_rawDesc = nil
	file_syncer_proto_syncer_proto_goTypes = nil
	file_syncer_proto_syncer_proto_depIdxs = nil
}
