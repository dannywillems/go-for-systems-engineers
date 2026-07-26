package exercises

import (
	"slices"
	"testing"
	"testing/quick"
)

func TestDedupTable(t *testing.T) {
	cases := []struct {
		in, want []int
	}{
		{[]int{1, 1, 2, 2, 2, 3, 1}, []int{1, 2, 3, 1}},
		{[]int{}, []int{}},
		{[]int{5}, []int{5}},
		{[]int{7, 7, 7}, []int{7}},
	}
	for _, c := range cases {
		if got := Dedup(c.in); !slices.Equal(got, c.want) {
			t.Errorf("Dedup(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

// Property: Dedup is idempotent. quick.Check generates random []int inputs.
func TestDedupIdempotent(t *testing.T) {
	idempotent := func(xs []int) bool {
		once := Dedup(xs)
		return slices.Equal(Dedup(once), once)
	}
	if err := quick.Check(idempotent, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDedupDoesNotMutate(t *testing.T) {
	in := []int{1, 1, 2}
	_ = Dedup(in)
	if !slices.Equal(in, []int{1, 1, 2}) {
		t.Fatalf("Dedup mutated its input: %v", in)
	}
}
