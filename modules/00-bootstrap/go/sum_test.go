package bootstrap

import "testing"

func TestSum(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 1},
		{10, 55},
		{1_000_000, 500_000_500_000},
	}
	for _, tt := range tests {
		if got := Sum(tt.n); got != tt.want {
			t.Errorf("Sum(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestWordSize(t *testing.T) {
	// Every target this repo is exercised on is 64-bit.
	if got := WordSizeBytes(); got != 8 {
		t.Errorf("WordSizeBytes() = %d, want 8 (non-64-bit target?)", got)
	}
}

// BenchmarkSum exists so Module 00 can prove the benchmark + benchstat harness
// end to end. sink defeats dead-code elimination of the result.
var sink int

func BenchmarkSum(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = Sum(1_000_000)
	}
	_ = sink
}
