package p2p

import (
	"io"
	"net"
	"time"

	log "github.com/cihub/seelog"
	"library/p2p/flowctrl"
)

const (
	// 最多从其它Peer发送请求
	maxOurRequests = 30
)

const (
	// HAVE 每当客户端下载了一个piece，即将该piece的下标作为have消息的负载构造have消息，并把该消息发送给所有建立连接的peer
	HAVE = iota

	// BITFIELD 交换位图
	BITFIELD

	// REQUEST 向该peer发送数据请求
	REQUEST

	// PIECE 当客户端收到某个peer的request消息后,则发送piece消息将文件数据传给该peer。
	PIECE
)

// 下载连接端
type peer struct {
	taskID  string   // 任务标识
	address string   // 对端地址
	conn    net.Conn // 物理连接
	client  bool     // 对端是否为客户端

	writeChan      chan []byte      // 连接的写Chan
	flowctrlWriter *flowctrl.Writer // 基于流控的写

	lastReadTime time.Time //
	have         *Bitset   // 已有的Piece

	ourRequests map[uint64]time.Time // What we requested, when we requested it
}

type peerMessage struct {
	peer    *peer
	message []byte // nil means an error occurred
}

// newPeer ...
func newPeer(c *PeerConn, speed int64) *peer {
	writeChan := make(chan []byte)
	return &peer{
		taskID:         c.taskID,
		conn:           c.conn,
		address:        c.remoteAddr.String(),
		client:         c.client,
		writeChan:      writeChan,
		flowctrlWriter: flowctrl.NewWriter(c.conn, speed),
		ourRequests:    make(map[uint64]time.Time, maxOurRequests),
	}
}

// Close ...
func (p *peer) Close() {
	log.Infof("[%s] Closing connection to %s", p.taskID, p.address)
	p.conn.Close()
	//close(p.writeChan)
}

func (p *peer) sendMessage(b []byte) {
	p.writeChan <- b
}

func (p *peer) keepAlive() {
	p.sendMessage([]byte{})
}

// This func is designed to be run as a goroutine. It
// listens for messages on a channel and sends them to a peer.
func (p *peer) peerWriter(errorChan chan peerMessage) {
	log.Infof("[%s] Writing messages to peer[%s]", p.taskID, p.address)
	var lastWriteTime time.Time

	for msg := range p.writeChan {
		now := time.Now()
		if len(msg) == 0 {
			// This is a keep-alive message.
			if now.Sub(lastWriteTime) < 2*time.Minute {
				continue
			}
			log.Tracef("[%s] Sending keep alive to peer[%s]", p.taskID, p.address)
		}
		lastWriteTime = now

		//log.Debugf("[%s] Sending message to peer[%s], length=%v", p.taskID, p.address, uint32(len(msg)))
		err := writeNBOUint32(p.flowctrlWriter, uint32(len(msg)))
		if err != nil {
			log.Error(err)
			break
		}
		_, err = p.flowctrlWriter.Write(msg)
		if err != nil {
			log.Errorf("[%s] Failed to write a message to peer[%s], length=%v, err=%v", p.taskID, p.address, len(msg), err)
			break
		}
	}

	log.Infof("[%s] Exiting Writing messages to peer[%s]", p.taskID, p.address)
	errorChan <- peerMessage{p, nil}
}

// This func is designed to be run as a goroutine. It
// listens for messages from the peer and forwards them to a channel.
func (p *peer) peerReader(msgChan chan peerMessage) {
	log.Infof("[%s] Reading messages from peer[%s]", p.taskID, p.address)
	for {
		var n uint32
		n, err := readNBOUint32(p.conn)
		if err != nil {
			break
		}
		if n > maxBlockLen {
			log.Error("[", p.taskID, "] Message size too large: ", n)
			break
		}

		var buf []byte
		if n == 0 {
			// keep-alive - we want an empty message
			buf = make([]byte, 1)
		} else {
			buf = make([]byte, n)
		}

		_, err = io.ReadFull(p.conn, buf)
		if err != nil {
			break
		}
		msgChan <- peerMessage{p, buf}
	}

	msgChan <- peerMessage{p, nil}
	log.Infof("[%s] Exiting reading messages from peer[%s]", p.taskID, p.address)
}

// SendBitfield 发送位图
func (p *peer) SendBitfield(bs *Bitset) {
	msg := make([]byte, len(bs.Bytes())+1)
	msg[0] = BITFIELD
	copy(msg[1:], bs.Bytes())
	log.Tracef("[%s] send BITFIELD to peer[%s]", p.taskID, p.address)
	p.sendMessage(msg)
}

// SendHave ...
func (p *peer) SendHave(piece uint32) {
	haveMsg := make([]byte, 5)
	haveMsg[0] = HAVE
	uint32ToBytes(haveMsg[1:5], piece)
	log.Tracef("[%s] send HAVE to peer[%s], piece=%v", p.taskID, p.address, piece)
	p.sendMessage(haveMsg)
}

// SendRequest ...
func (p *peer) SendRequest(piece, begin, length int) {
	req := make([]byte, 13)
	req[0] = byte(REQUEST)
	uint32ToBytes(req[1:5], uint32(piece))
	uint32ToBytes(req[5:9], uint32(begin))
	uint32ToBytes(req[9:13], uint32(length))
	requestIndex := (uint64(piece) << 32) | uint64(begin)

	p.ourRequests[requestIndex] = time.Now()
	log.Tracef("[%s] send REQUEST to peer[%s], piece=%v, begin=%v, length=%v",
		p.taskID, p.address, piece, begin, length)
	p.sendMessage(req)
}
