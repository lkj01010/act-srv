// Code generated by protoc-gen-go.
// source: pb/rpc.proto
// DO NOT EDIT!

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Game_FrameType int32

const (
	Game_Message Game_FrameType = 0
	Game_Kick    Game_FrameType = 1
	Game_Ping    Game_FrameType = 2
)

var Game_FrameType_name = map[int32]string{
	0: "Message",
	1: "Kick",
	2: "Ping",
}
var Game_FrameType_value = map[string]int32{
	"Message": 0,
	"Kick":    1,
	"Ping":    2,
}

func (x Game_FrameType) String() string {
	return proto.EnumName(Game_FrameType_name, int32(x))
}
func (Game_FrameType) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

type Game struct {
}

func (m *Game) Reset()                    { *m = Game{} }
func (m *Game) String() string            { return proto.CompactTextString(m) }
func (*Game) ProtoMessage()               {}
func (*Game) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

type Game_Frame struct {
	Type    Game_FrameType `protobuf:"varint,1,opt,name=Type,enum=pb.Game_FrameType" json:"Type,omitempty"`
	Message []byte         `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (m *Game_Frame) Reset()                    { *m = Game_Frame{} }
func (m *Game_Frame) String() string            { return proto.CompactTextString(m) }
func (*Game_Frame) ProtoMessage()               {}
func (*Game_Frame) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

func (m *Game_Frame) GetType() Game_FrameType {
	if m != nil {
		return m.Type
	}
	return Game_Message
}

func (m *Game_Frame) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func init() {
	proto.RegisterType((*Game)(nil), "pb.Game")
	proto.RegisterType((*Game_Frame)(nil), "pb.Game.Frame")
	proto.RegisterEnum("pb.Game_FrameType", Game_FrameType_name, Game_FrameType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GameService service

type GameServiceClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (GameService_StreamClient, error)
}

type gameServiceClient struct {
	cc *grpc.ClientConn
}

func NewGameServiceClient(cc *grpc.ClientConn) GameServiceClient {
	return &gameServiceClient{cc}
}

func (c *gameServiceClient) Stream(ctx context.Context, opts ...grpc.CallOption) (GameService_StreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GameService_serviceDesc.Streams[0], c.cc, "/pb.GameService/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &gameServiceStreamClient{stream}
	return x, nil
}

type GameService_StreamClient interface {
	Send(*Game_Frame) error
	Recv() (*Game_Frame, error)
	grpc.ClientStream
}

type gameServiceStreamClient struct {
	grpc.ClientStream
}

func (x *gameServiceStreamClient) Send(m *Game_Frame) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gameServiceStreamClient) Recv() (*Game_Frame, error) {
	m := new(Game_Frame)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for GameService service

type GameServiceServer interface {
	Stream(GameService_StreamServer) error
}

func RegisterGameServiceServer(s *grpc.Server, srv GameServiceServer) {
	s.RegisterService(&_GameService_serviceDesc, srv)
}

func _GameService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GameServiceServer).Stream(&gameServiceStreamServer{stream})
}

type GameService_StreamServer interface {
	Send(*Game_Frame) error
	Recv() (*Game_Frame, error)
	grpc.ServerStream
}

type gameServiceStreamServer struct {
	grpc.ServerStream
}

func (x *gameServiceStreamServer) Send(m *Game_Frame) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gameServiceStreamServer) Recv() (*Game_Frame, error) {
	m := new(Game_Frame)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _GameService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.GameService",
	HandlerType: (*GameServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _GameService_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pb/rpc.proto",
}

func init() { proto.RegisterFile("pb/rpc.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 178 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x48, 0xd2, 0x2f,
	0x2a, 0x48, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0xaa, 0xe7, 0x62,
	0x71, 0x4f, 0xcc, 0x4d, 0x95, 0xf2, 0xe4, 0x62, 0x75, 0x2b, 0x4a, 0xcc, 0x4d, 0x15, 0x52, 0xe3,
	0x62, 0x09, 0xa9, 0x2c, 0x48, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x33, 0x12, 0xd2, 0x2b, 0x48,
	0xd2, 0x03, 0x29, 0xd0, 0x03, 0xcb, 0x82, 0x64, 0x82, 0xc0, 0xf2, 0x42, 0x12, 0x5c, 0xec, 0xbe,
	0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x30, 0xae, 0x92,
	0x0e, 0x17, 0x27, 0x5c, 0xb1, 0x10, 0x37, 0x5c, 0x99, 0x00, 0x83, 0x10, 0x07, 0x17, 0x8b, 0x77,
	0x66, 0x72, 0xb6, 0x00, 0x23, 0x88, 0x15, 0x90, 0x99, 0x97, 0x2e, 0xc0, 0x64, 0x64, 0xcd, 0xc5,
	0x0d, 0x32, 0x3f, 0x38, 0xb5, 0xa8, 0x2c, 0x33, 0x39, 0x55, 0x48, 0x87, 0x8b, 0x2d, 0xb8, 0xa4,
	0x28, 0x35, 0x31, 0x57, 0x88, 0x0f, 0xd5, 0x6a, 0x29, 0x34, 0xbe, 0x06, 0xa3, 0x01, 0x63, 0x12,
	0x1b, 0xd8, 0x23, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x84, 0xb0, 0x3f, 0x7d, 0xd8, 0x00,
	0x00, 0x00,
}
