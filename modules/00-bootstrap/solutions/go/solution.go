// Package solutions is the corrigé for Module 00's exercises. It is excluded
// from the default build and test; CI runs it only via the `solutions` target
// to verify the answers stay correct.
package solutions

import "unsafe"

// Triangular returns the n-th triangular number in O(1) via the closed form.
func Triangular(n int) int {
	return n * (n + 1) / 2
}

// WordSizeBits returns sizeof(int) in bits: 64 on a 64-bit target.
func WordSizeBits() int {
	var x int
	return int(unsafe.Sizeof(x)) * 8
}
