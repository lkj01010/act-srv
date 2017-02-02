// Code generated by protoc-gen-go.
// source: game/logic/gamepb/game.proto
// DO NOT EDIT!

/*
Package gamepb is a generated protocol buffer package.

It is generated from these files:
	game/logic/gamepb/game.proto

It has these top-level messages:
	Player
	BattleInfoNotify
	Property
	Scene
	SceneNotify
	EnterGameReq
	EnterGameAck
*/
package gamepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Property_Type int32

const (
	Property_GEN   Property_Type = 0
	Property_HEART Property_Type = 1
)

var Property_Type_name = map[int32]string{
	0: "GEN",
	1: "HEART",
}
var Property_Type_value = map[string]int32{
	"GEN":   0,
	"HEART": 1,
}

func (x Property_Type) String() string {
	return proto.EnumName(Property_Type_name, int32(x))
}
func (Property_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

type Player struct {
	Id         string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Level      int32  `protobuf:"varint,3,opt,name=level" json:"level,omitempty"`
	Figure     int32  `protobuf:"varint,4,opt,name=figure" json:"figure,omitempty"`
	WeaponType int32  `protobuf:"varint,5,opt,name=weapon_type,json=weaponType" json:"weapon_type,omitempty"`
	Score      int32  `protobuf:"varint,6,opt,name=score" json:"score,omitempty"`
	X          int32  `protobuf:"varint,10,opt,name=x" json:"x,omitempty"`
	Y          int32  `protobuf:"varint,11,opt,name=y" json:"y,omitempty"`
	Life       int32  `protobuf:"varint,12,opt,name=life" json:"life,omitempty"`
	CurLife    int32  `protobuf:"varint,13,opt,name=cur_life,json=curLife" json:"cur_life,omitempty"`
	DefLevel   int32  `protobuf:"varint,14,opt,name=def_level,json=defLevel" json:"def_level,omitempty"`
	SpeedLevel int32  `protobuf:"varint,15,opt,name=speed_level,json=speedLevel" json:"speed_level,omitempty"`
}

func (m *Player) Reset()                    { *m = Player{} }
func (m *Player) String() string            { return proto.CompactTextString(m) }
func (*Player) ProtoMessage()               {}
func (*Player) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Player) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Player) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Player) GetLevel() int32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *Player) GetFigure() int32 {
	if m != nil {
		return m.Figure
	}
	return 0
}

func (m *Player) GetWeaponType() int32 {
	if m != nil {
		return m.WeaponType
	}
	return 0
}

func (m *Player) GetScore() int32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *Player) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Player) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Player) GetLife() int32 {
	if m != nil {
		return m.Life
	}
	return 0
}

func (m *Player) GetCurLife() int32 {
	if m != nil {
		return m.CurLife
	}
	return 0
}

func (m *Player) GetDefLevel() int32 {
	if m != nil {
		return m.DefLevel
	}
	return 0
}

func (m *Player) GetSpeedLevel() int32 {
	if m != nil {
		return m.SpeedLevel
	}
	return 0
}

type BattleInfoNotify struct {
	InfoList []*BattleInfoNotify_PlayerInfo `protobuf:"bytes,1,rep,name=info_list,json=infoList" json:"info_list,omitempty"`
}

func (m *BattleInfoNotify) Reset()                    { *m = BattleInfoNotify{} }
func (m *BattleInfoNotify) String() string            { return proto.CompactTextString(m) }
func (*BattleInfoNotify) ProtoMessage()               {}
func (*BattleInfoNotify) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BattleInfoNotify) GetInfoList() []*BattleInfoNotify_PlayerInfo {
	if m != nil {
		return m.InfoList
	}
	return nil
}

type BattleInfoNotify_PlayerInfo struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	CurLife int32  `protobuf:"varint,2,opt,name=cur_life,json=curLife" json:"cur_life,omitempty"`
	X       int32  `protobuf:"varint,3,opt,name=x" json:"x,omitempty"`
	Y       int32  `protobuf:"varint,4,opt,name=y" json:"y,omitempty"`
}

