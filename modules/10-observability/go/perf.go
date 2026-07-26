// Package observability is Module 10: measuring Go with the standard library and
// toolchain. The falsifiable spine is the allocation COUNT, which is
// deterministic (testing.AllocsPerRun forces GC and counts mallocs), so an
// optimization's allocation reduction is a fact rather than a vibe. Wall-clock
// time is noisy and belongs in a measured file run through benchstat.
package observability

import "strings"

// region:builder:start

// ConcatPlus builds a string with the += operator in a loop. Strings are
// immutable, so each += allocates a NEW backing array and copies everything so
// far: quadratic copying and O(n) allocations.
func ConcatPlus(parts []string) string {
	s := ""
	for _, p := range parts {
		s += p
	}
	return s
}

// BuilderGrow builds the same string with strings.Builder, pre-sizing the
// backing array once with Grow so the whole build costs a SINGLE allocation.
func BuilderGrow(parts []string) string {
	total := 0
	for _, p := range parts {
		total += len(p)
	}
	var b strings.Builder
	b.Grow(total)
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

// region:builder:end

// hotLoop is a small CPU-bound workload with no allocation and no I/O, used as
// the subject of the CPU profile (cmd/profile) and the throughput benchmark.
func hotLoop(n int) uint64 {
	var x uint64 = 1
	for i := 0; i < n; i++ {
		x = x*6364136223846793005 + 1442695040888963407
		x ^= x >> 33
	}
	return x
}

// HotLoop is the exported entry point for the profiling command.
func HotLoop(n int) uint64 { return hotLoop(n) }
