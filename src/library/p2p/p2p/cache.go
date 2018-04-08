package p2p

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/cihub/seelog"
)

// CacheProvider ...
type CacheProvider interface {
	NewCache(infohash string, numPieces int, pieceLength int, totalSize int64) FileCache
}

// FileCache ...
type FileCache interface {
	//Read what's cached, returns parts that weren't available to read.
	readAt(p []byte, offset int64) []chunk
	//Writes to cache, returns uncommitted data that has been trimmed.
	writeAt(p []byte, offset int64) []chunk
	//Marks a piece as committed to permanent storage.
	MarkCommitted(piece int)
	//Close the cache and free all the things
	Close()
}

type inttuple struct {
	a, b int
}

type accessTime struct {
	index int
	atime time.Time
}
type byTime []accessTime

func (a byTime) Len() int           { return len(a) }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTime) Less(i, j int) bool { return a[i].atime.Before(a[j].atime) }

//RAMCacheProvider provider creates a ram cache for each torrent.
//Each time a cache is created or closed, all cache
//are recalculated so they total <= capacity (in MiB).
type RAMCacheProvider struct {
	capacity int
	caches   map[string]*RAMCache
	m        *sync.Mutex
}

// NewRAMCacheProvider ...
func NewRAMCacheProvider(capacity int) CacheProvider {
	rc := &RAMCacheProvider{capacity, make(map[string]*RAMCache), new(sync.Mutex)}
	return rc
}

// NewCache ...
func (r *RAMCacheProvider) NewCache(infohash string, numPieces int, pieceSize int, torrentLength int64) FileCache {
	i := uint32(1)
	rc := &RAMCache{
		pieceSize:     pieceSize,
		atimes:        make([]time.Time, numPieces),
		store:         make([][]byte, numPieces),
		isBoxFull:     *NewBitset(numPieces),
		isBoxCommit:   *NewBitset(numPieces),
		isByteSet:     make([]Bitset, numPieces),
		torrentLength: torrentLength,
		cacheProvider: r,
		capacity:      &i,
		infohash:      infohash}

	r.m.Lock()
	r.caches[infohash] = rc
	r.rebalance(true)
	r.m.Unlock()
	return rc
}

//rebalance the cache capacity allocations; has to be called on each cache creation or deletion.
//'shouldTrim', if true, causes trimCommitted() to be called on all the caches. Recommended if a new cache was created
//because otherwise the old caches would stay over the new capacity until their next WriteAt happens.
func (r *RAMCacheProvider) rebalance(shouldTrim bool) {
	//Cache size is a diminishing return thing:
	//The more of it a torrent has, the less of a difference additional cache makes.
	//Thus, instead of scaling the distribution lineraly with torrent size, we'll do it by square-root
	log.Debug("Rebalancing caches...")
	var scalingTotal float64
	sqrts := make(map[string]float64)
	for i, cache := range r.caches {
		sqrts[i] = math.Sqrt(float64(cache.torrentLength))
		scalingTotal += sqrts[i]
	}

	scalingFactor := float64(r.capacity*1024*1024) / scalingTotal
	for i, cache := range r.caches {
		newCap := int(math.Floor(scalingFactor * sqrts[i] / float64(cache.pieceSize)))
		if newCap == 0 {
			newCap = 1 //Something's better than nothing!
		}
		log.Debugf("Setting cache '%s' to new capacity %v (%v MiB)", cache.infohash, newCap, float32(newCap*cache.pieceSize)/float32(1024*1024))
		cache.setCapacity(newCap)
	}

	if shouldTrim {
		for _, cache := range r.caches {
			cache.trimCommitted()
		}
	}
}

func (r *RAMCacheProvider) cacheClosed(infohash string) {
	r.m.Lock()
	delete(r.caches, infohash)
	r.rebalance(false)
	r.m.Unlock()
}

// RAMCache ...
//'pieceSize' is the size of the average piece
//'capacity' is how many pieces the cache can hold
//'actualUsage' is how many pieces the cache has at the moment
//'atime' is an array of access times for each stored box
//'store' is an array of "boxes" ([]byte of 1 piece each)
//'isBoxFull' indicates if a box entirely contains written data
//'isBoxCommit' indicates if a box has been committed to storage
//'isByteSet' for [i] indicates for box 'i' if a byte has been written to
//'torrentLength' is the number of bytes in the torrent
//'cacheProvider' is a pointer to the cacheProvider that created this cache
//'infohash' is the infohash of the torrent
type RAMCache struct {
	pieceSize     int
	capacity      *uint32 //Access only through getter/setter
	actualUsage   int
	atimes        []time.Time
	store         [][]byte
	isBoxFull     Bitset
	isBoxCommit   Bitset
	isByteSet     []Bitset
	torrentLength int64
	cacheProvider *RAMCacheProvider
	infohash      string
	m             sync.RWMutex
}

// Close ...
func (r *RAMCache) Close() {
	r.cacheProvider.cacheClosed(r.infohash)
	//We don't need to do anything else. The garbage collector will take care of it.
}

