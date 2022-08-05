package twoqueue

// Default2QRecentRatio is the ratio of the 2Q cache dedicated
// to recently added entries that have only been accessed once.
const DEFAULT_2Q_RECENT_RATIO = 0.25

// Default2QGhostEntries is the default ratio of ghost
// entries kept to track entries recently evicted
const DEFAULT_2Q_GHOST_ENTRIES = 0.50

// TwoQueueCache is a thread-safe fixed size 2Q cache.
// 2Q is an enhancement over the standard LRU cache
// in that it tracks both frequently and recently used
// entries separately. This avoids a burst in access to new
// entries from evicting frequently used entries. It adds some
// additional tracking overhead to the standard LRU cache, and is
// computationally about 2x the cost, and adds some metadata over
// head. The ARCCache is similar, but does not require setting any
// parameters.
type TwoQueueCache struct {
	size       int
	recentSize int
}
