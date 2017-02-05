package agent

import (
	"net"
	"github.com/lkj01010/goutils/log"

	"github.com/lkj01010/act-srv/consts"
	"github.com/lkj01010/act-srv/agent/logic"
)

type Client struct {
	conn net.Conn
	timeStamp int32
}

func NewClient() (*Client, error) {
	c, err := net.Dial("tcp", consts.AgentPort)
	if err != nil {
		log.Error("client dail error=", err)
	}
	return &Client{
		conn: c,
		timeStamp: 0,
	}, err
}

func (cli *Client)Send(bytes []byte) {
	if _, err := cli.conn.Write(bytes); err != nil {
		log.Error(err)
	}
	log.Debugf("send bytes: %v", bytes)
}

func (cli *Client)Close() {
	cli.conn.Close()
	log.Debug("agent client close")
}

func (cli *Client)Req_login() {
	cli.timeStamp++

	s_agent := &logic.S_agent{
		F_timeStamp: cli.timeStamp,
		F_proto: logic.Cmd["login_req"],
		F_payload: nil,
	}
	cli.Send(s_agent.Pack())
	//log.Debugf("send: %+v", s_agent)
}