package p2p

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"path/filepath"
	"time"

	log "github.com/cihub/seelog"
	"github.com/linclin/gopub/src/library/p2p/common"

	"runtime"
)

const (
	// 同一地址最大连接次数
	maxRetryConnectTimes = 2
)

// TaskSession ...
type TaskSession struct {
	// 全局信息
	g *global

	// 任务信息
	taskID    string
	task      *DispatchTask
	fileStore FileStore

	// 下载过程中的Pieces信息
	pieceSet        *Bitset // 本节点已存在Piece
	totalPieces     int     // 整个Piece个数
	totalSize       int64   // 所有文件大小
	lastPieceLength int     // 最一块Piece的长度
	goodPieces      int     // 已下载的Piece个数
	downloaded      uint64  // 已下载的字节数
	checkPieceTime  float64 // 检查Piece所花费的时间累计

	// 正在下载的Piece
	activePieces map[int]*ActivePiece

	// Peer信息
	addPeerChan     chan *PeerConn
	startChan       chan *StartTask
	peers           map[string]*peer
	peerMessageChan chan peerMessage

	// 重新连接定时器
	retryConnTimeChan <-chan time.Time
	indexInChain      int
	connFailCount     int

	//
	quitChan     chan struct{}
	endedChan    chan struct{}
	stopSessChan chan string // sessionmgnt

	//
	reportor   *reportor
	reportStep int

	//
	initedAt   time.Time
	startAt    time.Time
	finishedAt time.Time
}

// NewTaskSession ...
func NewTaskSession(g *global, dt *DispatchTask, stopSessChan chan string) (s *TaskSession, err error) {
	s = &TaskSession{
		g:      g,
		taskID: dt.TaskID,
		task:   dt,

		activePieces: make(map[int]*ActivePiece),
		peers:        make(map[string]*peer),

		addPeerChan:     make(chan *PeerConn, 5), // 不要阻塞
		startChan:       make(chan *StartTask),
		peerMessageChan: make(chan peerMessage, 5),

		quitChan:  make(chan struct{}),
		endedChan: make(chan struct{}),

		stopSessChan: stopSessChan,
		reportor:     newReportor(dt.TaskID, g.cfg),
	}
	return
}

func (s *TaskSession) init() error {
	log.Infof("[%s] Initing p2p session...", s.taskID)
	fileSystem, err := s.g.fsProvider.NewFS()
	if err != nil {
		return err
	}

	// 初始化存储
	m := s.task.MetaInfo
	s.fileStore, s.totalSize, err = NewFileStore(m, fileSystem)
	if err != nil {
		return err
	}

	s.totalPieces, s.lastPieceLength = countPieces(s.totalSize, m.PieceLen)
	return nil
}

func (s *TaskSession) initInServer() error {
	if err := s.init(); err != nil {
		return err
	}

	s.goodPieces = int(s.totalPieces)
	// 标识服务端都是下载完成的
	s.pieceSet = NewBitset(s.goodPieces)
	for index := 0; index < s.goodPieces; index++ {
		s.pieceSet.Set(index)
	}

	log.Infof("[%s] Inited p2p server session", s.taskID)
	s.initedAt = time.Now()
	return nil
}

func (s *TaskSession) initInClient() error {
	// 客户端与服务端的下载路径不同，修改路径
	exsited := false
	for _, fd := range s.task.MetaInfo.Files {
		fd.Path = s.g.cfg.DownDir
		exsited = common.FileExist(filepath.Join(s.g.cfg.DownDir, fd.Name))
	}

	if err := s.init(); err != nil {
		return err
	}

	//计算已经下载的块信息
	if exsited {
		var err error
		start := time.Now()
		s.goodPieces, _, s.pieceSet, err = checkPieces(s.fileStore, s.totalSize, s.task.MetaInfo)
		end := time.Now()
		s.checkPieceTime += end.Sub(start).Seconds()
		log.Infof("[%s] Computed missing pieces: total(%v), good(%v) (%.2f seconds)", s.taskID,
			s.totalPieces, s.goodPieces, s.checkPieceTime)
		if err != nil {
			return err
		}
	} else {
		s.pieceSet = NewBitset(s.totalPieces)
		s.goodPieces = 0
	}

	log.Infof("[%s] Inited p2p client session", s.taskID)
	s.initedAt = time.Now()
	return nil
}

