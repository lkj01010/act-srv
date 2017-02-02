package agent

import (
	"encoding/binary"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/lkj01010/goutils/log"

	. "github.com/lkj01010/act-srv/agent/types"
	"github.com/lkj01010/act-srv/utils"
	"github.com/lkj01010/act-srv/consts"
)

const (
	SERVICE = "[AGENT]"
)

func Startup() {
	// to catch all uncaught panic
	defer utils.PrintPanicStack()

	// open profiling
	go func() {
		log.Info(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	// startup
	//startup()
	// =>
	go sig_handler()
	// ]

	go tcpServer()

	// wait forever
	select {}
}

func tcpServer() {
	// resolve address & start listening
	tcpAddr, err := net.ResolveTCPAddr("tcp4", consts.AgentPort)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Info("listening on:", listener.Addr())

	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warning("accept failed:", err)
			continue
		}
		// set socket read buffer
		conn.SetReadBuffer(SO_RCVBUF)
		// set socket write buffer
		conn.SetWriteBuffer(SO_SNDBUF)
		// start a goroutine for every incoming connection for reading
		go handleClient(conn)

		// check server close signal
		select {
		case <-die:
			listener.Close()
			return
		default:
		}
	}
}

//func udpServer() {
//	l, err := kcp.Listen(_port)
//	checkError(err)
//	log.Info("udp listening on:", l.Addr())
//	lis := l.(*kcp.Listener)
//
//	// loop accepting
//	for {
//		conn, err := lis.AcceptKCP()
//		if err != nil {
//			log.Warning("accept failed:", err)
//			continue
//		}
//		// set kcp parameters
//		conn.SetWindowSize(32, 32)
//		conn.SetNoDelay(1, 20, 1, 1)
//		conn.SetKeepAlive(0) // require application ping
//		conn.SetStreamMode(true)
//
//		// start a goroutine for every incoming connection for reading
//		go handleClient(conn)
//	}
//}

// PIPELINE #1: handleClient
// the goroutine is used for reading incoming PACKETS
// each packet is defined as :
// | 2B size |     DATA       |
//
func handleClient(conn net.Conn) {
	defer utils.PrintPanicStack()
	// for reading the 2-Byte header
	header := make([]byte, 2)
	// the input channel for agent()
	in := make(chan []byte)
	defer func() {
		close(in) // session will close
	}()

	// create a new session object for the connection
	// and record it's IP address
	var sess Session
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("cannot get remote address:", err)
		return
	}
	sess.IP = net.ParseIP(host)
	log.Infof("new connection from:%v port:%v", host, port)

	// session die signal, will be triggered by agent()
	sess.Die = make(chan struct{})

	//// userId
	//userIdAcc++
	//sess.UserId = userIdAcc
	//
	//// game rpc
	//cli := pb.NewGameServiceClient(gameConn)
	//ctx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userid": fmt.Sprint(sess.UserId)}))
	//stream, err := cli.Stream(ctx)
	//if err != nil {
	//	log.Error(err)
	//	return nil
	//}
	//sess.Stream = stream

	// create a write buffer
	out := new_buffer(conn, sess.Die)
	go out.start()

	// start agent for PACKET processing
	wg.Add(1)
	go agent(&sess, in, out)

	// read loop
	for {
		// solve dead link problem:
		// physical disconnection without any communcation between client and server
		// will cause the read to block FOREVER, so a timeout is a rescue.
		conn.SetReadDeadline(time.Now().Add(TCP_READ_DEADLINE * time.Second))

		// read 2B header
		n, err := io.ReadFull(conn, header)
		if err != nil {
			log.Warningf("read header failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)

		// alloc a byte slice of the size defined in the header for reading data
		payload := make([]byte, size)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			log.Warningf("read payload failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}

		// deliver the data to the input queue of agent()
		select {
		case in <- payload: // payload queued
		case <-sess.Die:
			log.Warningf("connection closed by logic, flag:%v ip:%v", sess.Flag, sess.IP)
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
