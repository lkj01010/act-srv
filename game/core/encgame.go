package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/act-srv/game/core/gamepb"
	"github.com/lkj01010/act-srv/misc/packet"
	. "github.com/lkj01010/act-srv/com"
	"github.com/lkj01010/goutils/log"
)

func enc_g_enterGameNtf(userId int32) []byte {
	buf, _ :=  proto.Marshal(&gamepb.EnterGameNtf{
		UserId:  userId,
	})
	log.Debugf("[ntf enterGameNtf][userId=%+v]", userId)
	return packet.Pack(int16(Game_EnterGameNtf), buf, nil)
}

func enc_g_enterGameAck(userId int32, gameId int32) []byte {
	buf, _ :=proto.Marshal(&gamepb.EnterGameAck{
		UserId: userId,
		GameId: gameId,
	})
	log.Debugf("[enterGameAck][userId=%+v][gameId=%+v]", userId, gameId)
	return packet.Pack(int16(Game_EnterGameAck), buf, nil)
}

func enc_g_leaveGameNtf(userId int32) []byte {
	buf, _ :=  proto.Marshal(&gamepb.LeaveGameNtf{
		UserId:  userId,
	})
	log.Debugf("[ntf leaveGameNtf][userId=%+v]", userId)
	return packet.Pack(int16(Game_LeaveGameNtf), buf, nil)
}
