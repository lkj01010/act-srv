package com

//const (
//	Agent_Start        = "Agent_Start"
//	Agent_HeartbeatReq = "Agent_HeartbeatReq"
//	Agent_HeartbeatAck = "Agent_HeartbeatAck"
//	Agent_LoginReq     = "Agent_LoginReq"
//	Agent_LoginAck     = "Agent_LoginAck"
//	Agent_End = "Agent_End"
//
//	Game_Start        = "Game_Start"
//	Game_EnterGameReq = "Game_EnterGameReq"
//	Game_EnterGameAck = "Game_EnterGameAck"
//	Game_EnterGameNtf = "Game_EnterGameNtf"
//
//	Game_LeaveGameReq = "Game_LeaveGameReq"
//	Game_LeaveGameNtf = "Game_LeaveGameNtf"
//
//	Game_GameMgrEnd = "Game_GameMgrEnd"
//
//	Game_End = "Game_End"
//)
//
//var Cmd = map[string]int16{
//	Agent_Start:        0,
//	Agent_HeartbeatReq: 1, // 心跳包..
//	Agent_HeartbeatAck: 2, // 心跳包回复
//	Agent_LoginReq:     3, // 登陆
//	Agent_LoginAck:     4, // 登陆回执
//	Agent_End:          999,
//
//	Game_Start:         1000,
//	Game_EnterGameReq:  1001,
//	Game_EnterGameAck:  1002,
//	Game_EnterGameNtf:  1003,
//	Game_LeaveGameReq:  1004,
//	Game_LeaveGameNtf:  1005,
//	Game_GameMgrEnd: 1100,
//	Game_End:           1999,
//
//
//
//}
//
//var RCmd = map[int16]string{}
//
//func init() {
//	for k, v := range Cmd {
//		RCmd[v] = k
//	}
//}
//
const (
	DestAgent   = iota
	DestGameMgr
	DestGame
	DestInvalid
)
//
func CmdSendTo(cmd Cmd) int {
	if cmd < Game_Start {
		return DestAgent
	}
	if cmd >= Game_Start && cmd < Game_GameMgrEnd {
		return DestGameMgr
	} else if cmd >= Game_GameMgrEnd && cmd < Game_End {
		return DestGame
	} else {
		return DestInvalid
	}
}
