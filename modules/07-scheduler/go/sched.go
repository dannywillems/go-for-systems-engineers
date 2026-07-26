// Package sched probes Go's GMP scheduler: goroutines (G) multiplexed over OS
// threads (M) through a fixed set of logical processors (P = GOMAXPROCS), with
// work stealing, asynchronous preemption (since Go 1.14, so a tight loop no
// longer starves the scheduler), and a netpoller that parks goroutines blocked
// on I/O without consuming an M.
package sched

import (
	"math"
	"sort"
	"sync"
	"time"
)

// region:pool:start

// cpuWork is a small, fixed CPU-bound task (no allocation, no I/O), so a pool of
// these measures pure scheduling behavior.
func cpuWork() float64 {
	x := 0.0
	for i := range 2000 {
		x += math.Sqrt(float64(i))
	}
	return x
}

// PoolLatencies submits nTasks to a pool of `workers` goroutines and returns
// each task's LATENCY: the time from submission to completion. With more tasks
// than logical processors, tasks queue, and the tail (p99) shows the scheduler's
// queueing under CPU-bound load.
func PoolLatencies(workers, nTasks int) []time.Duration {
	tasks := make(chan time.Time, nTasks)
	lat := make([]time.Duration, nTasks)
	var wg sync.WaitGroup
	var idx int64
	var mu sync.Mutex

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for submitted := range tasks {
				_ = cpuWork()
				d := time.Since(submitted)
				mu.Lock()
				lat[idx] = d
				idx++
				mu.Unlock()
			}
		}()
	}
	for range nTasks {
		tasks <- time.Now()
	}
	close(tasks)
	wg.Wait()
	return lat
}

// region:pool:end

// Percentile returns the p-th percentile (0..100) of the durations.
func Percentile(ds []time.Duration, p float64) time.Duration {
	if len(ds) == 0 {
		return 0
	}
	s := append([]time.Duration(nil), ds...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	rank := int(math.Ceil(p/100*float64(len(s)))) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= len(s) {
		rank = len(s) - 1
	}
	return s[rank]
}