func (s *TaskSession) initPeersBitset() {
	// Enlarge any existing peers piece maps
	for _, p := range s.peers {
		if p.have.n != s.totalPieces {
			if p.have.n != 0 {
				log.Error("Expected p.have.n == 0")
				panic("Expected p.have.n == 0")
			}
			p.have = NewBitset(s.totalPieces)
		}
	}
}

// Start ...
func (s *TaskSession) Start(st *StartTask) {
	s.startChan <- st
}

func (s *TaskSession) startImp(st *StartTask) {
	if s.g.cfg.Server {
		s.startAt = time.Now()
		return
	}

	if s.totalPieces == s.goodPieces {
		// 本地文件的Piece与Block都下载完成，不再需要下载
		log.Infof("[%s] All piece has already download.", s.taskID)
		go s.reportStatus(float32(100))
		return
	}

	log.Infof("[%s] Starting p2p session...", s.taskID)
	// 更新路径
	s.task.LinkChain = st.LinkChain

	// 找到分发路径中位置
	net := s.g.cfg.Net
	self := fmt.Sprintf("%s:%v", net.IP, net.DataPort)
	addrs := s.task.LinkChain.DispatchAddrs
	count := len(addrs)
	for idx := count - 1; idx > 0; idx-- {
		if self == addrs[idx] {
			s.indexInChain = idx - 1
			break
		}
	}
	// 尝试与上一个节点建立连接
	s.tryNewPeer()
	s.initPeersBitset()
	s.startAt = time.Now()
	log.Infof("[%s] Started p2p client session", s.taskID)
}

// 寻找可用的地址并连接
func (s *TaskSession) tryNewPeer() {
	addrs := s.task.LinkChain.DispatchAddrs
	//先尝试服务端
	err := s.connectToPeer(addrs[0])
	if err != nil {
		if s.connFailCount >= maxRetryConnectTimes {
			s.indexInChain--
		}
		if s.indexInChain < 0 {
			s.indexInChain = 0
		}
		peer := addrs[s.indexInChain]
		err1 := s.connectToPeer(peer)
		//执行失败时尝试第一个
		if err1 != nil {
			//s.connectToPeer( addrs[0])
		}
	}

}

// 连接其它的Peer
func (s *TaskSession) connectToPeer(peer string) error {
	defer func() {
		if err := recover(); err != nil {
			var buf []byte = make([]byte, 1024)
			c := runtime.Stack(buf, false)
			log.Errorf("connectToPeer出错:", string(buf[0:c]))
			log.Errorf("Panic error:", err)
		}
	}()
	log.Debugf("[%s] Try connect to peer[%s]", s.taskID, peer)
	conn, err := net.DialTimeout("tcp", peer, 1*time.Second)
	if err != nil {
		log.Errorf("[%s] Failed to connect to peer[%s], error=%v", s.taskID, peer, err)
		conn.Close()
		s.connFailCount++
		s.retryConnTimeChan = time.After(50 * time.Microsecond)
		return err
	}

	// 发送消息头，用于认证
	err = writePHeader(conn, s.taskID, s.g.cfg)
	if err != nil {
		log.Errorf("[%s] Failed to send header to peer[%s], error=%v", s.taskID, peer, err)
		conn.Close()
		s.indexInChain-- //连接下一个
		s.retryConnTimeChan = time.After(50 * time.Microsecond)
		return err
	}

	// 阻塞接收响应
	bs := make([]byte, 1)
	_, err = conn.Read(bs)
	if err != nil {
		// 认证通过了，但没有返回正确的响应，Peer还没创建对应Task的Session
		log.Errorf("[%s] Failed to reading header from peer[%s], error=%v", s.taskID, peer, err)
		conn.Close()
		s.retryConnTimeChan = time.After(50 * time.Microsecond)
		return err
	}

	s.connFailCount = 0
	log.Infof("[%s] Success to connect to peer[%s]", s.taskID, peer)
	p2pconn := &PeerConn{
		conn:       conn,
		client:     false, // 对端是Server
		remoteAddr: conn.RemoteAddr(),
		taskID:     s.taskID,
	}

	s.addPeerImp(p2pconn)
	return nil
}

