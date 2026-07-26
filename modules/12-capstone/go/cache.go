// Package capstone is the cross-language capstone: a concurrent bounded cache
// over a simulated slow backend, driven by concurrent load. It exercises the
// whole runtime story (Modules 05-07) at once: shared mutable state under a
// lock, bounded memory (eviction), backpressure (a semaphore capping concurrent
// backend fetches), and graceful shutdown (context cancellation).
//
// The HTTP transport is intentionally omitted so the five implementations stay
// dependency-free and directly comparable; the concurrent cache is the
// substance a migration decision turns on, not the request plumbing.
package capstone

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Backend is the simulated slow dependency. Each fetch sleeps, and every call is
// counted, so a run can report how much the cache reduced backend load.
type Backend struct {
	latency time.Duration
	calls   atomic.Int64
}

func NewBackend(latency time.Duration) *Backend { return &Backend{latency: latency} }

func (b *Backend) fetch(key int) int {
	b.calls.Add(1)
	time.Sleep(b.latency)
	return key * key
}

func (b *Backend) Calls() int64 { return b.calls.Load() }

// region:cache:start

// Cache is a concurrent, bounded key->value cache over a Backend. It caps the
// number of entries (evicting an arbitrary one on overflow) and the number of
// concurrent backend fetches (a semaphore = backpressure).
type Cache struct {
	entries  map[int]int
	backend  *Backend
	sem      chan struct{} // backpressure: bounds concurrent fetches
	mu       sync.Mutex
	capacity int
}

func NewCache(capacity, maxInflight int, b *Backend) *Cache {
	return &Cache{
		entries:  make(map[int]int, capacity),
		backend:  b,
		sem:      make(chan struct{}, maxInflight),
		capacity: capacity,
	}
}

// Get returns the cached value or fetches it from the backend, respecting the
// context (graceful shutdown) and the in-flight backpressure limit.
func (c *Cache) Get(ctx context.Context, key int) (int, error) {
	c.mu.Lock()
	if v, ok := c.entries[key]; ok {
		c.mu.Unlock()
		return v, nil
	}
	c.mu.Unlock()

	// Backpressure: block until a fetch slot is free, or the context is done.
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	v := c.backend.fetch(key)

	c.mu.Lock()
	if len(c.entries) >= c.capacity {
		for k := range c.entries { // evict one arbitrary entry
			delete(c.entries, k)
			break
		}
	}
	c.entries[key] = v
	c.mu.Unlock()
	return v, nil
}

// region:cache:end

// Len reports the current number of cached entries (never exceeds capacity).
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
