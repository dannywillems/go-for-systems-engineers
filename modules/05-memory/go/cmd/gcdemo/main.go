// Command gcdemo runs a fixed allocation workload and reports throughput and GC
// behavior. Run it under different GOGC to see the throughput-vs-RSS trade:
//
//	GOGC=100 go run ./cmd/gcdemo   # default: more GC, less memory
//	GOGC=800 go run ./cmd/gcdemo   # fewer GCs, more HeapSys, often faster
package main

import (
	"fmt"
	"runtime"
	"time"
)

const n = 50_000_000

// alloc is noinline so the returned pointer genuinely escapes to the heap; each
// [8]int is short-lived, creating steady GC pressure.
//
//go:noinline
func alloc(i int) *[8]int {
	a := new([8]int)
	a[0] = i
	return a
}

func main() {
	var acc int
	start := time.Now()
	for i := range n {
		acc += alloc(i)[0]
	}
	elapsed := time.Since(start).Milliseconds()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("alloc %dM: %d ms  NumGC=%d  HeapSys=%d MiB  (acc=%d)\n",
		n/1_000_000, elapsed, m.NumGC, m.HeapSys>>20, acc)
}
