package simple

type LFUCache interface {
	Add(key, value interface{}) bool
	Get(key interface{}) (value interface{}, ok bool)
	Contains(key interface{}) (ok bool)
	Peek(key interface{}) (valeu interface{}, ok bool)
	Remove(key interface{}) bool
	Keys() []interface{}
	Len() int
	Size() float64
	Purge()

	// Returns the cache's current age factor
	Age() float64
}
