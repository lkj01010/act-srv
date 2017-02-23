package logic

import (
	. "github.com/lkj01010/act-srv/game/types"
)

//var CmdNames = [...]string{
//	proto.MessageName((*EnterGameReq)(nil)),
//	proto.MessageName((*EnterGameAck)(nil)),
//	proto.MessageName((*BattleInfoNotify)(nil)),
//	proto.MessageName((*SceneNotify)(nil)),
//}

var Cmd = map[string]int16{
	"enter_game_req": 1001,
	"enter_game_ack": 1002,
	"battle_info_notify": 1003,
	"scene_notify": 1004,
}

//var Cmd = map[string]int16{}
var RCmd = map[int16]string{}

var Handlers map[int16]func(*Session, []byte) []byte

func init() {
	for k, v := range Cmd {
		RCmd[v] = k
	}

	Handlers = map[int16]func(*Session, []byte) []byte{
		Cmd["enter_game_req"]: H_enter_game_req,
	}
}

