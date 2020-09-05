// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: event/payload.proto

package event

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/infobloxopen/protoc-gen-gorm/options"
	types "github.com/infobloxopen/protoc-gen-gorm/types"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Bet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           *types.UUIDValue     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreateTime   *timestamp.Timestamp `protobuf:"bytes,2,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	CustomerUuid string               `protobuf:"bytes,3,opt,name=customer_uuid,json=customerUuid,proto3" json:"customer_uuid,omitempty"`
	Stake        *Bet_Stake           `protobuf:"bytes,4,opt,name=stake,proto3" json:"stake,omitempty"`
	Selections   []*Selection         `protobuf:"bytes,5,rep,name=selections,proto3" json:"selections,omitempty"`
	Odds         string               `protobuf:"bytes,6,opt,name=odds,proto3" json:"odds,omitempty"`
}

func (x *Bet) Reset() {
	*x = Bet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_payload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bet) ProtoMessage() {}

func (x *Bet) ProtoReflect() protoreflect.Message {
	mi := &file_event_payload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bet.ProtoReflect.Descriptor instead.
func (*Bet) Descriptor() ([]byte, []int) {
	return file_event_payload_proto_rawDescGZIP(), []int{0}
}

func (x *Bet) GetId() *types.UUIDValue {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Bet) GetCreateTime() *timestamp.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Bet) GetCustomerUuid() string {
	if x != nil {
		return x.CustomerUuid
	}
	return ""
}

func (x *Bet) GetStake() *Bet_Stake {
	if x != nil {
		return x.Stake
	}
	return nil
}

func (x *Bet) GetSelections() []*Selection {
	if x != nil {
		return x.Selections
	}
	return nil
}

func (x *Bet) GetOdds() string {
	if x != nil {
		return x.Odds
	}
	return ""
}

type Selection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameId uint64 `protobuf:"varint,1,opt,name=game_id,json=gameId,proto3" json:"game_id,omitempty"`
	Market string `protobuf:"bytes,2,opt,name=market,proto3" json:"market,omitempty"`
	Odds   string `protobuf:"bytes,3,opt,name=odds,proto3" json:"odds,omitempty"`
}

func (x *Selection) Reset() {
	*x = Selection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_payload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Selection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Selection) ProtoMessage() {}

func (x *Selection) ProtoReflect() protoreflect.Message {
	mi := &file_event_payload_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Selection.ProtoReflect.Descriptor instead.
func (*Selection) Descriptor() ([]byte, []int) {
	return file_event_payload_proto_rawDescGZIP(), []int{1}
}

func (x *Selection) GetGameId() uint64 {
	if x != nil {
		return x.GameId
	}
	return 0
}

func (x *Selection) GetMarket() string {
	if x != nil {
		return x.Market
	}
	return ""
}

func (x *Selection) GetOdds() string {
	if x != nil {
		return x.Odds
	}
	return ""
}

type Bet_Stake struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value        string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Currency     string `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
	ExchangeRate string `protobuf:"bytes,3,opt,name=exchange_rate,json=exchangeRate,proto3" json:"exchange_rate,omitempty"`
}

func (x *Bet_Stake) Reset() {
	*x = Bet_Stake{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_payload_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bet_Stake) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bet_Stake) ProtoMessage() {}

func (x *Bet_Stake) ProtoReflect() protoreflect.Message {
	mi := &file_event_payload_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bet_Stake.ProtoReflect.Descriptor instead.
func (*Bet_Stake) Descriptor() ([]byte, []int) {
	return file_event_payload_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Bet_Stake) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Bet_Stake) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *Bet_Stake) GetExchangeRate() string {
	if x != nil {
		return x.ExchangeRate
	}
	return ""
}

var File_event_payload_proto protoreflect.FileDescriptor

var file_event_payload_proto_rawDesc = []byte{
	0x0a, 0x13, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x62, 0x6c, 0x6f,
	0x78, 0x6f, 0x70, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x67, 0x6f, 0x72, 0x6d, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x67, 0x6f,
	0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x62, 0x6c, 0x6f, 0x78, 0x6f, 0x70, 0x65,
	0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x67, 0x6f, 0x72,
	0x6d, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe4, 0x02, 0x0a, 0x03, 0x42, 0x65, 0x74, 0x12, 0x25, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x67, 0x6f, 0x72, 0x6d, 0x2e,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x3b, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65,
	0x72, 0x55, 0x75, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x42, 0x65, 0x74,
	0x2e, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x12, 0x30, 0x0a,
	0x0a, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x6f, 0x64, 0x64, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6f,
	0x64, 0x64, 0x73, 0x1a, 0x5e, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x23,
	0x0a, 0x0d, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52,
	0x61, 0x74, 0x65, 0x3a, 0x06, 0xba, 0xb9, 0x19, 0x02, 0x08, 0x01, 0x22, 0x58, 0x0a, 0x09, 0x53,
	0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x67, 0x61, 0x6d, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x64, 0x64,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6f, 0x64, 0x64, 0x73, 0x3a, 0x06, 0xba,
	0xb9, 0x19, 0x02, 0x08, 0x01, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x72, 0x74, 0x6b, 0x65, 0x2f, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x61, 0x72, 0x79, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x61, 0x64, 0x76,
	0x61, 0x6e, 0x63, 0x65, 0x64, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_event_payload_proto_rawDescOnce sync.Once
	file_event_payload_proto_rawDescData = file_event_payload_proto_rawDesc
)

func file_event_payload_proto_rawDescGZIP() []byte {
	file_event_payload_proto_rawDescOnce.Do(func() {
		file_event_payload_proto_rawDescData = protoimpl.X.CompressGZIP(file_event_payload_proto_rawDescData)
	})
	return file_event_payload_proto_rawDescData
}

var file_event_payload_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_event_payload_proto_goTypes = []interface{}{
	(*Bet)(nil),                 // 0: event.Bet
	(*Selection)(nil),           // 1: event.Selection
	(*Bet_Stake)(nil),           // 2: event.Bet.Stake
	(*types.UUIDValue)(nil),     // 3: gorm.types.UUIDValue
	(*timestamp.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_event_payload_proto_depIdxs = []int32{
	3, // 0: event.Bet.id:type_name -> gorm.types.UUIDValue
	4, // 1: event.Bet.create_time:type_name -> google.protobuf.Timestamp
	2, // 2: event.Bet.stake:type_name -> event.Bet.Stake
	1, // 3: event.Bet.selections:type_name -> event.Selection
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_event_payload_proto_init() }
func file_event_payload_proto_init() {
	if File_event_payload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_event_payload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bet); i {
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
		file_event_payload_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Selection); i {
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
		file_event_payload_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bet_Stake); i {
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
			RawDescriptor: file_event_payload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_event_payload_proto_goTypes,
		DependencyIndexes: file_event_payload_proto_depIdxs,
		MessageInfos:      file_event_payload_proto_msgTypes,
	}.Build()
	File_event_payload_proto = out.File
	file_event_payload_proto_rawDesc = nil
	file_event_payload_proto_goTypes = nil
	file_event_payload_proto_depIdxs = nil
}
