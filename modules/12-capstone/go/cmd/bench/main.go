// Command bench drives the concurrent cache under load and reports throughput,
// latency percentiles, and how much the cache reduced backend calls. Timings are
// non-deterministic; captured into measured.txt by scripts/capstone-bench.sh.
package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	capstone "github.com/dannywillems/go-for-systems-engineers/modules/12-capstone/go"
)

const (
	capacity    = 256
	maxInflight = 32
	keys        = 256 // == capacity: after warmup the hot set is fully cached
	workers     = 64
	perWorker   = 10_000
	latency     = 100 * time.Microsecond
)

func main() {
	b := capstone.NewBackend(latency)
	c := capstone.NewCache(capacity, maxInflight, b)
	ctx := context.Background()

	lat := make([][]time.Duration, workers)
	var wg sync.WaitGroup
	start := time.Now()
	for w := range workers {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			ls := make([]time.Duration, perWorker)
			//nolint:gosec // G115: w is a worker index in [0,workers), never negative.
			seed := uint64(w)*2654435761 | 1
			for i := range perWorker {
				seed = seed*6364136223846793005 + 1442695040888963407
				key := int(seed>>33) % keys
				s := time.Now()
				_, _ = c.Get(ctx, key)
				ls[i] = time.Since(s)
			}
			lat[w] = ls
		}(w)
	}
	wg.Wait()
	elapsed := time.Since(start)

	all := make([]time.Duration, 0, workers*perWorker)
	for _, ls := range lat {
		all = append(all, ls...)
	}
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
	n := len(all)
	pc := func(p float64) time.Duration {
		return all[min(int(p/100*float64(n)), n-1)].Round(100)
	}
	fmt.Printf("Go     %dk gets/%dw: %.0f kops/s  p50=%v p99=%v p999=%v  backend=%.1f%% of gets\n",
		n/1000, workers, float64(n)/elapsed.Seconds()/1000,
		pc(50), pc(99), pc(99.9), 100*float64(b.Calls())/float64(n))
}
