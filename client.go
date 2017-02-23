package main

import (
	//"testing"
	"github.com/lkj01010/act-srv/client"
	"github.com/lkj01010/goutils/log"
)

//func Test(t *testing.T) {
//	client, _ := client.NewClient()
//	defer client.Close()
//
//	client.Startup();
//}

func main() {
	if client, err := client.NewClient(); err != nil {
		log.Error(err)
	} else {
		defer client.Close()
		client.Startup()
	}
}
