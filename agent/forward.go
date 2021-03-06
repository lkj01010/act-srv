package agent

import (
	"errors"

	"github.com/lkj01010/log"
)

import (
	. "github.com/lkj01010/act-srv/agent/types"
	"github.com/lkj01010/act-srv/pb"
)

var (
	ERROR_STREAM_NOT_OPEN = errors.New("stream not opened yet")
)

// forward messages to game server
func forward(sess *Session, p []byte) error {
	frame := &pb.Game_Frame{
		Type:    pb.Game_Message,
		Message: p,
	}

	// check stream
	if sess.Stream == nil {
		return ERROR_STREAM_NOT_OPEN
	}

	// forward the frame to game
	//log.Infof("forward send=%+v", frame)
	if err := sess.Stream.Send(frame); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
