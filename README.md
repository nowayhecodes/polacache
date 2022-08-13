<div align="center">
    <img src="./asset/polacache.png"  width="300" alt="iters" />
</div>

<div align="center">
    <p style="font-size: 5.5rem;">Polacache</p>    
</div>

<div align="center" style="margin-top: -4rem;">
    <p style="font-size: 1.5rem;">Some caching algorithms in Go</p>    
</div>

### This package implements some caching algorithms as discribed below:

- LRU (Least Recently Used): Keeps track of the least recently used item in the cache and discards it. The LRU eviction algorithm evicts the page from the buffer which has not been accessed for the longest.

- LFU (Least Frequently Used): Does almost the same as LRU, but tracks the least frequently used. This also implements a Dynamic Aging and a GreedyDual-Size with Frequency cache policy.

- 2Q Cache: 2Q addresses the above-illustrated issues by introducing parallel buffers and supporting queues. 2Q algorithm works with two buffers. A primary LRU buffer and a secondary FIFO buffer. Instead of considering just recency as a factor, 2Q also considers access frequency while making the decision to ensure the page that is really warm gets a place in the LRU cache. It admits only hot pages to the main buffer and tests every page for a second reference.

