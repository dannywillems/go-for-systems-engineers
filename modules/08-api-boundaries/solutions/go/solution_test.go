package solutions

import "testing"

func TestNewRatioReduces(t *testing.T) {
	r, err := NewRatio(2, 4)
	if err != nil {
		t.Fatalf("NewRatio(2,4): %v", err)
	}
	if r.Num() != 1 || r.Den() != 2 {
		t.Fatalf("NewRatio(2,4) = %d/%d, want 1/2", r.Num(), r.Den())
	}
}

func TestNewRatioSign(t *testing.T) {
	r, err := NewRatio(3, -6)
	if err != nil {
		t.Fatalf("NewRatio(3,-6): %v", err)
	}
	if r.Num() != -1 || r.Den() != 2 {
		t.Fatalf("NewRatio(3,-6) = %d/%d, want -1/2 (den > 0)", r.Num(), r.Den())
	}
}

func TestNewRatioZeroDenominator(t *testing.T) {
	if _, err := NewRatio(1, 0); err == nil {
		t.Fatal("NewRatio(1,0) should return an error")
	}
}

func TestMulReduces(t *testing.T) {
	a, _ := NewRatio(1, 2)
	b, _ := NewRatio(2, 3)
	got := a.Mul(b)
	if got.Num() != 1 || got.Den() != 3 {
		t.Fatalf("(1/2)*(2/3) = %d/%d, want 1/3", got.Num(), got.Den())
	}
}
