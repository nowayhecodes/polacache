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

func New(size int) (*Cache, error) {
	return NewWithEviction(size, nil)
}

func NewWithEviction(size int, onEvicted func(key interface{}, value interface{})) (*Cache, error) {
	lru, err := simple.NewLRU(size, simple.EvictCallback(onEvicted))

	if err != nil {
		return nil, err
	}

	c := &Cache{
		lru: lru,
	}

	return c, nil
}
