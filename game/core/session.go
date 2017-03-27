package core

import (
	"google.golang.org/grpc/metadata"
	"strconv"
	"io"
	"errors"
	"github.com/lkj01010/goutils/log"

	. "github.com/lkj01010/act-srv/utils"
	. "github.com/lkj01010/act-srv/pb"
	"github.com/lkj01010/act-srv/game/registry"
	"github.com/lkj01010/act-srv/misc/packet"
	. "github.com/lkj01010/act-srv/com"
)

var (
	errIncorrectFrameType  = errors.New("incorrect frame type")
	ERROR_SERVICE_NOT_BIND = errors.New("service not bind")
)

//type ssFlagType int32
//const (
//	ssKickedOut ssFlagType = iota
//	ssInGame
//	ssFlagLen
//)
//
//type Msg struct {
//	Cmd     int16
//	Payload []byte
//}
//
//type MsgWithSession struct {
//	Ss 	*Session
//	Msg
//}

type HandleFunc interface{}

type Session struct {
	ToGameCh  chan HandleFunc
	DieCh     chan struct{}

	isDie    int32
	isInGame int32

	fnCh   chan interface{}
	stream GameService_StreamServer

	UserId int32
}

//func (ss *Session)StoreFlag(bit ssFlagType, value int32) {
//	atomic.StoreInt32(&ss.Flags[bit], value)
//}
//
//func (ss *Session)LoadFlag(bit ssFlagType, value int32) {
//	atomic.LoadInt32(&ss.Flags[bit])
//}

func NewSession() *Session {
	return &Session{
		DieCh:     make(chan struct{}),
		fnCh:      make(chan interface{}, SESSION_FN_CH_SIZE),
	}
}

func (ss *Session) GetUserId() int32 {
	return ss.UserId
}

// PIPELINE #1 stream receiver
// this function is to make the stream receiving SELECTABLE
func (ss *Session) recv(stream GameService_StreamServer, ssDie <-chan struct{}) <-chan *Game_Frame {
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
func (ss *Session) Stream(stream GameService_StreamServer) error {
	defer PrintPanicStack()

	ssDie := make(chan struct{})
	streamCh := ss.recv(stream, ssDie)

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
	Ipc.Register(ss.UserId, ss.fnCh)
	log.Debugf("[stream open][userId=%+v]", ss.UserId)
	ss.stream = stream

	// >> main message loop <<
	for {
		select {
		case frame, ok := <-streamCh:
			// frames from agent
			if ok {
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

					log.Debugf("[dispatch][cmd=%+v][payload=%+v][userId=%+v]", Cmd(c).String(), payload,
						ss.UserId)
					ss.dispatch(Cmd(c), payload)

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
			} else {
				// EOF
				ss.Close()
			}

		case fn := <-ss.fnCh:
			h := fn.(func(g *Session))
			h(ss)

		//case msg := <-ss.ToAgentCh:
		//	if err := stream.Send(&Game_Frame{Type: Game_Message, Message: msg}); err != nil {
		//		log.Error(err)
		//		return err
		//	}
		case <-ss.DieCh:
			if err := stream.Send(&Game_Frame{Type: Game_Kick}); err != nil {
				log.Error(err)
				return err
			}
			return nil
		}
	}
}

func (ss *Session) H_enterGame(gameCh chan HandleFunc) {
	if ss.ToGameCh != nil {
		log.Panic("never reach")
	}
	ss.ToGameCh = gameCh
	// 逻辑上对agent的ack交给game
	log.Debugf("[game <- H_sessionEnterAck][userId=%+v]", ss.UserId)
	ss.ToGameCh <- func(g *game) {
		g.H_sessionEnterAck(ss)
	}
}

func (ss *Session) Close() {
	if ss.ToGameCh != nil {
		log.Debugf("[game <- H_sessionLeave][userId=%+v]", ss.UserId)
		ss.ToGameCh <- func(g *game) {
			g.H_sessionLeave(ss)
		}
		ss.ToGameCh = nil
	} else {
		ss.Die()
	}
}

func (ss *Session) H_close() {
	//log.Debugf("[session H_close()][userId=%+v][ToGameCh=%+v]", ss.UserId, ss.ToGameCh)
	if ss.ToGameCh != nil {
		log.Debugf("[game <- H_sessionCloseAck][userId=%+v]", ss.UserId)
		ss.ToGameCh <- func(g *game) {
			g.H_sessionCloseAck(ss)
		}
		ss.ToGameCh = nil
		//log.Debugf("[session H_close() done][userId=%+v][ToGameCh=%+v]", ss.UserId, ss.ToGameCh)
	} else {
		ss.Die()
	}
}

func (ss *Session) Die() {
	//log.Debugf("[session Die()][userId=%+v", ss.UserId)
	close(ss.DieCh)
}

func (ss *Session) H_streamSend(msg []byte) {
	ss.stream.Send(&Game_Frame{Type: Game_Message, Message: msg})
}

func (ss *Session) dispatch(cmd Cmd, payload []byte) {
	dest := CmdSendTo(cmd)
	if dest == DestGameMgr {
		GameMgr.fnCh <- gameMgrDecoders[cmd](ss, payload)
	} else if dest == DestGame {
		if ss.ToGameCh != nil {
			ss.ToGameCh <- gameDecoders[cmd](ss, payload)
		} else {
			log.Error("dispatch to game when ToGameCh==nil")
		}
	} else {
		log.Errorf("[cmd not found][cmd=%+v]", cmd)
	}
}

//////////////////////////////////////////////////

var Ipc registry.Registry

func init() {
	Ipc.Init()
}
