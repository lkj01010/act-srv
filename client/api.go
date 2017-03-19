package client

import (
	. "github.com/lkj01010/act-srv/com"
)

var handlers map[int16]func(*session, []byte) ([]byte, error)

func init() {

	handlers = map[int16]func(*session, []byte) ([]byte, error){
		Cmd[Agent_HeartbeatAck]:  H_heartbeat_ack,
		Cmd[Agent_LoginAck]: H_login_ack,

		Cmd[Game_EnterGameNtf]: H_enter_game_ntf,
	}
}
