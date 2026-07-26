package conc

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSynchronizedCountersAreCorrect(t *testing.T) {
	if got := AtomicInc(); got != want {
		t.Errorf("AtomicInc = %d, want %d", got, want)
	}
	if got := MutexInc(); got != want {
		t.Errorf("MutexInc = %d, want %d", got, want)
	}
}

// Contended-counter benchmark: many goroutines incrementing one counter.
// b.RunParallel spreads iterations across GOMAXPROCS goroutines, so this
// measures the synchronization cost under real contention.

func BenchmarkAtomic(b *testing.B) {
	var c atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Add(1)
		}
	})
	_ = c.Load()
}

func BenchmarkMutex(b *testing.B) {
	var (
		mu sync.Mutex
		c  int64
	)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			c++
			mu.Unlock()
		}
	})
	_ = c
}