// AcceptNewPeer 接入其它的Peer连接
func (s *TaskSession) AcceptNewPeer(c *PeerConn) {
	// 先回一个连接响应
	_, err := c.conn.Write([]byte{byte(0xFF)})
	if err != nil {
		log.Errorf("[%s] Write connection init response to peer[%s] failed", s.taskID, c.remoteAddr.String())
		return
	}
	s.addPeerChan <- c
}

// 处理连接到其它成功的Peer，或者是其它Peer的接入
func (s *TaskSession) addPeerImp(c *PeerConn) {
	peerAddr := c.remoteAddr.String()
	log.Infof("[%s] Add new peer, peer[%s]", c.taskID, peerAddr)
	// 创建一个Peer对象
	ps := newPeer(c, s.task.Speed)

	// 位图
	ps.have = NewBitset(s.totalPieces)
	s.peers[peerAddr] = ps

	// 一个从连接上写消息，或读消息
	go ps.peerWriter(s.peerMessageChan)
	go ps.peerReader(s.peerMessageChan)

	// 连接建立之后， 把自己的位置信息给对端
	if s.pieceSet != nil {
		ps.SendBitfield(s.pieceSet)
	}
}

// 关闭Peer
func (s *TaskSession) closePeerAndTryReconn(peer *peer) {
	s.ClosePeer(peer)
	if !peer.client {
		s.tryNewPeer()
	}
}

// ClosePeer 关闭Peer
func (s *TaskSession) ClosePeer(peer *peer) {
	peer.Close()
	s.removeRequests(peer)
	delete(s.peers, peer.address)
}

// 删除REQUEST信息
func (s *TaskSession) removeRequests(p *peer) (err error) {
	for k := range p.ourRequests {
		piece := int(k >> 32)
		begin := int(k & 0xffffffff)
		block := begin / standardBlockLen
		log.Infof("[%s] Forgetting we requested block %v.%v", s.taskID, piece, block)
		s.removeRequest(piece, block)
	}
	p.ourRequests = make(map[uint64]time.Time, maxOurRequests)
	return
}

// 删除REQUEST信息
func (s *TaskSession) removeRequest(piece, block int) {
	v, ok := s.activePieces[piece]
	if ok && v.downloaderCount[block] > 0 {
		v.downloaderCount[block]--
	}
}

// 接收Peer消息并发送消息
func (s *TaskSession) doMessage(p *peer, message []byte) (err error) {
	if message == nil {
		return io.EOF // The reader or writer goroutine has exited
	}

	if len(message) == 0 { // keep alive
		return
	}

	err = s.generalMessage(message, p)
	return
}

func (s *TaskSession) generalMessage(message []byte, p *peer) (err error) {
	messageID := message[0]

	switch messageID {
	case HAVE: // 处理Peer发送过来的HAVE消息
		log.Tracef("[%s] Recv HAVE from peer[%s] ", p.taskID, p.address)
		if len(message) != 5 {
			return errors.New("Unexpected length")
		}
		n := bytesToUint32(message[1:])
		if n >= uint32(p.have.n) {
			return errors.New("have index is out of range")
		}
		p.have.Set(int(n))
		if !p.client {
			for i := 0; i < maxOurRequests; i++ {
				s.requestBlock(p) // 向请此Peer上请求发送块
			}
		}
	case BITFIELD: // 处理Peer发送过来的BITFIELD消息
		log.Tracef("[%s] Recv BITFIELD from peer[%s] isclient=%v", p.taskID, p.address, p.client)
		p.have = NewBitsetFromBytes(s.totalPieces, message[1:])
		if p.have == nil {
			return errors.New("Invalid bitfield data")
		}
		if !p.client {
			s.requestBlock(p) // 向Server Peer请求发送块
		}
	case REQUEST: // 处理Peer发送过来的REQUEST消息
		log.Tracef("[%s] Recv REQUEST from peer[%s] ", p.taskID, p.address)
		index, begin, length, err := s.decodeRequest(message, p)
		if err != nil {
			return err
		}
		return s.sendPiece(p, index, begin, length)
	case PIECE: // 处理Peer发送过来的PIECE消息
		log.Tracef("[%s] Recv PIECE from peer[%s]", p.taskID, p.address)
		index, begin, length, err := s.decodePiece(message, p)
		if err != nil {
			return err
		}

		if s.pieceSet.IsSet(int(index)) {
			log.Debugf("[%s] Recv PIECE from peer[%s] is already", p.taskID, p.address)
			err = s.requestBlock(p)
			break //  本Peer已存在此Piece，则继续
		}

		globalOffset := int64(index)*s.task.MetaInfo.PieceLen + int64(begin)
		_, err = s.fileStore.WriteAt(message[9:], globalOffset)
		if err != nil {
			return err
		}

		// 存储块的信息
		s.recordBlock(p, index, begin, uint32(length))
		err = s.requestBlock(p) // 继续向此Peer请求发送块信息
	default:
		return fmt.Errorf("Uknown message id: %d\n", messageID)
	}

	return
}

