// Command latency measures task-completion latency percentiles under CPU-bound
// load, at the natural pool size (GOMAXPROCS workers) versus an oversubscribed
// one, to show the scheduler's queueing tail. Non-deterministic; captured into
// measured.txt by scripts/sched-bench.sh.
package main

import (
	"fmt"
	"runtime"

	sched "github.com/dannywillems/go-for-systems-engineers/modules/07-scheduler/go"
)

const nTasks = 200_000

func report(label string, workers int) {
	lat := sched.PoolLatencies(workers, nTasks)
	fmt.Printf("%-22s p50=%-8v p99=%-8v p999=%-8v max=%v\n",
		label,
		sched.Percentile(lat, 50).Round(1e3),
		sched.Percentile(lat, 99).Round(1e3),
		sched.Percentile(lat, 99.9).Round(1e3),
		sched.Percentile(lat, 100).Round(1e3))
}

func main() {
	p := runtime.GOMAXPROCS(0)
	fmt.Printf("GOMAXPROCS=%d, %d CPU-bound tasks\n", p, nTasks)
	report(fmt.Sprintf("workers=GOMAXPROCS(%d)", p), p)
	report("workers=2 (limited)", 2)
	report(fmt.Sprintf("workers=4x(%d) oversub", 4*p), 4*p)
}
