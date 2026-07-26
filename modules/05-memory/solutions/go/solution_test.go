package solutions

import "testing"

func TestSafeHead(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	h := SafeHead(s, 2)
	_ = append(h, 99) //nolint:gocritic // testing that append cannot reach s
	if s[2] != 3 {
		t.Fatalf("aliased: s[2] = %d", s[2])
	}
}

func TestSumNoAlloc(t *testing.T) {
	xs := []int{1, 2, 3, 4, 5}
	if SumNoAlloc(xs) != 15 {
		t.Fatal("wrong sum")
	}
	if a := testing.AllocsPerRun(100, func() { _ = SumNoAlloc(xs) }); a != 0 {
		t.Fatalf("allocated %v", a)
	}
}
