package client

import (
	. "github.com/lkj01010/act-srv/com"
)

var handlers map[Cmd]func(*session, []byte) ([]byte, error)

func init() {

	handlers = map[Cmd]func(*session, []byte) ([]byte, error){
		Agent_HeartbeatAck: H_heartbeat_ack,
		Agent_LoginAck:     H_login_ack,

		Game_EnterGameNtf: H_enter_game_ntf,
	}
}
