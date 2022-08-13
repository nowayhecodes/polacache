package simple_test

import (
	"testing"

	LFU "github.com/nowayhecodes/polacache/lfu/simple"
)

func TestLFUDA(t *testing.T) {
	lfu := LFU.NewWithDA(2, nil)
	lfu.Add("a", "a")
	if v, _ := lfu.Get("a"); v != "a" {
		t.Errorf("Value was not saved: %v != 'a'", v)
	}
	if l := lfu.Len(); l != 1 {
		t.Errorf("Length was not updated: %v != 1", l)
	}

	lfu.Add("b", "b")
	if v, _ := lfu.Get("b"); v != "b" {
		t.Errorf("Value was not saved: %v != 'b'", v)
	}
	if l := lfu.Len(); l != 2 {
		t.Errorf("Length was not updated: %v != 2", l)
	}

	if v, ok := lfu.Get("b"); !ok {
		t.Errorf("Value was improperly evicted: %v != 'b'", v)
	}

	if ok := lfu.Remove("a"); !ok {
		t.Errorf("Item was not removed: a")
	}
	if v, _ := lfu.Get("a"); v != nil {
		t.Errorf("Value was not removed: %v", v)
	}
	if l := lfu.Len(); l != 1 {
		t.Errorf("Length was not updated: %v != 1", l)
	}
}

func TestEvictGDSF(t *testing.T) {
	lfu := LFU.NewWithGDSF(10, nil)
	lfu.Add("a", "aaaaaaaa")
	lfu.Add("b", "b")
	lfu.Add("c", "c")

	if lfu.Size() != 10 {
		t.Errorf("cache should have size 10 bytes at this point: %f", lfu.Size())
	}

	for i := 0; i < 10; i++ {
		lfu.Get("a")
	}

	for j := 0; j < 10; j++ {
		lfu.Add(j, j)
	}

	if ok := lfu.Contains("a"); ok {
		t.Errorf("cache should not have contained key a now")
	}

	lfu.Add("a", "aaaaaaaa")
	for i := 0; i < 50; i++ {
		lfu.Get("a")
	}

	for j := 0; j < 10; j++ {
		lfu.Add(j, j)
	}

	if ok := lfu.Contains("a"); !ok {
		t.Errorf("cache should have contained key a")
	}

	for j := 0; j < 10; j++ {
		lfu.Add(j, j)
	}

	if ok := lfu.Contains("a"); ok {
		t.Errorf("cache should NOT have contained key a now")
	}
}
