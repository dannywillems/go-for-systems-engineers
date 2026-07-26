package solutions

import "testing"

func TestMaxAreaShape(t *testing.T) {
	ss := []Shape{Circle{R: 1}, Square{S: 10}, Circle{R: 2}}
	if sq, ok := MaxAreaShape(ss).(Square); !ok || sq.S != 10 {
		t.Fatalf("MaxAreaShape wrong: %#v", MaxAreaShape(ss))
	}
	if MaxAreaShape(nil) != nil {
		t.Fatal("empty should be nil")
	}
}

func TestCountCircles(t *testing.T) {
	ss := []Shape{Circle{R: 1}, Square{S: 2}, Circle{R: 3}}
	if got := CountCircles(ss); got != 2 {
		t.Fatalf("CountCircles = %d, want 2", got)
	}
}
