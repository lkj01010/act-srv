package com

import (
)

const (
	DEFAULT_CH_IPC_SIZE = 16

)

const (
	SessKickedOut = 1 << iota
)

type Msg struct {
	Cmd     int16
	Payload []byte
}

type MsgWithSession struct {
	Ss 	*Session
	Msg
}

type HandleFunc interface {
}

type Session struct {
	ToGameCh chan HandleFunc
	ToAgentCh chan []byte
	Die chan struct{}

	Flag   int32 // 会话标记
	UserId int32
}

//func (g *Session) EnterScene(roomType int32, figure int32, userId int32) {
//	sceneMgr.inCh <- GMH_playerEnter(g, roomType, figure, userId)
//}


//func (ss *Session)Send(frame *pb.Game_Frame) {
//	ss.ToAgentCh <- frame
//}

func NewSession() *Session {
	return &Session{
		ToAgentCh: make(chan []byte, DEFAULT_CH_IPC_SIZE),
		Die: make(chan struct{}),
	}
}

func (ss *Session)GetUserId() int32 {
	return ss.UserId
}

