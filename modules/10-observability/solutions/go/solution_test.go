package solutions

import (
	"fmt"
	"strings"
	"testing"
)

func TestJoinOnceCorrect(t *testing.T) {
	parts := []string{"a", "bb", "ccc"}
	got := JoinOnce(parts, ",")
	want := strings.Join(parts, ",")
	if got != want {
		t.Fatalf("JoinOnce = %q, want %q", got, want)
	}
}

func TestJoinOnceSingleAllocation(t *testing.T) {
	parts := make([]string, 64)
	for i := range parts {
		parts[i] = fmt.Sprintf("chunk-%02d", i)
	}
	allocs := testing.AllocsPerRun(200, func() { _ = JoinOnce(parts, ";") })
	if allocs != 1 {
		t.Fatalf("JoinOnce did %.0f allocations, want exactly 1", allocs)
	}
}

func TestPercentile(t *testing.T) {
	xs := []int64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
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
