package logic

import (
	gl "github.com/lkj01010/act-srv/game/logic/gamelogic"
)

type Sender struct {
	userId int32
}

func NewSender(userId int32) Sender {
	return Sender{userId: userId}
}

func (s Sender) Send_player_enter_game(player *gl.Player) {

}
