package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/goutils/log"

	"github.com/lkj01010/act-srv/game/core/gamepb"
	. "github.com/lkj01010/act-srv/com"
	"sync/atomic"
)

var gameDecoders map[Cmd]func(ss *Session, payload []byte) HandleFunc

func init() {
	gameDecoders = make(map[Cmd]func(ss *Session, payload []byte) HandleFunc)

	gameDecoders[Game_EnterGameReq] = dec_g_enterGame
	gameDecoders[Game_LeaveGameReq] = dec_g_leaveGame
}

func dec_g_enterGame(ss *Session, payload []byte) HandleFunc {
	return func(g *game) {
		msg := new(gamepb.EnterGameReq)
		if err := proto.Unmarshal(payload, msg); err != nil {
			log.Errorf("unmarshal H_enter_game_req err=%+v", err)
			atomic.StoreInt32(&ss.isDie, 1)
			g.h_leaveGame(ss)
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

