package exercises

import (
	"sort"
	"testing"
)

func TestIncrement(t *testing.T) {
	if got := Increment(8, 10_000); got != 80_000 {
		t.Fatalf("Increment(8, 10000) = %d, want 80000", got)
	}
}

func TestMerge(t *testing.T) {
	mk := func(vs ...int) <-chan int {
		ch := make(chan int, len(vs))
		for _, v := range vs {
			ch <- v
		}
		close(ch)
		return ch
	}
	var got []int
	for v := range Merge(mk(1, 2), mk(3), mk(4, 5, 6)) {
		got = append(got, v)
	}
	sort.Ints(got)
	want := []int{1, 2, 3, 4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("Merge yielded %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Merge yielded %v, want %v", got, want)
		}
	}
}