func (s *TaskSession) decodeRequest(message []byte, p *peer) (index, begin, length uint32, err error) {
	if len(message) != 13 {
		err = errors.New("Unexpected message length")
		return
	}
	index = bytesToUint32(message[1:5])
	begin = bytesToUint32(message[5:9])
	length = bytesToUint32(message[9:13])
	if index >= uint32(p.have.n) {
		err = errors.New("piece out of range")
		return
	}
	if !s.pieceSet.IsSet(int(index)) {
		err = errors.New("we don't have that piece")
		return
	}
	if int64(begin) >= s.task.MetaInfo.PieceLen {
		err = errors.New("begin out of range")
		return
	}
	if int64(begin)+int64(length) > s.task.MetaInfo.PieceLen {
		err = errors.New("begin + length out of range")
		return
	}
	return
}

// 给Peer发送块消息
func (s *TaskSession) sendPiece(p *peer, index, begin, length uint32) (err error) {
	log.Debugf("[%s] Sending block to peer[%s], index=%v, begin=%v, length=%v",
		s.taskID, p.address, index, begin, length)
	buf := make([]byte, length+9)
	buf[0] = PIECE
	uint32ToBytes(buf[1:5], index)
	uint32ToBytes(buf[5:9], begin)
	_, err = s.fileStore.ReadAt(buf[9:],
		int64(index)*s.task.MetaInfo.PieceLen+int64(begin))
	if err != nil {
		log.Errorf("[%s] Read file failed, error=%v", s.taskID, err)
		return
	}
	p.sendMessage(buf)

	return
}

// 接收块消息
func (s *TaskSession) recordBlock(p *peer, piece, begin, length uint32) (err error) {
	block := begin / standardBlockLen
	log.Debugf("[%s] Received block from peer[%s] %v.%v", s.taskID, p.address, piece, block)

	requestIndex := (uint64(piece) << 32) | uint64(begin)
	delete(p.ourRequests, requestIndex)
	v, ok := s.activePieces[int(piece)]
	if !ok {
		log.Debugf("[%s] Received a block we already have from peer[%s], piece=%v.%v", s.taskID, p.address, piece, block)
		return
	}

	v.recordBlock(int(block))
	s.downloaded += uint64(length)
	if !v.isComplete() {
		return
	}

	// Piece完成下载，清理资源，提交文件
	delete(s.activePieces, int(piece))
	start := time.Now()
	good, pieceBytes, err := checkPiece(s.fileStore, s.totalSize, s.task.MetaInfo, int(piece))
	s.checkPieceTime += time.Now().Sub(start).Seconds()
	if !good || err != nil {
		log.Errorf("[%s] Closing peer[%s] that sent a bad piece=%v, error=%v", s.taskID, p.address, piece, err)
		go s.reportStatus(float32(-1))
		p.Close()
		return
	}

	// 提交文件存储
	s.fileStore.Commit(int(piece), pieceBytes, s.task.MetaInfo.PieceLen*int64(piece))
	s.pieceSet.Set(int(piece))
	s.goodPieces++

	var percentComplete float32
	if s.totalPieces > 0 {
		percentComplete = float32(s.goodPieces*100) / float32(s.totalPieces)
	}
	log.Debugf("[%s] Have %v of %v pieces %v%% complete", s.taskID, s.goodPieces, s.totalPieces,
		percentComplete)
	if s.goodPieces == s.totalPieces {
		s.finishedAt = time.Now() // 下载完成
		go s.reportStatus(percentComplete)
	} else {
		// 减少上报次数，减轻Server的压力
		if int(percentComplete) > s.reportStep {
			s.reportStep += 10
			go s.reportStatus(percentComplete)
		}
	}

	// 每当客户端下载了一个piece，即将该piece的下标作为have消息的负载构造have消息，
	// 并把该消息发送给所有建立连接的Peer。
	for _, p := range s.peers {
		if p.have != nil &&
			(int(piece) >= p.have.n || !p.have.IsSet(int(piece))) {
			p.SendHave(piece)
		}
	}

	return
}

