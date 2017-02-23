package gamelogic

type Session struct {
	Flag   int32 // 会话标记
	UserId int32

	Player *Player

	gameCh chan *Game_Frame
}