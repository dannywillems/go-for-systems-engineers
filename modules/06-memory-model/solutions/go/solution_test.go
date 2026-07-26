package solutions

import (
	"sort"
	"testing"
)

func TestIncrement(t *testing.T) {
	if Increment(8, 10_000) != 80_000 {
		t.Fatal("want 80000")
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
	if len(got) != 6 || got[0] != 1 || got[5] != 6 {
		t.Fatalf("Merge = %v", got)
	}
}
