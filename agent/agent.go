package agent

import (
	"time"
)

import (
	pb "github.com/lkj01010/act-srv/pb"
	. "github.com/lkj01010/act-srv/agent/types"
	"github.com/lkj01010/act-srv/utils"
	"github.com/lkj01010/goutils/log"
)

// PIPELINE #2: agent
// all the packets from handleClient() will be handled
func agent(ss *Session, recvCh chan []byte, sender *Buffer) {
	defer wg.Done() // will decrease waitgroup by one, useful for manual server shutdown
	defer utils.PrintPanicStack()

	// init session
	ss.ConnectTime = time.Now()
	ss.LastPacketTime = time.Now()
	// minute timer
	min_timer := time.After(time.Minute)

	// cleanup work
	defer func() {
		log.Info("agent srvDie")
		close(ss.Die)
		if ss.Stream != nil {
			ss.Stream.CloseSend()
		}
	}()

	// >> the main message loop <<
	// handles 4 types of message:
	//  1. from client
	//  2. from game service
	//  3. timer
	//  4. server shutdown signal
	for {
		select {
		case msg, ok := <-recvCh: // packet from network
			if !ok {
				return
			}
			//log.Debugf("recieve msg: %v", msg)

			ss.PacketCount++
			ss.PacketTime = time.Now()

			if result := proxy_user_request(ss, msg); result != nil {
				sender.send(ss, result)
			}
			ss.LastPacketTime = ss.PacketTime
		case frame := <-ss.StreamCh: // packets from game
			switch frame.Type {
			case pb.Game_Message:
				sender.send(ss, frame.Message)
			case pb.Game_Kick:
				ss.Flag |= SESS_KICKED_OUT
			}
		case <-min_timer: // minutes timer
			timer_work(ss, sender)
			min_timer = time.After(time.Minute)
		case <-srvDie:
			ss.Flag |= SESS_KICKED_OUT
		}

		// see if the player should be kicked sender.
		if ss.Flag&SESS_KICKED_OUT != 0 {
			return
		}
	}
}
