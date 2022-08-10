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
	lock sync.RWMutex

	// Total capacity of the cache
	size int

	// pref is the dynamic preference towards "recently" or "frequently"
	pref int

	// LRU for recently accessed items
	recently polacache.LRUCache

	// LRU evictions for recently accessed items
	recentlyEviction polacache.LRUCache

	// LRU for frequently accessed items
	frequently polacache.LRUCache

	// LRU evictions for frequently accessed items
	frequentlyEviction polacache.LRUCache
}

// Creates an ARC of given size
func NewARC(size int) (*ARCCache, error) {

	recentEvict, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	frequentEvict, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	recent, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	frequent, err := lru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	arc := &ARCCache{
		size:               size,
		pref:               0,
		recently:           recent,
		recentlyEviction:   recentEvict,
		frequently:         frequent,
		frequentlyEviction: frequentEvict,
	}

	return arc, nil
}

// Gets the given value from the cache
func (arc *ARCCache) Get(key interface{}) (value interface{}, ok bool) {
	arc.lock.Lock()
	defer arc.lock.Unlock()

	// Promotes a value to frequently used LRU,
	// if the recently used LRU contains the given value
	if val, ok := arc.recently.Peek(key); ok {
		arc.recently.Remove(key)
		arc.frequently.Add(key, val)

		return val, ok
	}

	// Checks if frequently used LRU contains the given value
	if val, ok := arc.frequently.Get(key); ok {
		return val, ok
	}

	return nil, false
}
