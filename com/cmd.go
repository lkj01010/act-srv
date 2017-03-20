package com

type Cmd int16

const (
	Agent_Start        Cmd = iota
	Agent_HeartbeatReq  // 心跳包..
	Agent_HeartbeatAck  // 心跳包回复
	Agent_LoginReq      // 登陆
	Agent_LoginAck      // 登陆回执
	Agent_End

	Game_Start        Cmd = iota + 1000
	Game_EnterGameReq
	Game_EnterGameAck
	Game_EnterGameNtf
	Game_LeaveGameReq
	Game_LeaveGameNtf
	Game_GameMgrEnd
	Game_End
)
