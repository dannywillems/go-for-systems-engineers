package exercises

import "testing"

// sumLoop is the independent oracle Triangular must agree with.
func sumLoop(n int) int {
	t := 0
	for i := 1; i <= n; i++ {
		t += i
	}
	return t
}

func TestTriangularMatchesSummation(t *testing.T) {
	for n := 0; n <= 10_000; n++ {
		if got, want := Triangular(n), sumLoop(n); got != want {
			t.Fatalf("Triangular(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestWordSizeBits(t *testing.T) {
	if got := WordSizeBits(); got != 64 {
		t.Fatalf("WordSizeBits() = %d, want 64", got)
	}
}
