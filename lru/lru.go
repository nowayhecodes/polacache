package lru

import (
	"sync"

	cache "github.com/nowayhecodes/polacache"
	"github.com/nowayhecodes/polacache/simple"
)

type Cache struct {
	lru  cache.LRUCache
	lock sync.RWMutex
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

// Gets a key's value from cache
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.lru.Get(key)
}

// Checks if a given is in the cache
func (c *Cache) Contains(key interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.lru.Contains(key)
}

// Returns the key value without updating
// the recentness of the key
func (c *Cache) Peek(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.lru.Peek(key)
}

// Reaps of the provided key from the cache
func (c *Cache) Remove(key interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.lru.Remove(key)
}

// Reaps of the oldest item from the cache
func (c *Cache) RemoveOldest() (interface{}, interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	key, value, ok := c.lru.RemoveOldest()
	return key, value, ok
}

// Gets the oldest item in the cache
func (c *Cache) GetOldest() (interface{}, interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	key, value, ok := c.lru.GetOldest()
	return key, value, ok
}

//Returns a slice of the keys in the cache, from oldest to newest
func (c *Cache) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.lru.Keys()
}

// This really doesn't need an explanation, right?
func (c *Cache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.lru.Len()
}
