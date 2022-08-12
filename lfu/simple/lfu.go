package simple

import (
	"container/list"
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
