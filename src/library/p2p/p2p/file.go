package p2p

import (
	"errors"
	"io"

	log "github.com/cihub/seelog"
)

// File Interface for a file.
// Multiple goroutines may access a File at the same time.
type File interface {
	io.ReaderAt
	io.WriterAt
	io.Closer
}

// FsProvider Interface for a provider of filesystems.
type FsProvider interface {
	NewFS() (FileSystem, error)
}

// FileSystem Interface for a file system. A file system contains files.
type FileSystem interface {
	Open(name []string, length int64) (file File, err error)
	io.Closer
}

// FileStore a file store.
type FileStore interface {
	io.ReaderAt
	io.WriterAt
	io.Closer
	SetCache(FileCache)
	Commit(int, []byte, int64)
}

type fileStore struct {
	fileSystem FileSystem
	offsets    []int64
	files      []fileEntry // Stored in increasing globalOffset order
	cache      FileCache
}

type fileEntry struct {
	length int64
	file   File
}

// NewFileStore 根据元数据信息打开所有文件
func NewFileStore(info *MetaInfo, fileSystem FileSystem) (f FileStore, totalSize int64, err error) {
	fs := &fileStore{}
	fs.fileSystem = fileSystem

	numFiles := len(info.Files)
	fs.files = make([]fileEntry, numFiles)
	fs.offsets = make([]int64, numFiles)

	for i, src := range info.Files {
		var file File
		file, err = fs.fileSystem.Open([]string{src.Path, src.Name}, src.Length)
		if err != nil {
			log.Errorf("Open file failed, file=%v/%v, error=%v", src.Path, src.Name, err)
			// Close all files opened up to now.
			for i2 := 0; i2 < i; i2++ {
				fs.files[i2].file.Close()
			}
			return
		}
		fs.files[i].file = file
		fs.files[i].length = src.Length
		fs.offsets[i] = totalSize
		totalSize += src.Length
	}
	f = fs
	return
}

// SetCache ...
func (f *fileStore) SetCache(cache FileCache) {
	f.cache = cache
}

func (f *fileStore) find(offset int64) int {
	// Binary search
	offsets := f.offsets
	low := 0
	high := len(offsets)
	for low < high-1 {
		probe := (low + high) / 2
		entry := offsets[probe]
		if offset < entry {
			high = probe
		} else {
			low = probe
		}
	}
	return low
}

// ReadAt ...
func (f *fileStore) ReadAt(p []byte, off int64) (int, error) {
	if f.cache == nil {
		return f.RawReadAt(p, off)
	}

	unfullfilled := f.cache.readAt(p, off)

	var retErr error
	for _, unf := range unfullfilled {
		_, err := f.RawReadAt(unf.data, unf.i)
		if err != nil {
			log.Error("Got an error on read (off=", unf.i, "len=", len(unf.data), ") from filestore:", err)
			retErr = err
		}
	}
	return len(p), retErr
}

// RawReadAt ...
func (f *fileStore) RawReadAt(p []byte, off int64) (n int, err error) {
	index := f.find(off)
	for len(p) > 0 && index < len(f.offsets) {
		chunk := int64(len(p))
		entry := &f.files[index]
		itemOffset := off - f.offsets[index]
		if itemOffset < entry.length {
			space := entry.length - itemOffset
			if space < chunk {
				chunk = space
			}
			var nThisTime int
			nThisTime, err = entry.file.ReadAt(p[0:chunk], itemOffset)
			n = n + nThisTime
			if err != nil {
				return
			}
			p = p[nThisTime:]
			off += int64(nThisTime)
		}
		index++
	}
	// At this point if there's anything left to read it means we've run off the
	// end of the file store. Read zeros. This is defined by the bittorrent protocol.
	for i := range p {
		p[i] = 0
	}
	return
}

// WriteAt ...
func (f *fileStore) WriteAt(p []byte, off int64) (int, error) {
	if f.cache != nil {
		needRawWrite := f.cache.writeAt(p, off)
		if needRawWrite != nil {
			for _, nc := range needRawWrite {
				f.RawWriteAt(nc.data, nc.i)
			}
		}
		return len(p), nil
	}
	return f.RawWriteAt(p, off)
}

// Commit ...
func (f *fileStore) Commit(pieceNum int, piece []byte, off int64) {
	if f.cache != nil {
		_, err := f.RawWriteAt(piece, off)
		if err != nil {
			log.Error("Error committing to storage:", err)
			return
		}
		f.cache.MarkCommitted(pieceNum)
	}
}

// RawWriteAt ...
func (f *fileStore) RawWriteAt(p []byte, off int64) (n int, err error) {
	index := f.find(off)
	for len(p) > 0 && index < len(f.offsets) {
		chunk := int64(len(p))
		entry := &f.files[index]
		itemOffset := off - f.offsets[index]
		if itemOffset < entry.length {
			space := entry.length - itemOffset
			if space < chunk {
				chunk = space
			}
			var nThisTime int
			nThisTime, err = entry.file.WriteAt(p[0:chunk], itemOffset)
			n += nThisTime
			if err != nil {
				return
			}
			p = p[nThisTime:]
			off += int64(nThisTime)
		}
		index++
	}
	// At this point if there's anything left to write it means we've run off the
	// end of the file store. Check that the data is zeros.
	// This is defined by the bittorrent protocol.
	for i := range p {
		if p[i] != 0 {
			err = errors.New("Unexpected non-zero data at end of store.")
			n = n + i
			return
		}
	}
	n = n + len(p)
	return
}

// Close ...
func (f *fileStore) Close() (err error) {
	for i := range f.files {
		f.files[i].file.Close()
	}
	if f.cache != nil {
		f.cache.Close()
		f.cache = nil
	}
	if f.fileSystem != nil {
		err = f.fileSystem.Close()
	}
	return
}
