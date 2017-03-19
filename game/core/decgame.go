package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/goutils/log"

	"github.com/lkj01010/act-srv/game/core/gamepb"
	. "github.com/lkj01010/act-srv/game/com"
	. "github.com/lkj01010/act-srv/com"
)

var gameDecoders map[int16]func(ss *Session, payload []byte) HandleFunc

func init() {
	gameDecoders = make(map[int16]func(ss *Session, payload []byte) HandleFunc)

	gameDecoders[Cmd[Game_EnterGameReq]] = dec_g_enterGame
	gameDecoders[Cmd[Game_LeaveGameReq]] = dec_g_leaveGame
}

func dec_g_enterGame(ss *Session, payload []byte) HandleFunc {
	return func(g *game) {
		msg := new(gamepb.EnterGameReq)
		if err := proto.Unmarshal(payload, msg); err != nil {
			log.Errorf("unmarshal H_enter_game_req err=%+v", err)
			close(ss.Die)
		}
		roomType := msg.GetRoomType()
		figure := msg.GetFigure()

		g.h_enterGame(ss, roomType, figure)
	}
}

func dec_g_leaveGame(ss *Session, payload []byte) HandleFunc {
	return func(g *game) {
		g.h_leaveGame(ss)
	}
}

