package main

import (
	//"testing"
	"github.com/lkj01010/act-srv/client"
)

//func Test(t *testing.T) {
//	client, _ := client.NewClient()
//	defer client.Close()
//
//	client.Startup();
//}

func main() {
	client, _ := client.NewClient()
	defer client.Close()

	client.Startup()
}
