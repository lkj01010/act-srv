package client

import (
	"net"
	"github.com/lkj01010/goutils/log"

	"github.com/lkj01010/act-srv/consts"
	"github.com/lkj01010/act-srv/agent/logic"
	"time"
	"encoding/binary"
	"io"
	"github.com/lkj01010/goutils/timer"
	"github.com/lkj01010/act-srv/misc/packet"
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

func (cli *Client) Send(bytes []byte) {
	if _, err := cli.conn.Write(bytes); err != nil {
		log.Error(err)
	}
	log.Debugf("send bytes: %v", bytes)
}

func (cli *Client) Close() {
	cli.conn.Close()
	log.Debug("agent client close")
}

func (cli *Client) Startup() {
	log.Infof("startup p: %+v", pool);
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

		log.Debugf("read header: %+v", header)

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
	timeStamp int32
	in        chan []byte
	die       chan struct{}
	w         io.Writer

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
	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGTERM)

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
			timer.Reset(3 * time.Second)
			sess.update() // game loop
		//case sig := <-sigCh:
		//	if sig == syscall.SIGTERM {
		//		log.Info("agent shutdown.")
		//		os.Exit(0)
		//	}
		}
	}
}

func (sess *session) Send(cmd int16, payload []byte) error {
	sess.timeStamp++
	s_agent := &logic.S_agent{
		F_timeStamp: sess.timeStamp,
		F_proto:     cmd,
		F_payload:   payload,
	}
	dataByte := s_agent.Pack()
	if _, err := sess.w.Write(dataByte); err != nil {
		log.Error("session send dail error=", err)
		return err
	}
	return nil
}

func (sess *session) proxy_msg(msg []byte) (int16, []byte) {
	log.Debugf("proxy msg: %v", msg)

	//return []byte("123")
	reader := packet.Reader(msg)
	var cmd int16
	var payload []byte
	var ret []byte
	var err error
	if cmd, err = reader.ReadS16(); err != nil {
		panic(err)
	}
	log.Debugf("prox msg, cmd=%+v", cmd)
	if payload, err = reader.ReadBytes(); err != nil {
		panic(err)
	}
	log.Debugf("proxy msg, payload", payload)

	if ret, err = handlers[cmd](sess, payload); err != nil {
		panic(err)
	}

	log.Debugf("handle ret=%+v", ret)

	data := pool.getNext()
	return data.cmd, data.payload
}

func (sess *session) update() {
	data := pool.getNext()
	sess.Send(data.cmd, data.payload)
}

func (s *session) Req_login() []byte {

	s_agent := &logic.S_agent{
		F_timeStamp: s.timeStamp,
		F_proto:     logic.Cmd["login_req"],
		F_payload:   nil,
	}
	return s_agent.Pack()
	//log.Debugf("send: %+v", s_agent)
}

func (s *session) Req_enter_game() (b []byte) {
	return
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
		cmd:     logic.Cmd["login_req"],
		payload: nil,
	})
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     logic.Cmd["heart_beat_req"],
		payload: nil,
	})
	p.cmdList = append(p.cmdList, cmdData{
		cmd:     logic.Cmd["heart_beat_req"],
		payload: nil,
	})
	log.Infof("p init: %+v", p);
}

func (p *cmdDataPool) getNext() (data cmdData) {
	//log.Infof("p getNext: %+v", p);
	if p.curCmdIndex < len(p.cmdList) - 1 {
		p.curCmdIndex++
	} else {
		p.curCmdIndex = 1
	}
	data = p.cmdList[p.curCmdIndex]

	//data = cmdData{
	//	cmd: logic.Cmd["login_req"],
	//	payload: nil,
	//}
	return
}

//////////////////////////////////////////////////
func checkErr(err error) {
	if err != nil {
		panic("error occured in protocol module")
	}
}
