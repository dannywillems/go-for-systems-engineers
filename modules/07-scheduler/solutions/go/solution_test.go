package solutions

import "testing"

func TestParallelMap(t *testing.T) {
	xs := make([]int, 1000)
	for i := range xs {
		xs[i] = i
	}
	got := ParallelMap(xs, 8, func(x int) int { return x * x })
	for i := range xs {
		if got[i] != i*i {
			t.Fatalf("got[%d] = %d", i, got[i])
		}
	}
}

func TestPercentile(t *testing.T) {
	xs := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	if Percentile(xs, 50) != 5 || Percentile(xs, 100) != 10 || xs[0] != 10 {
		t.Fatal("Percentile wrong or mutated input")
	}
}
