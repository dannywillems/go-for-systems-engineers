package exercises

import (
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(x int) string { return strconv.Itoa(x * 2) })
	want := []string{"2", "4", "6"}
	if len(got) != len(want) {
		t.Fatalf("Map = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Map = %v, want %v", got, want)
		}
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5, 6}, func(x int) bool { return x%2 == 0 })
	want := []int{2, 4, 6}
	if len(got) != len(want) {
		t.Fatalf("Filter = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Filter = %v, want %v", got, want)
		}
	}
}
