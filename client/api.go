package client

import (
)

var Cmd = map[string]int16{
	"heartbeat_req":   0,  // 心跳包..
	"heartbeat_ack":   1,  // 心跳包回复
	"login_req":        10, // 登陆
	"login_ack":        11, // 登陆回执
	"client_error_ack": 13, // 客户端错误
	//"get_seed_req":           30, // socket通信加密使用
	//"get_seed_ack":           31, // socket通信加密使用
	"proto_ping_req": 1001, //  ping
	"proto_ping_ack": 1002, //  ping回复
}

//var Cmd = map[string]int16{}
var RCmd = map[int16]string{}

var handlers map[int16]func(*session, []byte) ([]byte, error)

func init() {
	for k, v := range Cmd {
		RCmd[v] = k
	}

	handlers = map[int16]func(*session, []byte) ([]byte, error){
		Cmd["heartbeat_ack"]:  H_heartbeat_ack,
		Cmd["login_ack"]: H_login_ack,
	}
}
