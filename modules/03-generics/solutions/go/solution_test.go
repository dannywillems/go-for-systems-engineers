package solutions

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(x int) string { return strconv.Itoa(x * 2) })
	if len(got) != 3 || got[0] != "2" || got[2] != "6" {
		t.Fatalf("Map = %v", got)
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5, 6}, func(x int) bool { return x%2 == 0 })
	if len(got) != 3 || got[0] != 2 || got[2] != 6 {
		t.Fatalf("Filter = %v", got)
	}
}
