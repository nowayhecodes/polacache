package lru

import (
	"container/list"
	"fmt"
)

type EvictCallback func(key interface{}, value interface{})

type LRU struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EvictCallback
}

type entry struct {
	key   interface{}
	value interface{}
}

// Initialize a new LRU of the given size
func NewLRU(size int, onEvict EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, fmt.Errorf("must provide a positive integer for LRU size")
	}

	lru := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
	}
	return lru, nil
}

// Clears the cache
func (lru *LRU) Purge() {
	for k, v := range lru.items {
		if lru.onEvict != nil {
			lru.onEvict(k, v.Value.(*entry).value)
		}
		delete(lru.items, k)
	}
	lru.evictList.Init()
}

// Adds stuff to the cache returning true if an eviction occurs
func (lru *LRU) Add(key, value interface{}) (evicted bool) {
	if item, ok := lru.items[key]; ok {
		lru.evictList.MoveToFront(item)
		item.Value.(*entry).value = value
		return false
	}

	item := &entry{key, value}
	entry := lru.evictList.PushFront(item)
	lru.items[key] = entry

	evict := lru.evictList.Len() > lru.size

	if evict {
		lru.removeOldest()
	}

	return evict
}

// Gets a key's value from cache
func (lru *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if item, ok := lru.items[key]; ok {
		lru.evictList.MoveToFront(item)
		v, o := item.Value.(*entry)

		if v == nil || !o {
			return nil, false
		}

		return v.value, true
	}
	return
}

// Checks if a key is in the cache
func (lru *LRU) Contains(key interface{}) (ok bool) {
	_, ok = lru.items[key]
	return ok
}

// Peek returns a key's value from the cache (or undefined)
// without update the "recenty"-ness of the given key
func (lru *LRU) Peek(key interface{}) (value interface{}, ok bool) {
	var item *list.Element
	if item, ok = lru.items[key]; ok {
		return item.Value.(*entry).value, true
	}

	return nil, ok
}

// Removes the oldest item from the cache
func (lru *LRU) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	item := lru.evictList.Back()
	if item != nil {
		lru.removeElement(item)
		kv := item.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// Returns the oldest entry
func (lru *LRU) GetOldest() (key interface{}, value interface{}, ok bool) {
	item := lru.evictList.Back()
	if item != nil {
		kv := item.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// Returns a slice of the keys in cache (oldest -> newest)
func (lru *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(lru.items))
	i := 0

	for item := lru.evictList.Back(); item != nil; item = item.Prev() {
		keys[i] = item.Value.(*entry).key
		i++
	}
	return keys
}

func (lru *LRU) removeOldest() {
	item := lru.evictList.Back()
	if item != nil {
		lru.removeElement(item)
	}
}

func (lru *LRU) removeElement(element *list.Element) {
	lru.evictList.Remove(element)
	kv := element.Value.(*entry)
	delete(lru.items, kv.key)

	if lru.onEvict != nil {
		lru.onEvict(kv.key, kv.value)
	}
}
