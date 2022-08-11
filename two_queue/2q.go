package twoqueue

import (
	"fmt"
	"sync"

	"github.com/nowayhecodes/polacache"
	lru "github.com/nowayhecodes/polacache/simple"
)

// Default2QRecentRatio is the ratio of the 2Q cache dedicated
// to recently added entries that have only been accessed once.
const DEFAULT_2Q_RECENT_RATIO = 0.25

// Default2QGhostEntries is the default ratio of ghost
// entries kept to track entries recently evicted
const DEFAULT_2Q_GHOST_ENTRIES = 0.50

// TwoQueueCache is a thread-safe fixed size 2Q cache.
// 2Q is an enhancement over the standard LRU cache
// in that it tracks both frequently and recently used
// entries separately. This avoids a burst in access to new
// entries from evicting frequently used entries. It adds some
// additional tracking overhead to the standard LRU cache, and is
// computationally about 2x the cost, and adds some metadata over
// head. The ARCCache is similar, but does not require setting any
// parameters.
type TwoQueueCache struct {
	size        int
	recentSize  int
	recent      polacache.LRUCache
	frequent    polacache.LRUCache
	recentEvict polacache.LRUCache
	lock        sync.RWMutex
}

// Creates a 2Q cache of the given size with some default params
func New2Q(size int) (*TwoQueueCache, error) {
	return New2QWithParams(size, DEFAULT_2Q_RECENT_RATIO, DEFAULT_2Q_GHOST_ENTRIES)
}

// Creates a 2Q cache with the provided params
func New2QWithParams(size int, recentRatio float64, ghostRation float64) (*TwoQueueCache, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size")
	}

	if recentRatio < 0.0 || recentRatio > 1.0 {
		return nil, fmt.Errorf("invalid recent ratio")
	}

	if ghostRation < 0.0 || ghostRation > 1.0 {
		return nil, fmt.Errorf("invalid ghost ratio")
	}

	recentSize := int(float64(size) * recentRatio)
	evictSize := int(float64(size) * ghostRation)

	recent, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	frequent, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	recentEvict, err := lru.NewLRU(evictSize, nil)
	if err != nil {
		return nil, err
	}

	twoQ := &TwoQueueCache{
		size:        size,
		recentSize:  recentSize,
		recent:      recent,
		frequent:    frequent,
		recentEvict: recentEvict,
	}

	return twoQ, nil
}

// Adds a value to the cache
func (twoq *TwoQueueCache) Add(key, value interface{}) bool {
	twoq.lock.Lock()
	defer twoq.lock.Unlock()

	if twoq.frequent.Contains(key) {
		return twoq.frequent.Add(key, value)
	}

	if twoq.recent.Contains(key) {
		twoq.recent.Remove(key)
		return twoq.frequent.Add(key, value)
	}

	if twoq.recentEvict.Contains(key) {
		twoq.ensureSpace(true)
		twoq.recentEvict.Remove(key)
		return twoq.frequent.Add(key, value)
	}

	twoq.ensureSpace(false)
	return twoq.recent.Add(key, value)
}

// Gets a key's value from the cache
func (twoq *TwoQueueCache) Get(key interface{}) (value interface{}, ok bool) {
	twoq.lock.Lock()
	defer twoq.lock.Unlock()

	if val, ok := twoq.frequent.Get(key); ok {
		return val, ok
	}

	if val, ok := twoq.recent.Peek(key); ok {
		twoq.recent.Remove(key)
		twoq.frequent.Add(key, val)
		return val, ok
	}

	return nil, false
}

// This really doesn't need an explanation, right?
func (twoq *TwoQueueCache) Len() int {
	twoq.lock.RLock()
	defer twoq.lock.RUnlock()
	return twoq.recent.Len() + twoq.frequent.Len()
}

// Returns a slice of the keys in the cache.
func (twoq *TwoQueueCache) Keys() []interface{} {
	twoq.lock.RLock()
	defer twoq.lock.RUnlock()

	k1 := twoq.frequent.Keys()
	k2 := twoq.recent.Keys()
	return append(k1, k2...)
}

// Removes the given key from the cache
func (twoq *TwoQueueCache) Remove(key interface{}) bool {
	twoq.lock.Lock()
	defer twoq.lock.Unlock()

	if twoq.frequent.Remove(key) {
		return true
	}

	if twoq.recent.Remove(key) {
		return true
	}

	return twoq.recentEvict.Remove(key)
}

func (twoq *TwoQueueCache) ensureSpace(recentEvict bool) {
	recentLen := twoq.recent.Len()
	frequentLen := twoq.frequent.Len()

	if recentLen+frequentLen < twoq.size {
		return
	}

	if recentLen > 0 &&
		(recentLen > twoq.recentSize || (recentLen == twoq.recentSize && !recentEvict)) {
		key, _, _ := twoq.recent.RemoveOldest()
		twoq.recentEvict.Add(key, nil)
		return
	}

	twoq.frequent.RemoveOldest()
}
