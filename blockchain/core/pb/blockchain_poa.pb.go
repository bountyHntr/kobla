//go:build poa

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: blockchain/core/pb/blockchain_poa.proto

package pb

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

type TxStatus int32

const (
	TxStatus_Fail    TxStatus = 0
	TxStatus_Success TxStatus = 1
)

// Enum value maps for TxStatus.
var (
	TxStatus_name = map[int32]string{
		0: "Fail",
		1: "Success",
	}
	TxStatus_value = map[string]int32{
		"Fail":    0,
		"Success": 1,
	}
)

func (x TxStatus) Enum() *TxStatus {
	p := new(TxStatus)
	*p = x
	return p
}

func (x TxStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TxStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_blockchain_core_pb_blockchain_poa_proto_enumTypes[0].Descriptor()
}

func (TxStatus) Type() protoreflect.EnumType {
	return &file_blockchain_core_pb_blockchain_poa_proto_enumTypes[0]
}

func (x TxStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TxStatus.Descriptor instead.
func (TxStatus) EnumDescriptor() ([]byte, []int) {
	return file_blockchain_core_pb_blockchain_poa_proto_rawDescGZIP(), []int{0}
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp     int64          `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Number        int64          `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
	Transactions  []*Transaction `protobuf:"bytes,3,rep,name=transactions,proto3" json:"transactions,omitempty"`
	PrevBlockHash []byte         `protobuf:"bytes,4,opt,name=prev_block_hash,json=prevBlockHash,proto3" json:"prev_block_hash,omitempty"`
	Hash          []byte         `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
	Coinbase      []byte         `protobuf:"bytes,6,opt,name=coinbase,proto3" json:"coinbase,omitempty"`
	Signature     []byte         `protobuf:"bytes,7,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_core_pb_blockchain_poa_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_core_pb_blockchain_poa_proto_msgTypes[0]
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
	return file_blockchain_core_pb_blockchain_poa_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Block) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Block) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *Block) GetPrevBlockHash() []byte {
	if x != nil {
		return x.PrevBlockHash
	}
	return nil
}

func (x *Block) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Block) GetCoinbase() []byte {
	if x != nil {
		return x.Coinbase
	}
	return nil
}

func (x *Block) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender    []byte   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver  []byte   `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Amount    uint64   `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Data      []byte   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Hash      []byte   `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
	Status    TxStatus `protobuf:"varint,6,opt,name=status,proto3,enum=core.TxStatus" json:"status,omitempty"`
	Signature []byte   `protobuf:"bytes,7,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_core_pb_blockchain_poa_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_core_pb_blockchain_poa_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_blockchain_core_pb_blockchain_poa_proto_rawDescGZIP(), []int{1}
}

func (x *Transaction) GetSender() []byte {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *Transaction) GetReceiver() []byte {
	if x != nil {
		return x.Receiver
	}
	return nil
}

func (x *Transaction) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Transaction) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Transaction) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Transaction) GetStatus() TxStatus {
	if x != nil {
		return x.Status
	}
	return TxStatus_Fail
}

