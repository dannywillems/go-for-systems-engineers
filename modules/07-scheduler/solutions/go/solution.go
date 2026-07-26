// Package solutions is the corrigé for Module 07. Run via `make solutions`
// (and `go test -race`).
package solutions

import (
	"math"
	"sort"
	"sync"
)

func ParallelMap(xs []int, workers int, f func(int) int) []int {
	out := make([]int, len(xs))
	jobs := make(chan int, len(xs))
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				out[i] = f(xs[i]) // distinct index per job: no race
			}
		}()
	}
	for i := range xs {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return out
}

func Percentile(xs []int, p float64) int {
	if len(xs) == 0 {
		return 0
	}
	s := append([]int(nil), xs...)
	sort.Ints(s)
	rank := int(math.Ceil(p/100*float64(len(s)))) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= len(s) {
		rank = len(s) - 1
	}
	return s[rank]
}
