package core

import (
	. "github.com/lkj01010/act-srv/game/com"
)
//type UserIpc interface {
//	Send(userId int32, msg []byte)
//}
//
//// 需要在logic外层实例化
//var Ipc UserIpc

// ────────────────────────────────────────────────────────────────────────────────
// 		----- ioLoop ----
//		|				•
//		•				|
// session(player) -> scene ->
//			|			•
//			|			|
// 			•------> sceneMgr
// ────────────────────────────────────────────────────────────────────────────────


//////////////////////////////////////////////////

type Player struct {
	//sceneId int32
	// 对scene的操作chan
	ss *Session
}

func NewPlayer(ss *Session) *Player {
	return &Player{
		ss: ss,
	}
}

//func (p *Player) EnterScene(roomType int32, figure int32, userId int32) {
//	GameMgr.msgSsCh <- GMH_playerEnter(p.ss, roomType, figure, userId)
//}
//
//func (p *Player) LeaveScene() {
//	GameMgr.msgSsCh <- GMH_playerLeave(p)
//}
