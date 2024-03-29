// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: postman/games/wheel/wheel_service.proto

package wheel

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

type WheelPick int32

const (
	WheelPick_Grey   WheelPick = 0
	WheelPick_Blue   WheelPick = 1
	WheelPick_Yellow WheelPick = 2
	WheelPick_Red    WheelPick = 3
)

// Enum value maps for WheelPick.
var (
	WheelPick_name = map[int32]string{
		0: "Grey",
		1: "Blue",
		2: "Yellow",
		3: "Red",
	}
	WheelPick_value = map[string]int32{
		"Grey":   0,
		"Blue":   1,
		"Yellow": 2,
		"Red":    3,
	}
)

func (x WheelPick) Enum() *WheelPick {
	p := new(WheelPick)
	*p = x
	return p
}

func (x WheelPick) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WheelPick) Descriptor() protoreflect.EnumDescriptor {
	return file_postman_games_wheel_wheel_service_proto_enumTypes[0].Descriptor()
}

func (WheelPick) Type() protoreflect.EnumType {
	return &file_postman_games_wheel_wheel_service_proto_enumTypes[0]
}

func (x WheelPick) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WheelPick.Descriptor instead.
func (WheelPick) EnumDescriptor() ([]byte, []int) {
	return file_postman_games_wheel_wheel_service_proto_rawDescGZIP(), []int{0}
}

type JoinWheelRoundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string    `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Bet    float64   `protobuf:"fixed64,2,opt,name=bet,proto3" json:"bet,omitempty"`
	Pick   WheelPick `protobuf:"varint,3,opt,name=pick,proto3,enum=wheel_service.WheelPick" json:"pick,omitempty"`
}

func (x *JoinWheelRoundRequest) Reset() {
	*x = JoinWheelRoundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_postman_games_wheel_wheel_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinWheelRoundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinWheelRoundRequest) ProtoMessage() {}

func (x *JoinWheelRoundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_postman_games_wheel_wheel_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinWheelRoundRequest.ProtoReflect.Descriptor instead.
func (*JoinWheelRoundRequest) Descriptor() ([]byte, []int) {
	return file_postman_games_wheel_wheel_service_proto_rawDescGZIP(), []int{0}
}

func (x *JoinWheelRoundRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *JoinWheelRoundRequest) GetBet() float64 {
	if x != nil {
		return x.Bet
	}
	return 0
}

func (x *JoinWheelRoundRequest) GetPick() WheelPick {
	if x != nil {
		return x.Pick
	}
	return WheelPick_Grey
}

var File_postman_games_wheel_wheel_service_proto protoreflect.FileDescriptor

var file_postman_games_wheel_wheel_service_proto_rawDesc = []byte{
	0x0a, 0x27, 0x70, 0x6f, 0x73, 0x74, 0x6d, 0x61, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x2f,
	0x77, 0x68, 0x65, 0x65, 0x6c, 0x2f, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x68, 0x65, 0x65, 0x6c,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a, 0x15, 0x4a, 0x6f, 0x69, 0x6e, 0x57, 0x68, 0x65,
	0x65, 0x6c, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x65, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x62, 0x65, 0x74, 0x12, 0x2c, 0x0a, 0x04, 0x70, 0x69, 0x63,
	0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x57, 0x68, 0x65, 0x65, 0x6c, 0x50, 0x69, 0x63,
	0x6b, 0x52, 0x04, 0x70, 0x69, 0x63, 0x6b, 0x2a, 0x34, 0x0a, 0x09, 0x57, 0x68, 0x65, 0x65, 0x6c,
	0x50, 0x69, 0x63, 0x6b, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x72, 0x65, 0x79, 0x10, 0x00, 0x12, 0x08,
	0x0a, 0x04, 0x42, 0x6c, 0x75, 0x65, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x59, 0x65, 0x6c, 0x6c,
	0x6f, 0x77, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x52, 0x65, 0x64, 0x10, 0x03, 0x32, 0x60, 0x0a,
	0x0c, 0x57, 0x68, 0x65, 0x65, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x50, 0x0a,
	0x0e, 0x4a, 0x6f, 0x69, 0x6e, 0x57, 0x68, 0x65, 0x65, 0x6c, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12,
	0x24, 0x2e, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x4a, 0x6f, 0x69, 0x6e, 0x57, 0x68, 0x65, 0x65, 0x6c, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42,
	0x16, 0x5a, 0x14, 0x2e, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65,
	0x73, 0x2f, 0x77, 0x68, 0x65, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_postman_games_wheel_wheel_service_proto_rawDescOnce sync.Once
	file_postman_games_wheel_wheel_service_proto_rawDescData = file_postman_games_wheel_wheel_service_proto_rawDesc
)

func file_postman_games_wheel_wheel_service_proto_rawDescGZIP() []byte {
	file_postman_games_wheel_wheel_service_proto_rawDescOnce.Do(func() {
		file_postman_games_wheel_wheel_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_postman_games_wheel_wheel_service_proto_rawDescData)
	})
	return file_postman_games_wheel_wheel_service_proto_rawDescData
}

var file_postman_games_wheel_wheel_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_postman_games_wheel_wheel_service_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_postman_games_wheel_wheel_service_proto_goTypes = []interface{}{
	(WheelPick)(0),                // 0: wheel_service.WheelPick
	(*JoinWheelRoundRequest)(nil), // 1: wheel_service.JoinWheelRoundRequest
	(*emptypb.Empty)(nil),         // 2: google.protobuf.Empty
}
var file_postman_games_wheel_wheel_service_proto_depIdxs = []int32{
	0, // 0: wheel_service.JoinWheelRoundRequest.pick:type_name -> wheel_service.WheelPick
	1, // 1: wheel_service.WheelService.JoinWheelRound:input_type -> wheel_service.JoinWheelRoundRequest
	2, // 2: wheel_service.WheelService.JoinWheelRound:output_type -> google.protobuf.Empty
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_postman_games_wheel_wheel_service_proto_init() }
func file_postman_games_wheel_wheel_service_proto_init() {
	if File_postman_games_wheel_wheel_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_postman_games_wheel_wheel_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinWheelRoundRequest); i {
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
			RawDescriptor: file_postman_games_wheel_wheel_service_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_postman_games_wheel_wheel_service_proto_goTypes,
		DependencyIndexes: file_postman_games_wheel_wheel_service_proto_depIdxs,
		EnumInfos:         file_postman_games_wheel_wheel_service_proto_enumTypes,
		MessageInfos:      file_postman_games_wheel_wheel_service_proto_msgTypes,
	}.Build()
	File_postman_games_wheel_wheel_service_proto = out.File
	file_postman_games_wheel_wheel_service_proto_rawDesc = nil
	file_postman_games_wheel_wheel_service_proto_goTypes = nil
	file_postman_games_wheel_wheel_service_proto_depIdxs = nil
}
