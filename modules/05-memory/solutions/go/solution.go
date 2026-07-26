// Package solutions is the corrigé for Module 05. Run via `make solutions`.
package solutions

// SafeHead caps capacity with the full-slice expression, so append reallocates
// rather than writing into s's backing array.
func SafeHead(s []int, k int) []int {
	return s[:k:k]
}

// SumNoAlloc is a plain loop: no boxing, no allocation.
func SumNoAlloc(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}
