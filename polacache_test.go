package polacache

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var polacache = New(1 * time.Minute)

var testItem = Item{
	Key:   "polatest",
	Value: "testifying polacache",
}

func TestPolacache_GetSet(t *testing.T) {
	requires := require.New(t)

	polacache.Set(testItem, time.Now().Add(1*time.Hour).Unix())

	it, err := polacache.Get(testItem.Key)
	requires.NoError(err)
	requires.Equal(testItem.Value, it)

	polacache.stopCleanup()
}

func TestPolacache_ErrorMessage(t *testing.T) {
	requires := require.New(t)

	polacache.Delete(testItem.Key)

	it, err := polacache.Get(testItem.Key)
	requires.EqualError(err, "key polatest not in cache")
	requires.Equal(Item{}, it)

	polacache.stopCleanup()
}

func BenchmarkPolacache(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			polacache.Set(Item{
				Key:   strconv.Itoa(int(rand.Int63())),
				Value: "polacache test",
			}, time.Now().Add(1*time.Hour).Unix())
		}
	})

	polacache.stopCleanup()
}
