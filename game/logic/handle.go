package logic

import (
	. "github.com/lkj01010/act-srv/game/types"
	"github.com/lkj01010/act-srv/game/logic/gamepb"
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/goutils/log"
	"github.com/lkj01010/act-srv/misc/packet"
	"github.com/lkj01010/act-srv/game/logic/gamelogic"
)

func H_enter_game_req(sess *Session, payload []byte) []byte {
	pb := new(gamepb.EnterGameReq)
	if err := proto.Unmarshal(payload, pb); err != nil {
		log.Errorf("unmarshal H_enter_game_req err=%+v", err)
		sess.Flag |= SESS_KICKED_OUT
		return nil
	}
	rootType := pb.GetRoomType()
	figure := pb.GetFigure()
	log.Infof("EnterGameReq rootType=%+v, figure=%+v", rootType, figure)

	sess.Player = gamelogic.NewPlayer(NewSender(sess.UserId))

	sess.Player.EnterScene(rootType, figure, sess.UserId)

	sceneMgr.EnterScene

	ack := &gamepb.EnterGameAckTest{
		UserId: sess.UserId,
	}
	ret, err := proto.Marshal(ack)
	if err != nil {
		log.Error(err)
	}
	return packet.Pack(Cmd["enter_game_ack"], ret, nil)
}

func H_leave_game(sess *Session, payload []byte) []byte {

	return nil
}