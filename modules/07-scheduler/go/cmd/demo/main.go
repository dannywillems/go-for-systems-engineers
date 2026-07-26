// Command demo prints deterministic, machine-independent scheduler facts.
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// Goroutines are cheap and multiplexed: 100k of them exist at once here.
	// (A goroutine starts with a ~2 KB stack that grows on demand, versus a
	// ~1 MB OS-thread stack.)
	const n = 100_000
	var wg sync.WaitGroup
	release := make(chan struct{})
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-release // park until released
		}()
	}
	// Give them a moment to all park, then observe the count.
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("live goroutines with %d parked: %d\n", n, runtime.NumGoroutine())
	close(release)
	wg.Wait()
	fmt.Printf("after they finish: %d\n", runtime.NumGoroutine())

	// Asynchronous preemption (Go 1.14+): a goroutine in a tight loop with no
	// function calls used to starve the scheduler. Now the runtime preempts it,
	// so main keeps making progress. If this prints, preemption worked.
	var spin atomic.Bool
	spin.Store(true)
	go func() {
		for spin.Load() { //nolint:revive // deliberately tight; tests preemption
		}
	}()
	done := make(chan int, 1)
	go func() { done <- 42 }()
	got := <-done
	spin.Store(false)
	fmt.Printf("main progressed past a tight-loop goroutine (async preemption): got %d\n", got)
}
