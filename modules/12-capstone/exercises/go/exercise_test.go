package exercises

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func TestBoundedFetchCorrectAndBounded(t *testing.T) {
	const capacity = 16
	var mu sync.Mutex
	m := make(map[int]int, capacity)
	sem := NewSem(4)
	var calls atomic.Int64
	fetch := func(k int) int { calls.Add(1); return k * k }

	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				key := (w*500 + i) % 100
				v, err := BoundedFetch(context.Background(), &mu, m, capacity, sem, key, fetch)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if v != key*key {
					t.Errorf("key %d: got %d, want %d", key, v, key*key)
					return
				}
			}
		}(w)
	}
	wg.Wait()

	mu.Lock()
	got := len(m)
	mu.Unlock()
	if got > capacity {
		t.Fatalf("cache exceeded capacity: %d > %d", got, capacity)
	}
	if calls.Load() >= 8*500 {
		t.Fatalf("cache did nothing: %d backend calls", calls.Load())
	}
}

func TestAcquireRespectsContext(t *testing.T) {
	sem := NewSem(1)
	if err := sem.Acquire(context.Background()); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	// Semaphore is now full; a cancelled context must make Acquire return.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sem.Acquire(ctx); err == nil {
		t.Fatal("Acquire on a full semaphore with a cancelled context returned nil")
	}
}