func (s *TaskSession) decodePiece(message []byte, p *peer) (index, begin, length uint32, err error) {
	if len(message) < 9 {
		err = errors.New("unexpected message length")
		return
	}
	index = bytesToUint32(message[1:5])
	begin = bytesToUint32(message[5:9])
	length = uint32(len(message) - 9)

	if index >= uint32(p.have.n) {
		err = errors.New("piece out of range")
		return
	}

	if int64(begin) >= s.task.MetaInfo.PieceLen {
		err = errors.New("begin out of range")
		return
	}
	if int64(begin)+int64(length) > s.task.MetaInfo.PieceLen {
		err = errors.New("begin + length out of range")
		return
	}
	if length > maxBlockLen {
		err = errors.New("Block length too large")
		return
	}
	return
}

// 请求下载时，选择一个可用的Piece
func (s *TaskSession) choosePiece(p *peer) (piece int) {
	n := s.totalPieces
	start := rand.Intn(n)
	piece = s.checkRange(p, start, n)
	if piece == -1 {
		piece = s.checkRange(p, 0, start)
	}
	return
}

func (s *TaskSession) checkRange(p *peer, start, end int) (piece int) {
	clampedEnd := min(end, min(p.have.n, s.pieceSet.n))
	for i := start; i < clampedEnd; i++ {
		// 本Peer没有，但其它Peer存在时
		if (!s.pieceSet.IsSet(i)) && p.have.IsSet(i) {
			if _, ok := s.activePieces[i]; !ok {
				return i
			}
		}
	}
	return -1
}

// 构建请求块（本Peer缺失）信息
func (s *TaskSession) requestBlock(p *peer) (err error) {
	for k := range s.activePieces {
		if p.have.IsSet(k) {
			err = s.requestBlock2(p, k, false)
			if err != io.EOF {
				return
			}
		}
	}

	// No active pieces. (Or no suitable active pieces.) Pick one
	piece := s.choosePiece(p)
	if piece < 0 {
		for k := range s.activePieces {
			if p.have.IsSet(k) {
				err = s.requestBlock2(p, k, true)
				if err != io.EOF {
					return
				}
			}
		}
	}

	// 所有piece与block都下载完成了
	if piece < 0 {
		return
	}

	s.activePieces[piece] = NewActivePiece(s.pieceLength(piece))
	return s.requestBlock2(p, piece, false)

}

func (s *TaskSession) requestBlock2(p *peer, piece int, endGame bool) (err error) {
	v := s.activePieces[piece]
	block := v.chooseBlockToDownload(endGame)
	if block >= 0 {
		s.requestBlockImp(p, piece, block)
	} else {
		//log.Debugf("[%s] Request block from peer[%s], EOF", s.taskID, p.address)
		return io.EOF
	}
	return
}

// Request a block
func (s *TaskSession) requestBlockImp(p *peer, piece int, block int) {
	begin := block * standardBlockLen
	length := standardBlockLen
	if piece == s.totalPieces-1 {
		left := s.lastPieceLength - begin
		if left < length {
			length = left
		}
	}

	//log.Tracef("[%s] Requesting block from peer[%s], piece=%v.%v, length=%v", s.taskID, p.address, piece, block, length)
	p.SendRequest(piece, begin, length)
	return
}

func (s *TaskSession) pieceLength(piece int) int {
	if piece < s.totalPieces-1 {
		return int(s.task.MetaInfo.PieceLen)
	}
	return s.lastPieceLength
}

// Quit ...
func (s *TaskSession) Quit() (err error) {
	select {
	case s.quitChan <- struct{}{}:
	case <-s.endedChan: // 防quit阻塞
	}
	return
}

