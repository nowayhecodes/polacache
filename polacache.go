package polacache

import (
	"fmt"
	"sync"
	"time"
)

type item struct {
	key   string
	value interface{}
}

type cachedItem struct {
	item      item
	expiresAt int64
}

type cache struct {
	stop chan struct{}

	wg    sync.WaitGroup
	lock  sync.RWMutex
	items map[interface{}]cachedItem
}

func New(cleanupInterval time.Duration) *cache {
	c := &cache{
		items: make(map[interface{}]cachedItem),
		stop:  make(chan struct{}),
	}

	c.wg.Add(1)

	go func(cleanupInterval time.Duration) {
		defer c.wg.Done()
		c.cleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return c
}

func (c *cache) Set(i item, expiresAt int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[i.key] = cachedItem{
		item:      i,
		expiresAt: expiresAt,
	}
}

func (c *cache) Get(key interface{}) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cached, ok := c.items[key]
	if !ok {
		return item{}, fmt.Errorf("key %v not in cache", key)
	}

	return cached.item.value, nil
}

func (c *cache) Delete(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.items, key)
}

func (c *cache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-t.C:
			c.lock.Lock()
			for uid, cached := range c.items {
				if cached.expiresAt <= time.Now().Unix() {
					delete(c.items, uid)
				}
			}
			c.lock.Unlock()
		}
	}

}

func (c *cache) stopCleanup() {
	close(c.stop)
	c.wg.Wait()
}
