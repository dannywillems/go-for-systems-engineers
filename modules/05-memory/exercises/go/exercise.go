// Package exercises: Module 05 reader tasks. RED until you implement the stubs.
package exercises

// TODO(reader): return the first k elements of s as a slice that does NOT share
// s's backing array, so appending to the result cannot corrupt s. Use the
// full-slice expression s[:k:k] (or copy). The test appends to your result and
// checks s is unchanged.
func SafeHead(s []int, k int) []int {
	return s[:k] // BUG: this aliases s. Fix it.
}

// TODO(reader): sum xs. Keep it allocation-free (the test asserts zero allocs
// via testing.AllocsPerRun). A plain loop over the slice does not allocate;
// boxing into an interface would. Do not box.
func SumNoAlloc(xs []int) int {
	return 0 // replace me
}
