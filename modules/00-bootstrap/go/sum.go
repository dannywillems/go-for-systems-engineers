// Package bootstrap holds the trivial computation the three languages share so
// that Module 00 can prove, byte-for-byte, that all three toolchains build, run,
// and agree, and that the capture harness injects real program output.
package bootstrap

import "unsafe"

// region:sum:start
// Sum returns 1 + 2 + ... + n. The value is identical on every 64-bit target
// and in every language, which makes it a clean cross-toolchain fixture.
func Sum(n int) int {
	total := 0
	for i := 1; i <= n; i++ {
		total += i
	}
	return total
}

// WordSizeBytes is sizeof(int) on this target: 8 on any 64-bit platform,
// so it is stable across the author's arm64 and a CI amd64 runner.
func WordSizeBytes() int {
	var x int
	return int(unsafe.Sizeof(x))
}

// region:sum:end
