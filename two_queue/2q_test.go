package twoqueue_test

import (
	"math/rand"
	"testing"

	TwoQ "github.com/nowayhecodes/polacache/two_queue"
)

func Benchmark2Q_Rand(b *testing.B) {
	twoq, err := TwoQ.New2Q(8192)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hits, misses int

	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			twoq.Add(trace[i], trace[i])
		} else {
			_, ok := twoq.Get(trace[i])
			if ok {
				hits++
			} else {
				misses++
			}
		}
	}
	b.Logf("hits: %d | misses: %d | ratio: %f", hits, misses, float64(hits)/float64(misses))
}

func Benchmark2Q_Freq(b *testing.B) {
	twoq, err := TwoQ.New2Q(8192)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = rand.Int63() % 16384
		} else {
			trace[i] = rand.Int63() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		twoq.Add(trace[i], trace[i])
	}
	var hits, misses int
	for i := 0; i < b.N; i++ {
		_, ok := twoq.Get(trace[i])
		if ok {
			hits++
		} else {
			misses++
		}
	}
	b.Logf("hits: %d | misses: %d | ratio: %f", hits, misses, float64(hits)/float64(misses))
}

func Test2Q(t *testing.T) {
	twoq, err := TwoQ.New2Q(128)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		twoq.Add(i, i)
	}
	if twoq.Len() != 128 {
		t.Fatalf("bad len: %v", twoq.Len())
	}

	for i, k := range twoq.Keys() {
		if v, ok := twoq.Get(k); !ok || v != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := twoq.Get(i)
		if ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := twoq.Get(i)
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		twoq.Remove(i)
		_, ok := twoq.Get(i)
		if ok {
			t.Fatalf("should be deleted")
		}
	}

	twoq.Purge()
	if twoq.Len() != 0 {
		t.Fatalf("bad len: %v", twoq.Len())
	}
	if _, ok := twoq.Get(200); ok {
		t.Fatalf("should contain nothing")
	}
}