func (s *TaskSession) shutdown() (err error) {
	for _, peer := range s.peers {
		s.ClosePeer(peer)
	}

	if s.fileStore != nil {
		err = s.fileStore.Close()
		if err != nil {
			log.Errorf("[%s] Error closing filestore : %v", s.taskID, err)
		}
	}

	if s.reportor != nil {
		s.reportor.Close()
	}

	close(s.endedChan)
	return
}

// Init 初始化
func (s *TaskSession) Init() {
	// 开启缓存
	if s.fileStore != nil {
		cache := s.g.cacher.NewCache(s.taskID, s.totalPieces, int(s.task.MetaInfo.PieceLen), s.totalSize)
		s.fileStore.SetCache(cache)
	}

	if s.g.cfg.Server {
		if err := s.initInServer(); err != nil {
			log.Errorf("[%s] Init p2p server session failed, %v", s.taskID, err)
		}
	} else {
		if err := s.initInClient(); err != nil {
			log.Errorf("[%s] Init p2p client session failed, %v", s.taskID, err)
		}
	}

	keepAliveChan := time.Tick(60 * time.Second)
	tickDuration := 2 * time.Second
	tickChan := time.Tick(tickDuration)
	lastDownloaded := s.downloaded

	for {
		select {
		case conn := <-s.addPeerChan:
			s.addPeerImp(conn)
		case st := <-s.startChan:
			s.startImp(st)
		case pm := <-s.peerMessageChan:
			peer, message := pm.peer, pm.message
			peer.lastReadTime = time.Now()
			err2 := s.doMessage(peer, message)
			if err2 != nil {
				if err2 != io.EOF {
					log.Error("[", s.taskID, "] Closing peer[", peer.address, "] because ", err2)
					s.closePeerAndTryReconn(peer)
				} else {
					s.ClosePeer(peer)
				}
			}
		case <-keepAliveChan:
			if s.timeout() {
				// Session超时没有启动，需要stop
				s.stopSessChan <- s.taskID
				log.Info("[", s.taskID, "] P2p session is timeout")
			}
			s.peersKeepAlive()
		case <-tickChan:
			if !s.g.cfg.Server && s.totalPieces != s.goodPieces {
				speed := humanSize(float64(s.downloaded-lastDownloaded) / tickDuration.Seconds())
				lastDownloaded = s.downloaded
				log.Infof("[%s] downloaded: %d(%s/s), pieces: %d/%d, check pieces: (%.2f seconds)",
					s.taskID, s.downloaded, speed, s.goodPieces, s.totalPieces, s.checkPieceTime)
			}
		case <-s.retryConnTimeChan:
			s.tryNewPeer()
		case <-s.quitChan:
			log.Info("[", s.taskID, "] Quit p2p session")
			s.shutdown()
			return
		}
	}
}

func (s *TaskSession) doCheckRequests(p *peer) (err error) {
	now := time.Now()
	for k, v := range p.ourRequests {
		if now.Sub(v).Seconds() > 30 {
			piece := int(k >> 32)
			block := int(k&0xffffffff) / standardBlockLen
			log.Error("[", s.taskID, "] Timing out request of ", piece, ".", block)
			s.removeRequest(piece, block)
		}
	}
	return
}

func (s *TaskSession) peersKeepAlive() {
	now := time.Now()
	for _, peer := range s.peers {
		if peer.lastReadTime.Second() != 0 && now.Sub(peer.lastReadTime) > 3*time.Minute {
			log.Error("[", s.taskID, "] Closing peer [", peer.address, "] because timed out")
			s.ClosePeer(peer)
			continue
		}
		err2 := s.doCheckRequests(peer)
		if err2 != nil {
			if err2 != io.EOF {
				log.Error("[", s.taskID, "] Closing peer[", peer.address, "] because", err2)
			}
			s.ClosePeer(peer)
			continue
		}
		peer.keepAlive()
	}
}

// 检查是否超时了
func (s *TaskSession) timeout() bool {
	now := time.Now()
	if s.startAt.IsZero() && now.Sub(s.initedAt) >= 3*time.Minute {
		return true
	}

	if !s.finishedAt.IsZero() && now.Sub(s.finishedAt) >= 3*time.Minute {
		return true
	}
	return false
}

func (s *TaskSession) reportStatus(pecent float32) {
	s.reportor.DoReport(s.task.LinkChain.ServerAddr, pecent)
}
