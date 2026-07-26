package capstone

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCacheCorrectAndBounded(t *testing.T) {
	b := NewBackend(0)
	c := NewCache(16, 4, b)
	ctx := context.Background()

	// Concurrent load over a key space larger than the capacity.
	var wg sync.WaitGroup
	for w := range 8 {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := range 500 {
				key := (w*500 + i) % 100
				v, err := c.Get(ctx, key)
				if err != nil {
					t.Errorf("Get: %v", err)
					return
				}
				if v != key*key {
					t.Errorf("Get(%d) = %d, want %d", key, v, key*key)
					return
				}
			}
		}(w)
	}
	wg.Wait()

	if c.Len() > 16 {
		t.Errorf("cache size %d exceeds capacity 16", c.Len())
	}
	// The cache must have absorbed some load: fewer backend calls than requests.
	if b.Calls() >= 8*500 {
		t.Errorf("backend called %d times for %d requests; cache did nothing", b.Calls(), 8*500)
	}
}

func TestGracefulShutdown(t *testing.T) {
	b := NewBackend(50 * time.Millisecond)
	c := NewCache(4, 1, b) // maxInflight=1 so the 2nd fetch must wait
	ctx, cancel := context.WithCancel(context.Background())

	// Occupy the single fetch slot, then cancel; a second Get must observe the
	// cancellation instead of blocking forever.
	go func() { _, _ = c.Get(context.Background(), 1) }()
	time.Sleep(5 * time.Millisecond)
	cancel()
	if _, err := c.Get(ctx, 2); err == nil {
		t.Error("expected context cancellation error under backpressure")
	}
}
