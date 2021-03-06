package agent

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/lkj01010/goutils/log"
)

import (
	"github.com/lkj01010/act-srv/misc/packet"
	. "github.com/lkj01010/act-srv/agent/types"
	. "github.com/lkj01010/act-srv/agent/consts"
	"github.com/lkj01010/act-srv/utils"
)

// PIPELINE #3: buffer
// controls the packet sending for the client
type Buffer struct {
	ctrl    chan struct{} // receive exit signal
	pending chan []byte   // pending packets
	conn    net.Conn      // connection
	cache   []byte        // for combined syscall write
}

var (
	// for padding packet, random content
	// add some random content to confuse packet decrypter
	_padding [PADDING_SIZE]byte
)

func init() {
	go func() { // padding content update procedure
		for {
			for k := range _padding {
				_padding[k] = byte(<-utils.LCG)
			}
			//log.Info("Padding Updated:", _padding)
			<-time.After(PADDING_UPDATE_PERIOD * time.Second)
		}
	}()
}

// packet sending procedure
func (buf *Buffer) send(sess *Session, data []byte) {
	// in case of empty packet
	if data == nil {
		return
	}

	// padding
	// if the size of the data to return is tiny, pad with some random numbers
	// this strategy may change to randomize padding
	if len(data) < PADDING_LIMIT {
		// lkj:modify: del temperately
		//data = append(data, _padding[:]...)
	}

	// encryption
	// (NOT_ENCRYPTED) -> KEYEXCG -> ENCRYPT
	if sess.Flag&SESS_ENCRYPT != 0 { // encryption is enabled
		sess.Encoder.XORKeyStream(data, data)
	} else if sess.Flag&SESS_KEYEXCG != 0 { // key is exchanged, encryption is not yet enabled
		sess.Flag &^= SESS_KEYEXCG
		sess.Flag |= SESS_ENCRYPT
	}

	// queue the data for sending
	buf.pending <- data
	return
}

// packet sending goroutine
func (buf *Buffer) serve() {
	defer utils.PrintPanicStack()
	for {
		select {
		case data := <-buf.pending:
			buf.raw_send(data)
		case <-buf.ctrl: // receive session end signal
			close(buf.pending)
			// close the connection
			buf.conn.Close()
			return
		}
	}
}

// raw packet encapsulation and put it online
func (buf *Buffer) raw_send(data []byte) bool {
	// combine output to reduce syscall.write
	sz := len(data)
	binary.BigEndian.PutUint16(buf.cache, uint16(sz))
	copy(buf.cache[2:], data)

	// write data
	n, err := buf.conn.Write(buf.cache[:sz+2])
	if err != nil {
		log.Warningf("Error send reply data, bytes: %v reason: %v", n, err)
		return false
	}

	return true
}

// create a associated write buffer for a session
func new_buffer(conn net.Conn, ctrl chan struct{}) *Buffer {
	buf := Buffer{conn: conn}
	buf.pending = make(chan []byte)
	buf.ctrl = ctrl
	buf.cache = make([]byte, packet.PACKET_LIMIT+2)
	return &buf
}
