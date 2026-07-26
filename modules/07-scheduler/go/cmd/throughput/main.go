// Command throughput runs a CPU-bound parallel workload (sum of sqrt over a big
// range) across GOMAXPROCS goroutines and self-times it. This is the Go entry in
// the cross-language throughput comparison; captured into measured.txt.
package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

const total = 400_000_000

func main() {
	w := runtime.GOMAXPROCS(0)
	chunk := total / w
	sums := make([]float64, w)
	var wg sync.WaitGroup
	t0 := time.Now()
	for k := range w {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			s := 0.0
			for i := k * chunk; i < (k+1)*chunk; i++ {
				s += math.Sqrt(float64(i))
			}
			sums[k] = s
		}(k)
	}
	wg.Wait()
	acc := 0.0
	for _, s := range sums {
		acc += s
	}
	fmt.Printf("Go     sqrt-sum %dM / %d goroutines: %d ms\n",
		total/1_000_000, w, time.Since(t0).Milliseconds())
}
