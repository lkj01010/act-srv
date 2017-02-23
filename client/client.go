package client

import (
	"net"
	"time"
	"encoding/binary"
	"io"
	"github.com/lkj01010/goutils/log"

	"github.com/lkj01010/act-srv/consts"
	agentLogic "github.com/lkj01010/act-srv/agent/logic"
	gameLogic "github.com/lkj01010/act-srv/game/logic"
	"github.com/lkj01010/goutils/timer"
	"github.com/lkj01010/act-srv/misc/packet"
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/act-srv/game/logic/gamepb"
	"os"
	"os/signal"
)

const (
	READ_DEADLINE = 120
)

var pool cmdDataPool

func init() {
	pool = cmdDataPool{
		cmdList:     make([]cmdData, 0),
		curCmdIndex: 0,
	}
	pool.init()
}

type Client struct {
	conn net.Conn
	sess session
}

func NewClient() (*Client, error) {
	c, err := net.Dial("tcp", consts.AgentPort)
	if err != nil {
		log.Error("client dail error=", err)
	}
	return &Client{
		conn: c,
	}, err
}

func (cli *Client) Close() {
	cli.conn.Close()
	log.Debug("agent client close")
}

func (cli *Client) Startup() {
	conn := cli.conn

	// for reading the 2-Byte header
	header := make([]byte, 2)

	// the input channel for agent()
	in := make(chan []byte)
	defer func() {
		close(in) // session will close
	}()

	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("cannot get server address:", err)
		return
	}
	IP := net.ParseIP(host)
	log.Infof("connect server from:%v port:%v", host, port)

	// session
	sess := NewSession(in, conn)
	go sess.serve()

	// read loop
	for {
		conn.SetReadDeadline(time.Now().Add(READ_DEADLINE * time.Second))

		// read 2B header
		n, err := io.ReadFull(conn, header)
		if err != nil {
			log.Warningf("read header failed, ip:%v reason:%v size:%v", IP, err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)

		//log.Debugf("read header=%+v, size=%+v", header, size)

		// alloc a byte slice of the size defined in the header for reading data
		payload := make([]byte, size)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			log.Warningf("read payload failed, ip:%v reason:%v size:%v", IP, err, n)
			return
		}

		// deliver the data to the input queue of agent()
		select {
		case in <- payload: // payload queued
		case <-sess.die:
			log.Warningf("connection closed by logic, ip:%v", IP)
			return
		}
	}
}

//////////////////////////////////////////////////

type session struct {
	seqId uint32
	in    chan []byte
	die   chan struct{}
	w     io.Writer

	timer timer.Timer
}

func NewSession(in chan []byte, w io.Writer) *session {
	return &session{
		in:    in,
		w:     w,
		die:   make(chan struct{}),
		timer: timer.NewTimer(),
	}
}

func (sess *session) serve() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)

	ch := make(chan []byte, 1)

	defer func() {
		close(ch)
		close(sess.die)
	}()

	timer := time.NewTimer(0)
	for {
		select {
		case msg, ok := <-sess.in:
			if !ok {
				return
			}

			if cmd, payload := sess.proxy_msg(msg); cmd > 0 {
				//sess.Send(result)
				sess.Send(cmd, payload)
			}
		case <-timer.C:
			timer.Reset(5 * time.Second)
			sess.update() // game loop
		case sig := <-sigCh:
			if sig == os.Interrupt || sig == os.Kill {
				log.Info("agent shutdown")
				os.Exit(0)
			}
		}
	}
}

func (sess *session) Send(cmd int16, payload []byte) error {
	sess.seqId++
	s_agent := &agentLogic.S_agent{
		F_seqId:   sess.seqId,
		F_proto:   cmd,
		F_payload: payload,
	}
	dataByte := s_agent.Pack()
	log.Debugf("send cmd=%+v, payload=%+v", cmd, payload)
	if _, err := sess.w.Write(dataByte); err != nil {
		log.Error("session send dail error=", err)
		return err
	}
	return nil
}

func (sess *session) proxy_msg(msg []byte) (int16, []byte) {
	//log.Debugf("proxy_msg : %+v", msg)

	reader := packet.Reader(msg)
	var cmd int16
	var payload []byte
	var ret []byte
	var err error
	if cmd, err = reader.ReadS16(); err != nil {
		panic(err)
	}
	//log.Debugf("proxy msg, cmd=%+v", cmd)
	if payload, err = reader.ReadBytes(); err != nil {
		panic(err)
	}
	//log.Debugf("proxy msg, payload=%+v", payload)

	if ret, err = handlers[cmd](sess, payload); err != nil {
		panic(err)
	}

	log.Debugf("proxy_msg msg=%+v, ret=%+v", msg, ret)

	//data := pool.getNext()
	//return data.cmd, data.payload
	return 0, nil
}

func (sess *session) update() {
	data := pool.getNext()
	sess.Send(data.cmd, data.payload)
}

//////////////////////////////////////////////////
type cmdData struct {
	cmd     int16
	payload []byte
}

type cmdDataPool struct {
	cmdList     []cmdData
	curCmdIndex int
}

func (p *cmdDataPool) init() {
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     agentLogic.Cmd["login_req"],
		payload: nil,
	})
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     agentLogic.Cmd["heartbeat_req"],
		payload: nil,
	})
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     agentLogic.Cmd["heartbeat_req"],
		payload: nil,
	})

	payload, _ := proto.Marshal(&gamepb.EnterGameReq{
		RoomType: 1,
		Figure: 8,
	})
	//writer := packet.Writer()
	//writer.WriteBytes(payload)
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     gameLogic.Cmd["enter_game_req"],
		//payload: writer.Data(),
		payload: payload,
	})
}

func (p *cmdDataPool) getNext() (data cmdData) {
	data = p.cmdList[p.curCmdIndex]
	//log.Infof("getCmdIndex=%+v", p.curCmdIndex)

	if p.curCmdIndex < len(p.cmdList) - 1 {
		p.curCmdIndex++
	} else {
		p.curCmdIndex = 1
	}

	return
}

//////////////////////////////////////////////////
func checkErr(err error) {
	if err != nil {
		panic("error occured in protocol module")
	}
}