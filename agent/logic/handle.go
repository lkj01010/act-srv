package logic

import (
	"fmt"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/lkj01010/goutils/log"
)

import (
	"github.com/lkj01010/act-srv/misc/packet"
	pb "github.com/lkj01010/act-srv/pb"
	. "github.com/lkj01010/act-srv/agent/types"

	//sp "github.com/gonet2/libs/services"
	"github.com/lkj01010/act-srv/services"
	"github.com/lkj01010/act-srv/agent/logic/agentpb"
	"github.com/golang/protobuf/proto"
)

// 心跳包
func H_heartbeat_req(sess *Session) ([]byte, error) {
	//tbl, _ := PKT_auto_id(reader)
	ret, err := proto.Marshal(&agentpb.HeartbeatAck{
		SeqId: sess.PacketCount,
	})
	if err != nil {
		log.Errorf("marshal HeartbeatAck err=%+v", err)
		return nil, nil
	}
	return packet.Pack(Cmd["heartbeat_ack"], ret, nil), nil
	//return []byte("123")
}

// 密钥交换
// 加密建立方式: DH+RC4
// 注意:完整的加密过程包括 RSA+DH+RC4
// 1. RSA用于鉴定服务器的真伪(这步省略)
// 2. DH用于在不安全的信道上协商安全的KEY
// 3. RC4用于流加密
//func P_get_seed_req(sess *Session, reader *packet.Packet) []byte {
//	tbl, _ := PKT_seed_info(reader)
//	// KEY1
//	X1, E1 := dh.DHExchange()
//	KEY1 := dh.DHKey(X1, big.NewInt(int64(tbl.F_client_send_seed)))
//
//	// KEY2
//	X2, E2 := dh.DHExchange()
//	KEY2 := dh.DHKey(X2, big.NewInt(int64(tbl.F_client_receive_seed)))
//
//	ret := S_seed_info{int32(E1.Int64()), int32(E2.Int64())}
//	// 服务器加密种子是客户端解密种子
//	encoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", SALT, KEY2)))
//	if err != nil {
//		log.Error(err)
//		return nil
//	}
//	decoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", SALT, KEY1)))
//	if err != nil {
//		log.Error(err)
//		return nil
//	}
//	sess.Encoder = encoder
//	sess.Decoder = decoder
//	sess.Flag |= SESS_KEYEXCG
//	return packet.Pack(Code["get_seed_ack"], ret, nil)
//}

// 玩家登陆过程
func H_login_req(sess *Session) ([]byte, error) {
	// TODO: 登陆鉴权
	// 简单鉴权可以在agent直接完成，通常公司都存在一个用户中心服务器用于鉴权
	sess.UserId = genUserId()

	// TODO: 选择GAME服务器
	// 选服策略依据业务进行，比如小服可以固定选取某台，大服可以采用HASH或一致性HASH
	sess.GSID = DEFAULT_GSID

	//// 连接到已选定GAME服务器
	//conn := sp.GetServiceWithId(sp.DEFAULT_SERVICE_PATH+"/game", sess.GSID)
	//if conn == nil {
	//	log.Error("cannot get game service:", sess.GSID)
	//	return nil
	//}
	cli := pb.NewGameServiceClient(services.GetGameConn())

	// 开启到游戏服的流
	ctx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userid": fmt.Sprint(sess.UserId)}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sess.Stream = stream

	// 读取GAME返回消息的goroutine
	fetcher_task := func(sess *Session) {
		for {
			in, err := sess.Stream.Recv()
			if err == io.EOF { // 流关闭
				log.Debug(err)
				return
			}
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case sess.MQ <- *in:
			case <-sess.Die:
			}
		}
	}
	go fetcher_task(sess)

	if ret, err := proto.Marshal(&agentpb.LoginAck{UserId: sess.UserId}); err != nil {
		log.Error(err)
		return nil, nil
	} else {
		return packet.Pack(Cmd["login_ack"], ret, nil), nil
	}

	//return packet.Pack(Cmd["login_ack"], S_user_snapshot{F_uid: sess.UserId}, nil)
}
