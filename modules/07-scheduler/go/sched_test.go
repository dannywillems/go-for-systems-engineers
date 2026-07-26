package sched

import (
	"testing"
	"time"
)

func TestPercentile(t *testing.T) {
	ds := []time.Duration{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if got := Percentile(ds, 50); got != 5 {
		t.Errorf("p50 = %v, want 5", got)
	}
	if got := Percentile(ds, 100); got != 10 {
		t.Errorf("p100 = %v, want 10", got)
	}
	if got := Percentile(ds, 10); got != 1 {
		t.Errorf("p10 = %v, want 1", got)
	}
}

func TestPoolLatenciesLength(t *testing.T) {
	lat := PoolLatencies(4, 1000)
	if len(lat) != 1000 {
		t.Fatalf("got %d latencies, want 1000", len(lat))
	}
	for _, d := range lat {
		if d < 0 {
			t.Fatal("negative latency")
		}
	}
}