func (m *BattleInfoNotify_PlayerInfo) Reset()                    { *m = BattleInfoNotify_PlayerInfo{} }
func (m *BattleInfoNotify_PlayerInfo) String() string            { return proto.CompactTextString(m) }
func (*BattleInfoNotify_PlayerInfo) ProtoMessage()               {}
func (*BattleInfoNotify_PlayerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *BattleInfoNotify_PlayerInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BattleInfoNotify_PlayerInfo) GetCurLife() int32 {
	if m != nil {
		return m.CurLife
	}
	return 0
}

func (m *BattleInfoNotify_PlayerInfo) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *BattleInfoNotify_PlayerInfo) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

type Property struct {
	Type Property_Type `protobuf:"varint,1,opt,name=type,enum=gamepb.Property_Type" json:"type,omitempty"`
	X    int32         `protobuf:"varint,2,opt,name=x" json:"x,omitempty"`
	Y    int32         `protobuf:"varint,3,opt,name=y" json:"y,omitempty"`
}

func (m *Property) Reset()                    { *m = Property{} }
func (m *Property) String() string            { return proto.CompactTextString(m) }
func (*Property) ProtoMessage()               {}
func (*Property) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Property) GetType() Property_Type {
	if m != nil {
		return m.Type
	}
	return Property_GEN
}

func (m *Property) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Property) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

type Scene struct {
	PropertyList []*Property `protobuf:"bytes,1,rep,name=property_list,json=propertyList" json:"property_list,omitempty"`
}

func (m *Scene) Reset()                    { *m = Scene{} }
func (m *Scene) String() string            { return proto.CompactTextString(m) }
func (*Scene) ProtoMessage()               {}
func (*Scene) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Scene) GetPropertyList() []*Property {
	if m != nil {
		return m.PropertyList
	}
	return nil
}

type SceneNotify struct {
	PropertyIncList []*Property `protobuf:"bytes,1,rep,name=property_inc_list,json=propertyIncList" json:"property_inc_list,omitempty"`
	PropertyDecList []*Property `protobuf:"bytes,2,rep,name=property_dec_list,json=propertyDecList" json:"property_dec_list,omitempty"`
}

func (m *SceneNotify) Reset()                    { *m = SceneNotify{} }
func (m *SceneNotify) String() string            { return proto.CompactTextString(m) }
func (*SceneNotify) ProtoMessage()               {}
func (*SceneNotify) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SceneNotify) GetPropertyIncList() []*Property {
	if m != nil {
		return m.PropertyIncList
	}
	return nil
}

func (m *SceneNotify) GetPropertyDecList() []*Property {
	if m != nil {
		return m.PropertyDecList
	}
	return nil
}

type EnterGameReq struct {
	RoomType int32 `protobuf:"varint,1,opt,name=room_type,json=roomType" json:"room_type,omitempty"`
}

func (m *EnterGameReq) Reset()                    { *m = EnterGameReq{} }
func (m *EnterGameReq) String() string            { return proto.CompactTextString(m) }
func (*EnterGameReq) ProtoMessage()               {}
func (*EnterGameReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *EnterGameReq) GetRoomType() int32 {
	if m != nil {
		return m.RoomType
	}
	return 0
}

