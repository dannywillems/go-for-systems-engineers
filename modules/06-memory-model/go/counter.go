// Package conc demonstrates the Go Memory Model: races are undefined behavior,
// the -race detector is dynamic (it only sees races that actually occur), and
// Go — unlike Rust or Swift 6 — does NOT reject racy programs at compile time.
//
// The Go Memory Model (https://go.dev/ref/mem) defines happens-before via
// channel send/receive, mutex Lock/Unlock, and sync/atomic. A read that is not
// ordered-before-or-after a concurrent write by one of these is a data race,
// and "a program with data races has undefined behavior" — including torn reads
// of multiword values (interfaces, slices).
package conc

import (
	"sync"
	"sync/atomic"
)

const (
	goroutines = 8
	perG       = 100_000
	want       = goroutines * perG
)

// region:counters:start

// RacyInc increments a shared int from many goroutines with NO synchronization.
// This is a data race: the result is nondeterministic AND the program has
// undefined behavior. -race flags it (see the captured report); the compiler
// does not. Rust and Swift 6 reject the analogous program at compile time.
func RacyInc() int {
	var counter int // shared, unguarded
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perG {
				counter++ // read-modify-write with no happens-before edge
			}
		}()
	}
	wg.Wait()
	return counter
}

// AtomicInc uses sync/atomic, which establishes the happens-before edges the
// memory model requires. Always returns exactly want.
func AtomicInc() int64 {
	var counter atomic.Int64
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perG {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
	return counter.Load()
}

// MutexInc uses a mutex; Lock/Unlock also establish happens-before.
func MutexInc() int {
	var (
		mu      sync.Mutex
		counter int
	)
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perG {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return counter
}

// region:counters:end
