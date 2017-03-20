package game

import (
	"net"
	"os"

	"github.com/lkj01010/goutils/log"
	"google.golang.org/grpc"
	. "github.com/lkj01010/act-srv/consts"
	. "github.com/lkj01010/act-srv/pb"
	"github.com/lkj01010/act-srv/services"
	"github.com/lkj01010/act-srv/game/core"
)

func Startup() {
	// 监听
	lis, err := net.Listen("tcp", GamePort)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	log.Info("gamesrv listening on ", lis.Addr())

	// 注册服务
	s := grpc.NewServer()
	ins := core.NewSession()
	RegisterGameServiceServer(s, ins)

	//// 初始化Services
	//sp.Init("snowflake")
	// 开始服务
	s.Serve(lis)

	services.CloseGameConn()
}
