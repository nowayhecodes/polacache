package arc

import (
	"sync"

	"github.com/nowayhecodes/polacache"
	lru "github.com/nowayhecodes/polacache/simple"
)

// ARC Cache is a thread-safe fixed size Adaptative Replacement Cache.
// ARC enhances the standard LRU cache algorithm by tracking both
// frequency and recency of use, avoiding a burst in access to new
// entries from evicting the frequently used older entries.
type ARCCache struct {
	size int
	pref int // pref is the dynamic preference towards T1 or T2
	lock sync.RWMutex

	T1         polacache.LRUCache
	evictForT1 polacache.LRUCache

	T2         polacache.LRUCache
	evictForT2 polacache.LRUCache
}

// Creates an ARC of given size
func NewARC(size int) (*ARCCache, error) {

	evictForT1, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	evictForT2, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	t1, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	t2, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	arc := &ARCCache{
		size:       size,
		pref:       0,
		T1:         t1,
		evictForT1: evictForT1,
		T2:         t2,
		evictForT2: evictForT2,
	}

	return arc, nil
}
