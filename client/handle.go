package client

import (
	"github.com/lkj01010/act-srv/agent/logic/agentpb"
	"github.com/lkj01010/goutils/log"
	"github.com/golang/protobuf/proto"
	"github.com/lkj01010/act-srv/game/core/gamepb"
)

func H_heartbeat_ack(sess *session, payload []byte) ([]byte, error) {
	pb := new(agentpb.HeartbeatAck)
	if err := proto.Unmarshal(payload, pb); err != nil {
		log.Errorf("unmarshal HeartbeatAck err=%+v", err)
		return nil, err
	}
	//seqId := pb.GetSeqId()
	//log.Infof("H_heartbeat_ack seqId=%+v", seqId)
	return nil, nil
}

func H_login_ack(sess *session, payload []byte) ([]byte, error) {
	pb := new(agentpb.LoginAck)
	if err := proto.Unmarshal(payload, pb); err != nil {
		log.Errorf("unmarshal LoginAck err=%+v", err)
		return nil, err
	}
	//userId := pb.GetUserId()
	//log.Infof("H_login_ack userId=%+v", userId)
	return nil, nil
}

func H_enter_game_ntf(sess *session, payload []byte) ([]byte, error) {
	pb := new(gamepb.EnterGameNtf)
	if err := proto.Unmarshal(payload, pb); err != nil {
		log.Errorf("unmarshal EnterGameNtf err=%+v", err)
		return nil, err
	}
	//userId := pb.GetUserId()
	//log.Infof("H_enter_game_ntf userId=%+v", userId)
	return nil, nil
}
