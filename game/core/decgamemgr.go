package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/goutils/log"

	"github.com/lkj01010/act-srv/game/core/gamepb"
	. "github.com/lkj01010/act-srv/com"
)

var gameMgrDecoders map[Cmd](func(ss *Session, payload []byte) HandleFunc)

func init() {
	gameMgrDecoders = make(map[Cmd]func(ss *Session, payload []byte) HandleFunc)

	gameMgrDecoders[Game_EnterGameReq] = dec_gm_enterGame
}

func dec_gm_enterGame(ss *Session, payload []byte) HandleFunc {
	//log.Info("dec_gm_enterGame")
	return func(gm *gameManager) {
		//log.Info("dec_gm_enterGame =>")
		msg := new(gamepb.EnterGameReq)
		if err := proto.Unmarshal(payload, msg); err != nil {
			log.Errorf("unmarshal dec_gm_enterGame err=%+v", err)
			close(ss.Die)
			// note: 如果这时 isInGame==1, 那么给game发送leaveFunc
		}
		roomType := msg.GetRoomType()
		figure := msg.GetFigure()

		gm.h_enterGame(ss, roomType, figure)
	}
}
