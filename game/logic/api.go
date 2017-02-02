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
	"enter_game_req": 0,
	"enter_game_ack": 1,
	"battle_info_notify": 2,
	"scene_notify": 3,
}

//var Cmd = map[string]int16{}
var RCmd = map[int16]string{}

var Handlers map[int16]func(*Session) []byte

func init() {
	for k, v := range Cmd {
		RCmd[v] = k
	}

	Handlers = map[int16]func(*Session) []byte{
		0: P_enter_game_req,
	}
}

