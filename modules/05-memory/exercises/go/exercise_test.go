package exercises

import "testing"

func TestSafeHeadDoesNotAlias(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	h := SafeHead(s, 2)
	_ = append(h, 99) //nolint:gocritic // deliberately test that append cannot reach s
	if s[2] != 3 {
		t.Fatalf("SafeHead result aliases s: s[2] = %d, want 3", s[2])
	}
}

func TestSumNoAllocValueAndAllocs(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	if got := SumNoAlloc(xs); got != 15 {
		t.Fatalf("SumNoAlloc = %d, want 15", got)
	}
	allocs := testing.AllocsPerRun(100, func() { _ = SumNoAlloc(xs) })
	if allocs != 0 {
		t.Fatalf("SumNoAlloc allocated %v times, want 0 (did you box into an interface?)", allocs)
	}
}
