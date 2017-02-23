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
	. "github.com/lkj01010/act-srv/game/logic"
	"github.com/lkj01010/act-srv/game/registry"
	. "github.com/lkj01010/act-srv/game/types"
	"github.com/lkj01010/act-srv/misc/packet"
	. "github.com/lkj01010/act-srv/pb"
	"github.com/lkj01010/act-srv/services"
	. "github.com/lkj01010/act-srv/utils"
)

var (
	ERROR_INCORRECT_FRAME_TYPE = errors.New("incorrect frame type")
	ERROR_SERVICE_NOT_BIND     = errors.New("service not bind")
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
func (s *server) recv(stream GameService_StreamServer, sess_die chan struct{}) <-chan *Game_Frame {
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

			log.Infof("Recv in=%+v, err=%+v", in, err)
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case ch <- in:
			case <-sess_die: // 关闭的channel可以立即取出数据,应该是nil
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
	var sess Session
	sess_die := make(chan struct{})
	ch_agent := s.recv(stream, sess_die)
	ch_ipc := make(chan *Game_Frame, DEFAULT_CH_IPC_SIZE)

	defer func() {
		registry.Unregister(sess.UserId)
		close(sess_die)
		log.Debug("stream end:", sess.UserId)
	}()

	// read metadata from context
	md, ok := metadata.FromContext(stream.Context())
	if !ok {
		log.Error("cannot read metadata from context")
		return ERROR_INCORRECT_FRAME_TYPE
	}
	// read key
	if len(md["userid"]) == 0 {
		log.Error("cannot read key:userid from metadata")
		return ERROR_INCORRECT_FRAME_TYPE
	}
	// parse userid
	userid, err := strconv.Atoi(md["userid"][0])
	if err != nil {
		log.Error(err)
		return ERROR_INCORRECT_FRAME_TYPE
	}

	// register user
	sess.UserId = int32(userid)
	registry.Register(sess.UserId, ch_ipc)
	log.Debug("userid", sess.UserId, "logged in")

	// >> main message loop <<
	for {
		select {
		case frame, ok := <-ch_agent:
			// frames from agent
			if !ok {
				// EOF
				return nil
			}
			switch frame.Type {
			case Game_Message: // the passthrough message from client->agent->game
				// locate handler by proto number
				reader := packet.Reader(frame.Message)
				c, err := reader.ReadS16()
				if err != nil {
					log.Error(err)
					return err
				}
				handle := Handlers[c]
				if handle == nil {
					log.Error("service not bind:", c)
					return ERROR_SERVICE_NOT_BIND

				}

				// handle request
				//ret := handle(&sess, reader)
				payload, err := reader.ReadBytes()
				if err != nil {
					log.Error(err)
					return err
				}
				ret := handle(&sess, payload)

				// construct frame & return message from logic
				if ret != nil {
					if err := stream.Send(&Game_Frame{Type: Game_Message, Message: ret}); err != nil {
						log.Error(err)
						return err
					}
				}

				// session control by logic
				if sess.Flag&SESS_KICKED_OUT != 0 {
					// logic kick out
					if err := stream.Send(&Game_Frame{Type: Game_Kick}); err != nil {
						log.Error(err)
						return err
					}
					return nil
				}
			case Game_Ping:
				if err := stream.Send(&Game_Frame{Type: Game_Ping, Message: frame.Message}); err != nil {
					log.Error(err)
					return err
				}
				log.Debug("pinged")
			default:
				log.Error("incorrect frame type:", frame.Type)
				return ERROR_INCORRECT_FRAME_TYPE
			}
		case frame := <-ch_ipc:
			// forward async messages from interprocess(goroutines) communication
			if err := stream.Send(frame); err != nil {
				log.Error(err)
				return err
			}
		}
	}
}
