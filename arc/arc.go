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

// Adds a value to the cache
func (arc *ARCCache) Add(key, value interface{}) bool {
	arc.lock.Lock()
	defer arc.lock.Unlock()

	delta := 1

	if arc.recently.Contains(key) {
		if !arc.recently.Remove(key) {
			return false
		}
		return arc.frequently.Add(key, value)
	}

	if arc.frequently.Contains(key) {
		return arc.frequently.Add(key, value)
	}

	if arc.recentlyEviction.Contains(key) {
		recentlyEvictionLen := arc.recentlyEviction.Len()
		frequentlyEvictionLen := arc.frequentlyEviction.Len()

		if frequentlyEvictionLen > recentlyEvictionLen {
			delta = frequentlyEvictionLen / recentlyEvictionLen
		}

		if arc.pref+delta >= arc.size {
			arc.pref = arc.size
		} else {
			arc.pref += delta
		}

		if arc.recently.Len()+arc.frequently.Len() >= arc.size {
			arc.replace(false)
		}

		arc.recentlyEviction.Remove(key)

		return arc.frequently.Add(key, value)
	}

	if arc.frequentlyEviction.Contains(key) {
		recentEvictLen := arc.recently.Len()
		freqEvictLen := arc.frequently.Len()

		if recentEvictLen > freqEvictLen {
			delta = recentEvictLen / freqEvictLen
		}

		if delta >= arc.pref {
			arc.pref = 0
		} else {
			arc.pref -= delta
		}

		if arc.recently.Len()+arc.frequently.Len() >= arc.size {
			arc.replace(true)
		}

		arc.frequentlyEviction.Remove(key)

		return arc.frequently.Add(key, value)
	}

	if arc.recently.Len()+arc.frequently.Len() >= arc.size {
		arc.replace(false)
	}

	if arc.recentlyEviction.Len() > arc.size-arc.pref {
		arc.recentlyEviction.RemoveOldest()
	}

	if arc.frequentlyEviction.Len() > arc.pref {
		arc.frequentlyEviction.RemoveOldest()
	}

	return arc.recently.Add(key, value)
}

// This really doesn't need an explanation, right?
func (arc *ARCCache) Len() int {
	arc.lock.RLock()
	defer arc.lock.RUnlock()
	return arc.recently.Len() + arc.frequently.Len()
}

// Returns the cached keys
func (arc *ARCCache) Keys() []interface{} {
	arc.lock.RLock()
	defer arc.lock.RUnlock()

	k1 := arc.recently.Keys()
	k2 := arc.frequently.Keys()

	return append(k1, k2...)
}

// Adaptively evict from either recently or frequently
// based on current value of pref
func (arc *ARCCache) replace(freqEvictContainsKey bool) {
	recentLen := arc.recently.Len()

	if recentLen > 0 && (recentLen > arc.pref || (recentLen == arc.pref && freqEvictContainsKey)) {
		k, _, ok := arc.recently.RemoveOldest()

		if ok {
			arc.recently.Add(k, nil)
		}
	} else {
		k, _, ok := arc.frequently.RemoveOldest()
		if ok {
			arc.frequentlyEviction.Add(k, nil)
		}
	}
}