type EnterGameAck struct {
	SceneInfo       *Scene    `protobuf:"bytes,1,opt,name=scene_info,json=sceneInfo" json:"scene_info,omitempty"`
	User            *Player   `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	OtherPlayerList []*Player `protobuf:"bytes,3,rep,name=other_player_list,json=otherPlayerList" json:"other_player_list,omitempty"`
}

func (m *EnterGameAck) Reset()                    { *m = EnterGameAck{} }
func (m *EnterGameAck) String() string            { return proto.CompactTextString(m) }
func (*EnterGameAck) ProtoMessage()               {}
func (*EnterGameAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *EnterGameAck) GetSceneInfo() *Scene {
	if m != nil {
		return m.SceneInfo
	}
	return nil
}

func (m *EnterGameAck) GetUser() *Player {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *EnterGameAck) GetOtherPlayerList() []*Player {
	if m != nil {
		return m.OtherPlayerList
	}
	return nil
}

func init() {
	proto.RegisterType((*Player)(nil), "gamepb.Player")
	proto.RegisterType((*BattleInfoNotify)(nil), "gamepb.BattleInfoNotify")
	proto.RegisterType((*BattleInfoNotify_PlayerInfo)(nil), "gamepb.BattleInfoNotify.PlayerInfo")
	proto.RegisterType((*Property)(nil), "gamepb.Property")
	proto.RegisterType((*Scene)(nil), "gamepb.Scene")
	proto.RegisterType((*SceneNotify)(nil), "gamepb.SceneNotify")
	proto.RegisterType((*EnterGameReq)(nil), "gamepb.EnterGameReq")
	proto.RegisterType((*EnterGameAck)(nil), "gamepb.EnterGameAck")
	proto.RegisterEnum("gamepb.Property_Type", Property_Type_name, Property_Type_value)
}

func init() { proto.RegisterFile("game/logic/gamepb/game.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 537 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x53, 0xed, 0x6e, 0xd3, 0x40,
	0x10, 0xc4, 0x5f, 0x69, 0xb2, 0x4e, 0x93, 0xf4, 0x04, 0xc8, 0x14, 0xa4, 0x46, 0xe6, 0x4f, 0x10,
	0x28, 0x95, 0x82, 0xf8, 0x83, 0x10, 0xa2, 0x88, 0xa8, 0x54, 0x8a, 0xaa, 0xca, 0xed, 0x7f, 0xcb,
	0xb5, 0xd7, 0xc1, 0xc2, 0xf1, 0x99, 0xf3, 0x05, 0xea, 0x47, 0xe0, 0x15, 0x10, 0x2f, 0xc0, 0x5b,
	0xa2, 0xdb, 0xb3, 0x43, 0x12, 0x24, 0xf8, 0x95, 0xdb, 0xd9, 0xd9, 0xf1, 0xce, 0xdc, 0x05, 0x9e,
	0x2c, 0xa3, 0x15, 0x9e, 0xe6, 0x7c, 0x99, 0xc5, 0xa7, 0xea, 0x58, 0xde, 0xd2, 0xcf, 0xb4, 0x14,
	0x5c, 0x72, 0xd6, 0xd1, 0x90, 0xff, 0xc3, 0x84, 0xce, 0x55, 0x1e, 0xd5, 0x28, 0xd8, 0x00, 0xcc,
	0x2c, 0xf1, 0x8c, 0xb1, 0x31, 0xe9, 0x05, 0x66, 0x96, 0x30, 0x06, 0x76, 0x11, 0xad, 0xd0, 0x33,
	0x09, 0xa1, 0x33, 0xbb, 0x0f, 0x4e, 0x8e, 0x5f, 0x31, 0xf7, 0xac, 0xb1, 0x31, 0x71, 0x02, 0x5d,
	0xb0, 0x87, 0xd0, 0x49, 0xb3, 0xe5, 0x5a, 0xa0, 0x67, 0x13, 0xdc, 0x54, 0xec, 0x04, 0xdc, 0x6f,
	0x18, 0x95, 0xbc, 0x08, 0x65, 0x5d, 0xa2, 0xe7, 0x50, 0x13, 0x34, 0x74, 0x53, 0x97, 0x24, 0x57,
	0xc5, 0x5c, 0xa0, 0xd7, 0xd1, 0x72, 0x54, 0xb0, 0x3e, 0x18, 0x77, 0x1e, 0x10, 0x62, 0xdc, 0xa9,
	0xaa, 0xf6, 0x5c, 0x5d, 0xd5, 0x6a, 0xa9, 0x3c, 0x4b, 0xd1, 0xeb, 0x13, 0x40, 0x67, 0xf6, 0x08,
	0xba, 0xf1, 0x5a, 0x84, 0x84, 0x1f, 0x12, 0x7e, 0x10, 0xaf, 0xc5, 0x42, 0xb5, 0x1e, 0x43, 0x2f,
	0xc1, 0x34, 0xd4, 0x3b, 0x0f, 0xa8, 0xd7, 0x4d, 0x30, 0x5d, 0xd0, 0xda, 0x27, 0xe0, 0x56, 0x25,
	0x62, 0xd2, 0xb4, 0x87, 0x7a, 0x3d, 0x82, 0x88, 0xe0, 0xff, 0x32, 0x60, 0xf4, 0x3e, 0x92, 0x32,
	0xc7, 0x8b, 0x22, 0xe5, 0x97, 0x5c, 0x66, 0x69, 0xcd, 0xde, 0x41, 0x2f, 0x2b, 0x52, 0x1e, 0xe6,
	0x59, 0x25, 0x3d, 0x63, 0x6c, 0x4d, 0xdc, 0xd9, 0xd3, 0xa9, 0x4e, 0x73, 0xba, 0x4f, 0x9e, 0xea,
	0x68, 0x15, 0x10, 0x74, 0xd5, 0xd4, 0x22, 0xab, 0xe4, 0xf1, 0x35, 0xc0, 0x1f, 0xfc, 0xaf, 0xd8,
	0xb7, 0xdd, 0x98, 0xbb, 0x6e, 0x28, 0x18, 0x6b, 0x27, 0x18, 0xbb, 0x09, 0xc6, 0x5f, 0x41, 0xf7,
	0x4a, 0xf0, 0x12, 0x85, 0xac, 0xd9, 0x33, 0xb0, 0x29, 0x70, 0x25, 0x3a, 0x98, 0x3d, 0x68, 0xb7,
	0x6b, 0xfb, 0x53, 0x95, 0x7d, 0x40, 0x14, 0x2d, 0x69, 0xee, 0x48, 0x5a, 0xad, 0xe4, 0x31, 0xd8,
	0x74, 0x4b, 0x07, 0x60, 0x9d, 0xcf, 0x2f, 0x47, 0xf7, 0x58, 0x0f, 0x9c, 0x8f, 0xf3, 0xb3, 0xe0,
	0x66, 0x64, 0xf8, 0x6f, 0xc1, 0xb9, 0x8e, 0xb1, 0x40, 0xf6, 0x0a, 0x0e, 0xcb, 0x46, 0x77, 0x3b,
	0x92, 0xd1, 0xfe, 0x47, 0x83, 0x7e, 0x4b, 0x53, 0x19, 0xf8, 0xdf, 0x0d, 0x70, 0x49, 0xa0, 0x49,
	0xf5, 0x0d, 0x1c, 0x6d, 0x64, 0xb2, 0x22, 0xfe, 0xb7, 0xd4, 0xb0, 0xa5, 0x5e, 0x14, 0xb1, 0x52,
	0xdb, 0x99, 0x4e, 0xb0, 0x99, 0x36, 0xff, 0x37, 0xfd, 0x01, 0x69, 0xda, 0x7f, 0x0e, 0xfd, 0x79,
	0x21, 0x51, 0x9c, 0x47, 0x2b, 0x0c, 0xf0, 0x8b, 0x7a, 0x34, 0x82, 0xf3, 0x55, 0xb8, 0xc9, 0xd0,
	0x09, 0xba, 0x0a, 0x50, 0x61, 0xf8, 0x3f, 0x8d, 0x2d, 0xf6, 0x59, 0xfc, 0x99, 0xbd, 0x00, 0xa8,
	0x94, 0x91, 0x50, 0xdd, 0x2f, 0xd1, 0xdd, 0xd9, 0x61, 0xfb, 0x51, 0xb2, 0x18, 0xf4, 0x88, 0x40,
	0xb7, 0xed, 0x83, 0xbd, 0xae, 0x50, 0x50, 0xe4, 0xee, 0x6c, 0xb0, 0x59, 0x8e, 0xde, 0x43, 0x40,
	0x3d, 0xf6, 0x1a, 0x8e, 0xb8, 0xfc, 0x84, 0x22, 0x2c, 0x09, 0xd5, 0x6e, 0x2c, 0x72, 0xb3, 0x3f,
	0x30, 0x24, 0xa2, 0x2e, 0x94, 0x97, 0xdb, 0x0e, 0xfd, 0xbd, 0x5f, 0xfe, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0xd5, 0xc7, 0x51, 0x37, 0xfe, 0x03, 0x00, 0x00,
}