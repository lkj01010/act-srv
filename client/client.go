package client

import (
	"net"
	"time"
	"encoding/binary"
	"io"
	"github.com/lkj01010/goutils/log"

	"github.com/lkj01010/act-srv/consts"
	. "github.com/lkj01010/act-srv/com"
	agentLogic "github.com/lkj01010/act-srv/agent/logic"
	"github.com/lkj01010/goutils/timer"
	"github.com/lkj01010/act-srv/misc/packet"
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/act-srv/game/core/gamepb"
	"sync"
)

const (
	READ_DEADLINE = 120
	MSG_SEND_INTERVAL = 1
)

type Client struct {
	conn net.Conn
	sess session
}

func NewClient() (*Client, error) {
	c, err := net.Dial("tcp", consts.AgentPort)
	if err != nil {
		log.Error("client dail error=", err)
	}
	cli := &Client{
		conn: c,
	}
	return cli, nil
}

func (cli *Client) Close() {
	cli.conn.Close()
	log.Debug("agent client close")
}

func (cli *Client) Startup(wg *sync.WaitGroup) {
	conn := cli.conn

	// for reading the 2-Byte header
	header := make([]byte, 2)

	// the input channel for agent()
	in := make(chan []byte)
	defer func() {
		close(in) // session will close
		wg.Done()
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
	msgPool msgPool
}

func NewSession(in chan []byte, w io.Writer) *session {
	ss := &session{
		in:    in,
		w:     w,
		die:   make(chan struct{}),
		timer: timer.NewTimer(),
		msgPool: msgPool{
			cmdList:     make([]cmdData, 0),
			curCmdIndex: 0,
		},
	}
	ss.msgPool.init()
	return ss
}

func (ss *session) serve() {
	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, os.Interrupt, os.Kill)

	ch := make(chan []byte, 1)

	defer func() {
		close(ch)
		close(ss.die)
	}()

	timer := time.NewTimer(0)
	for {
		select {
		case msg, ok := <-ss.in:
			if !ok {
				return
			}

			if cmd, payload := ss.proxy_msg(msg); cmd > 0 {
				//ss.Send(result)
				ss.Send(cmd, payload)
			}
		case <-timer.C:
			timer.Reset(MSG_SEND_INTERVAL * time.Second)
			ss.update() // game loop
		//case sig := <-sigCh:
		//	if sig == os.Interrupt || sig == os.Kill {
		//		log.Info("agent shutdown")
		//		//os.Exit(0)
		//	}
		}
	}
}

func (ss *session) Send(cmd int16, payload []byte) error {
	ss.seqId++
	s_agent := &agentLogic.S_agent{
		F_seqId:   ss.seqId,
		F_proto:   cmd,
		F_payload: payload,
	}
	dataByte := s_agent.Pack()
	log.Debugf("[send][cmd=%+v][payload=%+v]", RCmd[cmd], payload)
	if _, err := ss.w.Write(dataByte); err != nil {
		log.Error("session send dail error=", err)
		return err
	}
	return nil
}

func (ss *session) proxy_msg(msg []byte) (int16, []byte) {
	//log.Debugf("proxy_msg : %+v", msg)

	reader := packet.Reader(msg)
	var cmd int16
	var payload []byte
	var ret []byte
	var err error
	if cmd, err = reader.ReadS16(); err != nil {
		panic(err)
	}
	if payload, err = reader.ReadBytes(); err != nil {
		panic(err)
	}
	log.Debugf("[receive][cmd=%+v][payload=%+v]", RCmd[cmd], payload)

	if ret, err = handlers[cmd](ss, payload); err != nil {
		panic(err)
	}

	log.Debugf("[handle ret=%+v]", ret)

	//data := pool.getCurCmdData()
	//return data.cmd, data.payload
	return 0, nil
}

func (ss *session) update() {
	data := ss.msgPool.getCurCmdData()
	if data != nil {
		ss.Send(data.cmd, data.payload)
	}
}

//////////////////////////////////////////////////
type cmdData struct {
	cmd     int16
	payload []byte
}

type msgPool struct {
	cmdList     []cmdData
	curCmdIndex int
}

func (p *msgPool) init() {
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     Cmd[Agent_LoginReq],
		payload: nil,
	})
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     Cmd[Agent_HeartbeatReq],
		payload: nil,
	})
	//p.cmdList = append(p.cmdList, cmdData{
	//	cmd:     agentLogic.Cmd[agentLogic.HeartbeatReq],
	//	payload: nil,
	//})

	payload, _ := proto.Marshal(&gamepb.EnterGameReq{
		RoomType: 1,
		Figure:   8,
	})
	//writer := packet.Writer()
	//writer.WriteBytes(payload)
	p.cmdList = append(p.cmdList, cmdData{
		cmd: Cmd[Game_EnterGameReq],
		//payload: writer.Data(),
		payload: payload,
	})
}

func (p *msgPool) getCurCmdData() (data *cmdData) {
	if p.curCmdIndex < len(p.cmdList) {
		data = &p.cmdList[p.curCmdIndex]
	} else {
		data = nil
	}
	p.curCmdIndex++
	return
}

//////////////////////////////////////////////////
func checkErr(err error) {
	if err != nil {
		panic("error occured in protocol module")
	}
}
