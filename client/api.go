package client

import (
	agentLogic "github.com/lkj01010/act-srv/agent/logic"
	gameLogic "github.com/lkj01010/act-srv/game/logic"
)

var handlers map[int16]func(*session, []byte) ([]byte, error)

func init() {

	handlers = map[int16]func(*session, []byte) ([]byte, error){
		agentLogic.Cmd["heartbeat_ack"]:  H_heartbeat_ack,
		agentLogic.Cmd["login_ack"]: H_login_ack,

		gameLogic.Cmd["enter_game_ack"]: H_enter_game_ack,
	}
}
