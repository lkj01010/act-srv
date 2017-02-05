package test

import (
	"testing"
	"github.com/lkj01010/act-srv/agent"
	"time"
)

func TestSmoke(t *testing.T) {
	client, _ := agent.NewClient()
	defer client.Close()

	for i := 0; i < 1; i++ {
		//client.Send([]byte("123"))
		client.Req_login()
		time.Sleep(5*time.Second)
	}
}
