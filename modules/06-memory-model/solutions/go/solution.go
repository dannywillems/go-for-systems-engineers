// Package solutions is the corrigé for Module 06. Run via `make solutions`
// (and `go test -race` locally to confirm race-freedom).
package solutions

import (
	"sync"
	"sync/atomic"
)

func Increment(workers, per int) int {
	var counter atomic.Int64
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range per {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
	return int(counter.Load())
}

func Merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, ch := range chans {
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}
	// The sender closes; a separate goroutine waits for all fan-in producers.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
