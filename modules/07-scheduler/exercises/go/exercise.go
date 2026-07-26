// Package exercises: Module 07 reader tasks. RED until you implement the stubs.
// Run the test with -race to prove your pool has no data race.
package exercises

// TODO(reader): apply f to every element of xs using `workers` goroutines, and
// return the results in the SAME order as xs. Each goroutine must write a
// distinct index (no shared-write race). Do not leak goroutines.
func ParallelMap(xs []int, workers int, f func(int) int) []int {
	return nil // replace me
}

// TODO(reader): return the p-th percentile (0..100) of xs. Sort a copy, take the
// value at rank ceil(p/100 * len) - 1 (clamped to a valid index). Do not mutate
// the input.
func Percentile(xs []int, p float64) int {
	return 0 // replace me
}
