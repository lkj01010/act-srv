package logic

import (
	. "github.com/lkj01010/act-srv/agent/types"
)

var Cmd = map[string]int16{
	"heartbeat_req":         0, // 心跳包..
	"heartbeat_ack":         1, // 心跳包回复
	"login_req":            10, // 登陆
	"login_ack":            11, // 登陆回执
	"client_error_ack":       13, // 客户端错误
	//"get_seed_req":           30, // socket通信加密使用
	//"get_seed_ack":           31, // socket通信加密使用
	"proto_ping_req":         1001, //  ping
	"proto_ping_ack":         1002, //  ping回复
}

var RCmd = map[int16]string{
	0:    "heartbeat_req", // 心跳包..
	1:    "heartbeat_ack", // 心跳包回复
	10:   "login_req", // 登陆
	11:   "login_ack", // 登陆回执
	13:   "client_error_ack", // 客户端错误
	//30:   "get_seed_req", // socket通信加密使用
	//31:   "get_seed_ack", // socket通信加密使用
	1001: "proto_ping_req", //  ping
	1002: "proto_ping_ack", //  ping回复
}

var Handlers map[int16]func(*Session) []byte
var UserId int32

func init() {
	Handlers = map[int16]func(*Session) []byte{
		0:  P_heartbeat_req,
		10: P_login_req,
		//30: P_get_seed_req,
	}
	UserId = 0
}

func genUserId() (id int32) {
	UserId++
	id = UserId
	return
}


