// Package exercises holds Module 00's reader tasks. The tests in this package
// are RED by design: `go test ./...` here fails until you implement the stubs.
// This is intentional; CI does not run exercises. Fix them, then compare with
// ../../solutions/go.
package exercises

// TODO(reader): return the n-th triangular number using the closed form
// n(n+1)/2, NOT a loop. The test checks it against a summation for many n, which
// is the cheapest possible property test: two independent computations of the
// same value must agree.
func Triangular(n int) int {
	return 0 // replace me
}

// TODO(reader): return the machine word size in BITS on this target (64 on any
// 64-bit platform). Use unsafe.Sizeof of an int and convert bytes to bits.
func WordSizeBits() int {
	return 0 // replace me
}
