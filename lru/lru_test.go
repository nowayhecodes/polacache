package lru_test

import (
	"math/rand"
	"testing"

	"github.com/nowayhecodes/polacache/lru"
)

func BenchmarkLRU_Randon(b *testing.B) {
	lruCache, err := lru.New(8192)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)

	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hits, misses int
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			lruCache.Add(trace[i], trace[i])
		} else {
			_, ok := lruCache.Get(trace[i])
			if ok {
				hits++
			} else {
				misses++
			}
		}
	}

	b.Logf("hits: %d | misses: %d | ratio: %f", hits, misses, float64(hits)/float64(misses))
}

func TestLRU(t *testing.T) {
	evictionCounter := 0
	onEvicted := func(k, v interface{}) {
		if k != v {
			t.Fatalf("Evict values aren't equal (%v != %v)", k, v)
		}
		evictionCounter++
	}

	lruCache, err := lru.NewWithEviction(128, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		lruCache.Add(i, i)
	}

	if lruCache.Len() != 128 {
		t.Fatalf("lru has the wrong length: %v", lruCache.Len())
	}

	if evictionCounter != 128 {
		t.Fatalf("wrong eviction count: %v", evictionCounter)
	}

	lruCache.Purge()
	if lruCache.Len() != 0 {
		t.Fatalf("wrong length for lru cache: %v", lruCache.Len())
	}
}
