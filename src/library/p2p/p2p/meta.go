package p2p

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/cihub/seelog"
)

// fileSystemAdapter FileSystem接口适配
type fileSystemAdapter struct {
}

// Open ...
func (f *fileSystemAdapter) Open(name []string, length int64) (file File, err error) {
	var ff *os.File
	ff, err = os.Open(path.Clean(path.Join(name...)))
	if err != nil {
		return
	}
	stat, err := ff.Stat()
	if err != nil {
		return
	}
	actualSize := stat.Size()
	if actualSize != length {
		err = fmt.Errorf("Unexpected file size %v. Expected %v", actualSize, length)
		return
	}
	file = ff
	return
}

// Close ...
func (f *fileSystemAdapter) Close() error {
	return nil
}

func (m *MetaInfo) addFiles(fileInfo os.FileInfo, file string, idx int) (err error) {
	fileDict := FileDict{Length: fileInfo.Size()}
	cleanFile := path.Clean(file)
	fileDict.Path, fileDict.Name = path.Split(cleanFile)
	fileDict.Sum, err = sha1Sum(file)
	if err != nil {
		return err
	}
	m.Files[idx] = &fileDict
	return
}

// CreateFileMeta ...
func CreateFileMeta(roots []string, pieceLen int64) (mi *MetaInfo, err error) {
	mi = &MetaInfo{Files: make([]*FileDict, len(roots))}
	for idx, f := range roots {
		var fileInfo os.FileInfo
		fileInfo, err = os.Stat(f)
		if err != nil {
			log.Errorf("File not exist file=%s, error=%v", f, err)
			return
		}

		if fileInfo.IsDir() {
			return nil, fmt.Errorf("Not support dir")
		}

		err = mi.addFiles(fileInfo, f, idx)
		if err != nil {
			return nil, err
		}
		mi.Length += fileInfo.Size()
	}

	if pieceLen == 0 {
		pieceLen = choosePieceLength(mi.Length)
	}
	mi.PieceLen = pieceLen

	fileStore, fileStoreLength, err := NewFileStore(mi, &fileSystemAdapter{})
	if err != nil {
		return nil, err
	}
	defer fileStore.Close()
	if fileStoreLength != mi.Length {
		return nil, fmt.Errorf("Filestore total length %v, expected %v", fileStoreLength, mi.Length)
	}

	var sums []byte
	sums, err = computeSums(fileStore, mi.Length, mi.PieceLen)
	if err != nil {
		return nil, err
	}
	mi.Pieces = sums
	log.Debugf("File totallength=%v, piecelength=%v", mi.Length, pieceLen)
	return mi, nil
}

func sha1Sum(file string) (sum string, err error) {
	var f *os.File
	f, err = os.Open(file)
	if err != nil {
		log.Errorf("Open file failed, file=%s, error=%v", file, err)
		return
	}
	defer f.Close()
	hash := sha1.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		log.Errorf("Summary file by sha1 failed, file=%s, error=%v", file, err)
		return
	}
	sum = string(hash.Sum(nil))
	return
}

const (
	minimumPieceLength   = 16 * 1024
	targetPieceCountLog2 = 10
	targetPieceCountMin  = 1 << targetPieceCountLog2

	// Target piece count should be < targetPieceCountMax
	targetPieceCountMax = targetPieceCountMin << 1
)

// Choose a good piecelength.
func choosePieceLength(totalLength int64) (pieceLength int64) {
	// Must be a power of 2.
	// Must be a multiple of 16KB
	// Prefer to provide around 1024..2048 pieces.
	pieceLength = minimumPieceLength
	pieces := totalLength / pieceLength
	for pieces >= targetPieceCountMax {
		pieceLength <<= 1
		pieces >>= 1
	}
	return
}
