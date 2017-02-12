package client

import (
	"github.com/lkj01010/act-srv/agent/logic/agentpb"
	"github.com/lkj01010/goutils/log"
	"github.com/golang/protobuf/proto"
)

func H_heartbeat_ack(s *session, payload []byte) ([]byte, error) {
	pb := new(agentpb.HeartbeatAck)
	if err := proto.Unmarshal(payload, pb); err != nil {
		log.Errorf("unmarsha HeartbeatAck err=%+v", err)
		return nil, err
	}
	foo := pb.GetFoo()
	log.Infof("H_heartbeat_ack foo=%+v", foo)
	return nil, nil
}

func H_login_ack(s *session, payload []byte) ([]byte, error) {
	pb := new(agentpb.LoginAck)
	if err := proto.Unmarshal(payload, pb); err != nil {
		log.Errorf("unmarshal LoginAck err=%+v", err)
		return nil, err
	}
	userId := pb.GetUserId()
	log.Infof("H_heartbeat_ack userId=%+v", userId)
	return nil, nil
}
