package main

import (
	//"testing"
	"github.com/lkj01010/act-srv/client"
	"github.com/lkj01010/goutils/log"
	"sync"
	"os"
	"os/signal"
)

//func Test(t *testing.T) {
//	client, _ := client.NewClient()
//	defer client.Close()
//
//	client.Startup();
//}

func main() {
	var wg sync.WaitGroup

	runClient := func() {
		if client, err := client.NewClient(); err != nil {
			log.Error(err)
		} else {
			defer client.Close()
			client.Startup(&wg)
		}
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go runClient()
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, os.Kill)
		<-sigCh
		log.Info("all clients shutdown")
		os.Exit(0)
		//for {
		//	msg := <-sigCh
		//	switch msg {
		//	case syscall.SIGTERM: // 关闭agent
		//		log.Info("clients shutdown.")
		//		os.Exit(0)
		//	}
		//}
	}()

	wg.Wait()
}
