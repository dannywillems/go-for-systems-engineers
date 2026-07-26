// Package exercises: Module 10 reader tasks. RED until you implement the stubs.
// The point is a SINGLE-allocation join, verified by testing.AllocsPerRun.
package exercises

// TODO(reader): join parts with sep between them, doing exactly ONE heap
// allocation regardless of len(parts). Pre-size a strings.Builder with the exact
// total length (sum of part lengths plus separators), then WriteString. The
// test asserts JoinOnce allocates exactly 1 time.
func JoinOnce(parts []string, sep string) string {
	return "" // replace me
}

// TODO(reader): return the p-th percentile (0..100) of a benchmark's latency
// samples in nanoseconds. Sort a COPY (do not mutate xs), take the value at rank
// ceil(p/100 * len) - 1, clamped to a valid index.
func Percentile(xs []int64, p float64) int64 {
	return 0 // replace me
}
