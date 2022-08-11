// Package Polacache provides two different implementations of LRU Cache:
// the standard LRU cache and the 2Q cache
//
// The standard LRU cache doesn't need presentations.
//
// The 2Q cache separately tracks frequently used and recently used entries,
// avoiding the burst of accesses from taking out frequently used entries.
//
// The caches here are locked while operating, meaning that is thread-safe
// for it's consumers.
package polacache
