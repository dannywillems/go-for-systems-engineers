package exercises

import "testing"

func TestParallelMap(t *testing.T) {
	xs := make([]int, 1000)
	for i := range xs {
		xs[i] = i
	}
	got := ParallelMap(xs, 8, func(x int) int { return x * x })
	if len(got) != len(xs) {
		t.Fatalf("len = %d, want %d", len(got), len(xs))
	}
	for i := range xs {
		if got[i] != i*i {
			t.Fatalf("got[%d] = %d, want %d (order preserved?)", i, got[i], i*i)
		}
	}
}

func TestPercentile(t *testing.T) {
	xs := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	if got := Percentile(xs, 50); got != 5 {
		t.Fatalf("p50 = %d, want 5", got)
	}
	if got := Percentile(xs, 100); got != 10 {
		t.Fatalf("p100 = %d, want 10", got)
	}
	if xs[0] != 10 {
		t.Fatal("Percentile mutated its input")
	}
}
