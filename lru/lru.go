package lru

import (
	"sync"

	cache "github.com/nowayhecodes/polacache"
	"github.com/nowayhecodes/polacache/simple"
)

type Cache struct {
	lru  cache.LRUCache
	lock sync.Mutex
}

// Creates a LRU cache of a given size
// (internally returns a NewWithEviction, with a nil eviction callback)
func New(size int) (*Cache, error) {
	return NewWithEviction(size, nil)
}

// Creates a LRU cache of a given size with a given eviction callback
func NewWithEviction(size int, onEvicted func(key, value interface{})) (*Cache, error) {
	lru, err := simple.NewLRU(size, simple.EvictCallback(onEvicted))

	if err != nil {
		return nil, err
	}

	c := &Cache{
		lru: lru,
	}

	return c, nil
}

// Purges the cache
func (c *Cache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.lru.Purge()
}

// Adds an entry to the cache, returning the if an eviction occured
func (c *Cache) Add(key, value interface{}) (evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.lru.Add(key, value)
}
