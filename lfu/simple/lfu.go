package simple

import (
	"container/list"
	"encoding/binary"
	"fmt"
)

type EvictCallback func(key, value interface{})
type cachePolicy func(element *item, cacheAge float64) float64

type LFU struct {
	size        float64
	currentSize float64
	items       map[interface{}]*item
	frequently  *list.List
	onEvict     EvictCallback
	age         float64
	policy      cachePolicy
}

type item struct {
	key            interface{}
	value          interface{}
	size           float64
	hits           float64
	priorityKey    float64
	frequentlyNode *list.Element
}

type listEntry struct {
	entries     map[*item]byte
	priorityKey float64
}

// Returns a LFU with Dynamic Aging of the given size in bytes, using the default LFU eviction policy
func New(size float64, onEvict EvictCallback) *LFU {
	return &LFU{
		size:        size,
		currentSize: 0,
		items:       make(map[interface{}]*item),
		frequently:  list.New(),
		onEvict:     onEvict,
		age:         0,
		policy:      lfuPolicy,
	}
}

// Returns a new LFU with Dynamic Aging of the given size in bytes, using the GDSF eviction policy
func NewWithGDSF(size float64, onEvict EvictCallback) *LFU {
	return &LFU{
		size:        size,
		currentSize: 0,
		items:       make(map[interface{}]*item),
		frequently:  list.New(),
		onEvict:     onEvict,
		age:         0,
		policy:      gdsfPolicy,
	}
}

// Returns a new LFU with Dynamic Aging of the given size in bytes, using LFUDA eviction policy
func NewWithDA(size float64, onEvict EvictCallback) *LFU {
	return &LFU{
		size:        size,
		currentSize: 0,
		items:       make(map[interface{}]*item),
		frequently:  list.New(),
		onEvict:     onEvict,
		age:         0,
		policy:      lfuDynamicAgingPolicy,
	}
}

// Adds a value to the cache, returning true if an eviction occurs.
func (lfu *LFU) Add(key, value interface{}) bool {
	evicted := false
	if entry, ok := lfu.items[key]; ok {
		entry.value = value
		lfu.increment(entry)
	} else {
		calculatedBytes := calculateBytes(value)

		if lfu.size < calculatedBytes {
			return false
		}

		for {
			if lfu.currentSize+calculatedBytes > lfu.size {
				lfu.evict()
				evicted = true
			} else {
				break
			}
		}

		entry := new(item)
		entry.size = calculatedBytes
		entry.key = key
		entry.value = value
		lfu.items[key] = entry
		lfu.currentSize += calculatedBytes
		lfu.increment(entry)
	}
	return evicted
}

// Gets a key's value from the cache
func (lfu *LFU) Get(key interface{}) (interface{}, bool) {
	if entry, ok := lfu.items[key]; ok {
		lfu.increment(entry)
		return entry.value, true
	}
	return nil, false
}

func lfuDynamicAgingPolicy(element *item, cacheAge float64) float64 {
	return element.hits + cacheAge
}

func gdsfPolicy(element *item, cacheAge float64) float64 {
	return (element.hits / element.size) + cacheAge
}

func lfuPolicy(element *item, cacheAge float64) float64 {
	return element.hits
}

func calculateBytes(value interface{}) float64 {
	if b, ok := value.([]byte); ok {
		return float64(len(b))
	} else if b := binary.Size(value); b != -1 {
		return float64(b)
	} else {
		return float64(len([]byte(fmt.Sprintf("%v", value))))
	}
}

func (lfu *LFU) evict() bool {
	if place := lfu.frequently.Front(); place != nil {
		for entry := range place.Value.(*listEntry).entries {
			if lfu.age < entry.priorityKey {
				lfu.age = entry.priorityKey
			}

			lfu.Remove(entry.key)
			return true
		}
	}
	return false
}

func (lfu *LFU) increment(entry *item) {
	old := entry.frequentlyNode
	cursor := entry.frequentlyNode

	var next *list.Element

	if cursor == nil {
		next = lfu.frequently.Front()
	} else {
		next = cursor.Next()
	}

	entry.hits++
	entry.priorityKey = lfu.policy(entry, lfu.age)

	for {
		if next == nil || next.Value.(*listEntry).priorityKey > entry.priorityKey {
			li := new(listEntry)
			li.priorityKey = entry.priorityKey
			li.entries = make(map[*item]byte)

			if cursor != nil {
				next = lfu.frequently.InsertAfter(li, cursor)

			} else {
				next = lfu.frequently.PushFront(li)
			}
			break
		} else if next.Value.(*listEntry).priorityKey == entry.priorityKey {
			break
		} else if entry.priorityKey > next.Value.(*listEntry).priorityKey {
			cursor = next
			next = cursor.Next()
		}
	}

	entry.frequentlyNode = next
	next.Value.(*listEntry).entries[entry] = 1

	if old != nil {
		lfu.removeEntry(old, entry)
	}
}

func (lfu *LFU) removeEntry(place *list.Element, entry *item) {
	entries := place.Value.(*listEntry).entries
	delete(entries, entry)

	if len(entries) == 0 {
		lfu.frequently.Remove(place)
	}
}
