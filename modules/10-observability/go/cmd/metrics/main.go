// Command metrics prints a readout from the stable runtime/metrics API. The
// metric NAMES are a versioned, documented contract; the values are
// machine-dependent, so this feeds measured.txt, not a gated output.
package main

import (
	"fmt"

	obs "observability"
)

func main() {
	// Allocate to move the heap counters, then read them back.
	parts := make([]string, 1024)
	for i := range parts {
		parts[i] = "xxxxxxxx"
	}
	_ = obs.ConcatPlus(parts)

	names := []string{
		"/gc/heap/allocs:bytes",
		"/gc/heap/objects:objects",
		"/sched/goroutines:goroutines",
		"/memory/classes/total:bytes",
	}
	for _, n := range names {
		fmt.Printf("%-34s = %d\n", n, obs.ReadMetric(n))
	}
}
