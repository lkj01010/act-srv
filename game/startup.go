package game

import (
	"errors"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/lkj01010/goutils/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	. "github.com/lkj01010/act-srv/consts"
	"github.com/lkj01010/act-srv/misc/packet"
	. "github.com/lkj01010/act-srv/pb"
	"github.com/lkj01010/act-srv/services"
	. "github.com/lkj01010/act-srv/utils"
	. "github.com/lkj01010/act-srv/game/com"
	"github.com/lkj01010/act-srv/game/core"
)

var (
	errIncorrectFrameType  = errors.New("incorrect frame type")
	ERROR_SERVICE_NOT_BIND = errors.New("service not bind")
)

func Startup() {
	// 监听
	lis, err := net.Listen("tcp", GamePort)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	log.Info("gamesrv listening on ", lis.Addr())

	// 注册服务
	s := grpc.NewServer()
	ins := new(server)
	RegisterGameServiceServer(s, ins)

	//// 初始化Services
	//sp.Init("snowflake")
	// 开始服务
	s.Serve(lis)

	services.CloseGameConn()
}

type server struct{}

// PIPELINE #1 stream receiver
// this function is to make the stream receiving SELECTABLE
func (s *server) recv(stream GameService_StreamServer, ssDie chan struct{}) <-chan *Game_Frame {
	ch := make(chan *Game_Frame, 1)
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// client closed
				return
			}

			//log.Infof("Recv in=%+v, err=%+v", in, err)
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case ch <- in:
			case <-ssDie: // 关闭的channel可以立即取出数据,应该是nil
				// mid:
				return
			}
		}
	}()
	return ch
}

// PIPELINE #2 stream processing
// the center of game logic
func (s *server) Stream(stream GameService_StreamServer) error {
	defer PrintPanicStack()
	// session init
	ss := NewSession()

	ssDie := make(chan struct{})
	streamCh := s.recv(stream, ssDie)

	defer func() {
		Ipc.Unregister(ss.UserId)
		close(ssDie)
		log.Debugf("[stream end][userid=%+v]", ss.UserId)
	}()

	// read metadata from context
	md, ok := metadata.FromContext(stream.Context())
	if !ok {
		log.Error("cannot read metadata from context")
		return errIncorrectFrameType
	}
	// read key
	if len(md["userid"]) == 0 {
		log.Error("cannot read key:userid from metadata")
		return errIncorrectFrameType
	}
	// parse userid
	userId, err := strconv.Atoi(md["userid"][0])
	if err != nil {
		log.Error(err)
		return errIncorrectFrameType
	}

	// register user
	ss.UserId = int32(userId)
	Ipc.Register(ss.UserId, ss.ToAgentCh)
	log.Debugf("[stream open][userId=%+v]", ss.UserId)

	// >> main message loop <<
	for {
		select {
		case frame, ok := <-streamCh:
			// frames from agent
			if !ok {
				// EOF
				ss.Flag |= SessKickedOut
				log.Debug("streamCh is closed")
				return nil
			}
			switch frame.Type {
			case Game_Message:
				// locate handler by proto number
				reader := packet.Reader(frame.Message)
				c, err := reader.ReadS16()
				if err != nil {
					log.Error(err)
					return err
				}

				payload, err := reader.ReadBytes()
				if err != nil {
					log.Error(err)
					return err
				}

				log.Debugf("[dispatch][cmd=%v][payload=%+v]", c, payload)
				core.Dispatch(ss, c, payload)

			case Game_Ping:
				if err := stream.Send(&Game_Frame{Type: Game_Ping, Message: frame.Message}); err != nil {
					log.Error(err)
					return err
				}
				log.Debug("pinged")
			default:
				log.Errorf("[incorrect frame type=%+v]", frame.Type)
				return errIncorrectFrameType
			}

		case msg := <-ss.ToAgentCh:
			if err := stream.Send(&Game_Frame{Type: Game_Message, Message: msg}); err != nil {
				log.Error(err)
				return err
			}
		case <-ss.Die:
			if err := stream.Send(&Game_Frame{Type: Game_Kick}); err != nil {
				log.Error(err)
				return err
			}
			return nil
		}
	}
}
