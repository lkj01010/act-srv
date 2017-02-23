package logic

import (
	. "github.com/lkj01010/act-srv/agent/types"
)

var Cmd = map[string]int16{
	"heartbeat_req":         1, // 心跳包..
	"heartbeat_ack":         2, // 心跳包回复
	"login_req":            3, // 登陆
	"login_ack":            4, // 登陆回执
}

var RCmd = map[int16]string{}
var Handlers map[int16]func(*Session) ([]byte, error)
var UserId int32

func init() {
	for k, v := range Cmd {
		RCmd[v] = k
	}

	Handlers = map[int16]func(*Session) ([]byte, error) {
		Cmd["heartbeat_req"]: H_heartbeat_req,
		Cmd["login_req"]: H_login_req,
	}
	UserId = 0
}

func genUserId() (id int32) {
	UserId++
	id = UserId
	return
}


