package agent

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/lkj01010/goutils/log"
)

import (
	"github.com/lkj01010/act-srv/utils"
)

var (
	wg sync.WaitGroup
	// server close signal
	srvDie = make(chan struct{})
)

// handle unix signals
func sig_handler() {
	defer utils.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM: // 关闭agent
			close(srvDie)
			log.Info("sigterm received")
			log.Info("waiting for agents close, please wait...")
			wg.Wait()
			log.Info("agent shutdown.")
			os.Exit(0)
		}
	}
}
