package services

import (
	"sync"
	"google.golang.org/grpc"
	"github.com/lkj01010/act-srv/consts"
	"github.com/lkj01010/goutils/log"
	"time"
)

var (
	once sync.Once
	gameConn *grpc.ClientConn
)

func GetGameConn() *grpc.ClientConn {
	once.Do(func() {
		conn, err := grpc.Dial(consts.GamePort, grpc.WithTimeout(time.Millisecond), grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("connect gamesrv err: %+v", err)
		}
		gameConn = conn
		log.Info("grpc game service connected")
	})

	return gameConn
}

func CloseGameConn() {
	if gameConn != nil {
		gameConn.Close()
	}
}