func (x *Transaction) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type ChainStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height         int64    `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	AddressFrom    string   `protobuf:"bytes,2,opt,name=address_from,json=addressFrom,proto3" json:"address_from,omitempty"`
	KnownAddresses []string `protobuf:"bytes,3,rep,name=known_addresses,json=knownAddresses,proto3" json:"known_addresses,omitempty"`
}

func (x *ChainStatus) Reset() {
	*x = ChainStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockchain_core_pb_blockchain_poa_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChainStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChainStatus) ProtoMessage() {}

func (x *ChainStatus) ProtoReflect() protoreflect.Message {
	mi := &file_blockchain_core_pb_blockchain_poa_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChainStatus.ProtoReflect.Descriptor instead.
func (*ChainStatus) Descriptor() ([]byte, []int) {
	return file_blockchain_core_pb_blockchain_poa_proto_rawDescGZIP(), []int{2}
}

func (x *ChainStatus) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *ChainStatus) GetAddressFrom() string {
	if x != nil {
		return x.AddressFrom
	}
	return ""
}

func (x *ChainStatus) GetKnownAddresses() []string {
	if x != nil {
		return x.KnownAddresses
	}
	return nil
}

var File_blockchain_core_pb_blockchain_poa_proto protoreflect.FileDescriptor

var file_blockchain_core_pb_blockchain_poa_proto_rawDesc = []byte{
	0x0a, 0x27, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x70, 0x62, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f,
	0x70, 0x6f, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x6f, 0x72, 0x65, 0x22,
	0xea, 0x01, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x35, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x70, 0x72, 0x65, 0x76, 0x5f, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0d, 0x70, 0x72, 0x65, 0x76, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12,
	0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0xc7, 0x01, 0x0a,
	0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x26, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x54, 0x78, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x71, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x21, 0x0a,
	0x0c, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x46, 0x72, 0x6f, 0x6d,
	0x12, 0x27, 0x0a, 0x0f, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x2a, 0x21, 0x0a, 0x08, 0x54, 0x78, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x61, 0x69, 0x6c, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x10, 0x01, 0x42, 0x16, 0x5a, 0x14,
	0x2e, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blockchain_core_pb_blockchain_poa_proto_rawDescOnce sync.Once
	file_blockchain_core_pb_blockchain_poa_proto_rawDescData = file_blockchain_core_pb_blockchain_poa_proto_rawDesc
)

func file_blockchain_core_pb_blockchain_poa_proto_rawDescGZIP() []byte {
	file_blockchain_core_pb_blockchain_poa_proto_rawDescOnce.Do(func() {
		file_blockchain_core_pb_blockchain_poa_proto_rawDescData = protoimpl.X.CompressGZIP(file_blockchain_core_pb_blockchain_poa_proto_rawDescData)
	})
	return file_blockchain_core_pb_blockchain_poa_proto_rawDescData
}

var file_blockchain_core_pb_blockchain_poa_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_blockchain_core_pb_blockchain_poa_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_blockchain_core_pb_blockchain_poa_proto_goTypes = []interface{}{
	(TxStatus)(0),       // 0: core.TxStatus
	(*Block)(nil),       // 1: core.Block
	(*Transaction)(nil), // 2: core.Transaction
	(*ChainStatus)(nil), // 3: core.ChainStatus
}
var file_blockchain_core_pb_blockchain_poa_proto_depIdxs = []int32{
	2, // 0: core.Block.transactions:type_name -> core.Transaction
	0, // 1: core.Transaction.status:type_name -> core.TxStatus
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_blockchain_core_pb_blockchain_poa_proto_init() }
func file_blockchain_core_pb_blockchain_poa_proto_init() {
	if File_blockchain_core_pb_blockchain_poa_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_blockchain_core_pb_blockchain_poa_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_blockchain_core_pb_blockchain_poa_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
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
		file_blockchain_core_pb_blockchain_poa_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChainStatus); i {
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
			RawDescriptor: file_blockchain_core_pb_blockchain_poa_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_blockchain_core_pb_blockchain_poa_proto_goTypes,
		DependencyIndexes: file_blockchain_core_pb_blockchain_poa_proto_depIdxs,
		EnumInfos:         file_blockchain_core_pb_blockchain_poa_proto_enumTypes,
		MessageInfos:      file_blockchain_core_pb_blockchain_poa_proto_msgTypes,
	}.Build()
	File_blockchain_core_pb_blockchain_poa_proto = out.File
	file_blockchain_core_pb_blockchain_poa_proto_rawDesc = nil
	file_blockchain_core_pb_blockchain_poa_proto_goTypes = nil
	file_blockchain_core_pb_blockchain_poa_proto_depIdxs = nil
}
