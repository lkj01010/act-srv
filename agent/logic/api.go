package logic

import (
	. "github.com/lkj01010/act-srv/agent/types"
	. "github.com/lkj01010/act-srv/com"
)

var Handlers map[Cmd]func(*Session) ([]byte, error)
var UserId int32

func init() {
	Handlers = map[Cmd]func(*Session) ([]byte, error) {
		Agent_HeartbeatReq: H_heartbeat_req,
		Agent_LoginReq: H_login_req,
	}
	UserId = 0
}

func genUserId() (id int32) {
	UserId++
	id = UserId
	return
}


