package simple_test

import (
	"testing"

	"github.com/nowayhecodes/polacache/simple"
)

func TestLRU(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		if k != v {
			t.Fatalf("Evict values not equal (%v != %v)", k, v)
		}
		evictCounter++
	}

	lru, err := simple.NewLRU(128, onEvicted)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		lru.Add(i, i)
	}

	if lru.Len() != 128 {
		t.Fatalf("the length differs in size. expected %v, but got %v", 128, lru.Len())
	}

	lru.Purge()

	if lru.Len() != 0 {
		t.Fatalf("the cache wasn't purged")
	}
}

func TestLRUContains(t *testing.T) {
	lru, err := simple.NewLRU(2, nil)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	lru.Add(1, 1)
	lru.Add(2, 2)

	if !lru.Contains(1) {
		t.Errorf("1 should be in cache")
	}

	lru.Add(3, 3)

	if lru.Contains(1) {
		t.Errorf("The recent-ness of 1 shouldn't have been updated")
	}
}
