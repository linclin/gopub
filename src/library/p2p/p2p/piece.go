package p2p

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"runtime"
)

const (
	// 每个Piece分成多个Block，每次下载块的大小
	standardBlockLen = 32 * 1024

	// 最大块的长度
	maxBlockLen = 128 * 1024
)

type chunk struct {
	i    int64
	data []byte
}

func countPieces(totalSize, pieceLen int64) (totalPieces, lastPieceLength int) {
	totalPieces = int(totalSize / pieceLen)
	lastPieceLength = int(totalSize % pieceLen)
	if lastPieceLength == 0 { // last piece is a full piece
		lastPieceLength = int(pieceLen)
	} else {
		totalPieces++
	}
	return
}

// 根据元数据信息，在文件中检查已下载的位图信息，有多少好的Piece，有多少块的Piece
func checkPieces(fs FileStore, totalLength int64, m *MetaInfo) (good, bad int, goodBits *Bitset, err error) {
	pieceLen := m.PieceLen
	totalPieces, _ := countPieces(totalLength, pieceLen)
	goodBits = NewBitset(int(totalPieces))
	ref := m.Pieces
	refLen := len(ref)
	if refLen != totalPieces*sha1.Size {
		err = errors.New(fmt.Sprint("Incorrect MetaInfo.Pieces length ", totalPieces*sha1.Size, "actual length ", refLen))
		return
	}
	currentSums, err := computeSums(fs, totalLength, pieceLen)
	if err != nil {
		return
	}
	for i := 0; i < totalPieces; i++ {
		base := i * sha1.Size
		end := base + sha1.Size
		if checkEqual([]byte(ref[base:end]), currentSums[base:end]) {
			good++
			goodBits.Set(int(i))
		} else {
			bad++
		}
	}
	return
}

// computeSums reads the file content and computes the SHA1 hash for each
// piece. Spawns parallel goroutines to compute the hashes, since each
// computation takes ~30ms.
func computeSums(fs FileStore, totalLength int64, pieceLength int64) (sums []byte, err error) {
	// Calculate the SHA1 hash for each piece in parallel goroutines.
	hashes := make(chan chunk)
	results := make(chan chunk, 3)
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go hashPiece(hashes, results)
	}

	// Read file content and send to "pieces", keeping order.
	numPieces := (totalLength + pieceLength - 1) / pieceLength
	go func() {
		for i := int64(0); i < numPieces; i++ {
			piece := make([]byte, pieceLength, pieceLength)
			if i == numPieces-1 {
				piece = piece[0 : totalLength-i*pieceLength]
			}
			// Ignore errors.
			fs.ReadAt(piece, i*pieceLength)
			hashes <- chunk{i: i, data: piece}
		}
		close(hashes)
	}()

	// Merge back the results.
	sums = make([]byte, sha1.Size*numPieces)
	for i := int64(0); i < numPieces; i++ {
		h := <-results
		copy(sums[h.i*sha1.Size:], h.data)
	}
	return
}

func hashPiece(h chan chunk, result chan chunk) {
	hasher := sha1.New()
	for piece := range h {
		hasher.Reset()
		_, err := hasher.Write(piece.data)
		if err != nil {
			result <- chunk{piece.i, nil}
		} else {
			result <- chunk{piece.i, hasher.Sum(nil)}
		}
	}
}

func computePieceSum(fs FileStore, totalLength int64, pieceLength int64, pieceIndex int) (sum []byte, piece []byte, err error) {
	numPieces := (totalLength + pieceLength - 1) / pieceLength
	hasher := sha1.New()
	piece = make([]byte, pieceLength)
	if int64(pieceIndex) == numPieces-1 {
		piece = piece[0 : totalLength-int64(pieceIndex)*pieceLength]
	}
	_, err = fs.ReadAt(piece, int64(pieceIndex)*pieceLength)
	if err != nil {
		return
	}
	_, err = hasher.Write(piece)
	if err != nil {
		return
	}
	sum = hasher.Sum(nil)
	return
}

func checkPiece(fs FileStore, totalLength int64, m *MetaInfo, pieceIndex int) (good bool, piece []byte, err error) {
	ref := m.Pieces
	var currentSum []byte
	currentSum, piece, err = computePieceSum(fs, totalLength, m.PieceLen, pieceIndex)
	if err != nil {
		return
	}
	base := pieceIndex * sha1.Size
	end := base + sha1.Size
	refSha1 := []byte(ref[base:end])
	good = checkEqual(refSha1, currentSum)
	if !good {
		err = fmt.Errorf("reference sha1: %v != piece sha1: %v", refSha1, currentSum)
	}
	return
}

// ActivePiece 正在下载的Piece
type ActivePiece struct {
	downloaderCount []int // -1 means piece is already downloaded
	pieceLength     int
}

// NewActivePiece ...
func NewActivePiece(pieceLength int) *ActivePiece {
	pieceCount := (pieceLength + standardBlockLen - 1) / standardBlockLen
	return &ActivePiece{make([]int, pieceCount), pieceLength}
}

func (a *ActivePiece) chooseBlockToDownload(endgame bool) (index int) {
	if endgame {
		return a.chooseBlockToDownloadEndgame()
	}
	return a.chooseBlockToDownloadNormal()
}

func (a *ActivePiece) chooseBlockToDownloadNormal() (index int) {
	for i, v := range a.downloaderCount {
		if v == 0 {
			a.downloaderCount[i]++
			return i
		}
	}
	return -1
}

func (a *ActivePiece) chooseBlockToDownloadEndgame() (index int) {
	index, minCount := -1, -1
	for i, v := range a.downloaderCount {
		if v >= 0 && (minCount == -1 || minCount > v) {
			index, minCount = i, v
		}
	}
	if index > -1 {
		a.downloaderCount[index]++
	}
	return
}

func (a *ActivePiece) recordBlock(index int) (requestCount int) {
	requestCount = a.downloaderCount[index]
	a.downloaderCount[index] = -1
	return
}

func (a *ActivePiece) isComplete() bool {
	for _, v := range a.downloaderCount {
		if v != -1 {
			return false
		}
	}
	return true
}