func (r *RAMCache) readAt(p []byte, off int64) []chunk {
	r.m.RLock()
	defer r.m.RUnlock()
	var unfulfilled []chunk

	boxI := int(off / int64(r.pieceSize))
	boxOff := int(off % int64(r.pieceSize))

	for i := 0; i < len(p); {
		if r.store[boxI] == nil { //definitely not in cache
			end := len(p[i:])
			if end > r.pieceSize-boxOff {
				end = r.pieceSize - boxOff
			}
			if len(unfulfilled) > 0 {
				last := unfulfilled[len(unfulfilled)-1]
				if last.i+int64(len(last.data)) == off+int64(i) {
					unfulfilled = unfulfilled[:len(unfulfilled)-1]
					i = int(last.i - off)
					end += len(last.data)
				}
			}
			unfulfilled = append(unfulfilled, chunk{off + int64(i), p[i : i+end]})
			i += end
		} else if r.isBoxFull.IsSet(boxI) { //definitely in cache
			i += copy(p[i:], r.store[boxI][boxOff:])
		} else { //Bah, do it byte by byte.
			missing := []*inttuple{&inttuple{-1, -1}}
			end := len(p[i:]) + boxOff
			if end > r.pieceSize {
				end = r.pieceSize
			}
			for j := boxOff; j < end; j++ {
				if r.isByteSet[boxI].IsSet(j) {
					p[i] = r.store[boxI][j]
				} else {
					lastIT := missing[len(missing)-1]
					if lastIT.b == i {
						lastIT.b = i + 1
					} else {
						missing = append(missing, &inttuple{i, i + 1})
					}
				}
				i++
			}
			for _, intt := range missing[1:] {
				unfulfilled = append(unfulfilled, chunk{off + int64(intt.a), p[intt.a:intt.b]})
			}
		}
		boxI++
		boxOff = 0
	}
	return unfulfilled
}

// writeAt
func (r *RAMCache) writeAt(p []byte, off int64) []chunk {
	r.m.Lock()
	defer r.m.Unlock()
	boxI := int(off / int64(r.pieceSize))
	boxOff := int(off % int64(r.pieceSize))

	for i := 0; i < len(p); {
		if r.store[boxI] == nil {
			r.store[boxI] = make([]byte, r.pieceSize)
			r.actualUsage++
		}
		copied := copy(r.store[boxI][boxOff:], p[i:])
		i += copied
		r.atimes[boxI] = time.Now()
		if copied == r.pieceSize {
			r.isBoxFull.Set(boxI)
		} else {
			if r.isByteSet[boxI].n == 0 {
				r.isByteSet[boxI] = *NewBitset(r.pieceSize)
			}
			for j := boxOff; j < boxOff+copied; j++ {
				r.isByteSet[boxI].Set(j)
			}
		}
		boxI++
		boxOff = 0
	}
	if r.actualUsage > r.getCapacity() {
		return r.trim()
	}
	return nil
}

// MarkCommitted ...
func (r *RAMCache) MarkCommitted(piece int) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.store[piece] != nil {
		r.isBoxFull.Set(piece)
		r.isBoxCommit.Set(piece)
		r.isByteSet[piece] = *NewBitset(0)
	}
}

func (r *RAMCache) removeBox(boxI int) {
	r.isBoxFull.Clear(boxI)
	r.isBoxCommit.Clear(boxI)
	r.isByteSet[boxI] = *NewBitset(0)
	r.store[boxI] = nil
	r.actualUsage--
}

func (r *RAMCache) getCapacity() int {
	return int(atomic.LoadUint32(r.capacity))
}

func (r *RAMCache) setCapacity(capacity int) {
	atomic.StoreUint32(r.capacity, uint32(capacity))
}

//Trim stuff that's already been committed
//Return true if we got underneath capacity, false if not.
func (r *RAMCache) trimCommitted() bool {
	r.m.Lock()
	defer r.m.Unlock()
	for i := 0; i < r.isBoxCommit.Len(); i++ {
		if r.isBoxCommit.IsSet(i) {
			r.removeBox(i)
		}
		if r.actualUsage <= r.getCapacity() {
			return true
		}
	}
	return false
}

//Trim excess data. Returns any uncommitted chunks that were trimmed
func (r *RAMCache) trim() []chunk {
	if r.trimCommitted() {
		return nil
	}

	var retVal []chunk

	//Still need more space? figure out what's oldest
	//RawWrite it to storage, and clear that then
	tATA := make([]accessTime, 0, r.actualUsage)

	for i, atime := range r.atimes {
		if r.store[i] != nil {
			tATA = append(tATA, accessTime{i, atime})
		}
	}

	sort.Sort(byTime(tATA))

	deficit := r.actualUsage - r.getCapacity()
	for i := 0; i < deficit; i++ {
		deadBox := tATA[i].index
		data := r.store[deadBox]
		if r.isBoxFull.IsSet(deadBox) { //Easy, the whole box has to go
			retVal = append(retVal, chunk{int64(deadBox) * int64(r.pieceSize), data})
		} else { //Ugh, we'll just trim anything unused from the start and the end, and send that.
			off := int64(0)
			endData := r.pieceSize
			//Trim out any unset bytes at the beginning
			for j := 0; j < r.pieceSize; j++ {
				if !r.isByteSet[deadBox].IsSet(j) {
					off++
				} else {
					break
				}
			}

			//Trim out any unset bytes at the end
			for j := r.pieceSize - 1; j > 0; j-- {
				if !r.isByteSet[deadBox].IsSet(j) {
					endData--
				} else {
					break
				}
			}
			retVal = append(retVal, chunk{int64(deadBox)*int64(r.pieceSize) + off, data[off:endData]})
		}
		r.removeBox(deadBox)
	}
	return retVal
}
