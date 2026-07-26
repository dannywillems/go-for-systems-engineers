// Package solutions is the corrigé for Module 11. Run via `make solutions M=11`.
package solutions

// Dedup removes consecutive duplicates, preserving order, without mutating xs.
func Dedup(xs []int) []int {
	out := make([]int, 0, len(xs))
	for i, x := range xs {
		if i == 0 || x != xs[i-1] {
			out = append(out, x)
		}
	}
	return out
}
