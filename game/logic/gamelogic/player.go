package gamelogic

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

type Sender interface {
	Send_player_enter_game(player *Player)
}

//////////////////////////////////////////////////

type Player struct {
	//sceneId int32
	// 对scene的操作chan
	Sender
	scene *Scene
}

func NewPlayer(sender Sender) *Player {
	return &Player{
		Sender: sender,
	}
}

func (player *Player) EnterScene(roomType int32, figure int32, userId int32) {
	sceneMgr.Perform(SMH_playerEnter(player))
}

func (player *Player) LeaveScene() {
	sceneMgr.Perform(SMH_playerLeave(player))
}
