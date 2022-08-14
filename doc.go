// Polacache is deadly a simple and thread-safe map cache.
//
// In it's constructor, you set a cleanupInterval, which launchs
// a goroutine to perform the cleanup loop.
//
//
// Example:
//
// 	package main
//
// 	import (
// 		"time"
//
// 		pc "github.com/nowayhecodes/polacache"
// 	)
//
// 	func main() {
// 		cache := pc.New(1 * time.Minute)
//
// 		exampleItem := pc.Item{
// 			Key:   "example",
// 			Value: 42,
// 		}
//
// 		cache.Set(exampleItem, time.Now().Add(1*time.Hour).Unix())
// 		cache.Get(exampleItem.Key)
// 		cache.Delete(exampleItem.Key)
//
// 	}
//
// 	...
package polacache
