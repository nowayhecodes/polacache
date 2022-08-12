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

func lfuDynamicAgingPolicy(element *item, cacheAge float64) float64 {
	return element.hits + cacheAge
}

func gdsfPolicy(element *item, cacheAge float64) float64 {
	return (element.hits / element.size) + cacheAge
}

func lfuPolicy(element *item, cacheAge float64) float64 {
	return element.hits
}

func calculeBytes(value interface{}) float64 {
	if b, ok := value.([]byte); ok {
		return float64(len(b))
	} else if b := binary.Size(value); b != -1 {
		return float64(b)
	} else {
		return float64(len([]byte(fmt.Sprintf("%v", value))))
	}
}
