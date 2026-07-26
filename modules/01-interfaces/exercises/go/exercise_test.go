package exercises

import "testing"

func TestMaxAreaShape(t *testing.T) {
	ss := []Shape{Circle{R: 1}, Square{S: 10}, Circle{R: 2}}
	got := MaxAreaShape(ss)
	if sq, ok := got.(Square); !ok || sq.S != 10 {
		t.Fatalf("MaxAreaShape = %#v, want Square{S:10}", got)
	}
	if MaxAreaShape(nil) != nil {
		t.Fatal("MaxAreaShape(nil) should be nil")
	}
}

func TestCountCircles(t *testing.T) {
	ss := []Shape{Circle{R: 1}, Square{S: 2}, Circle{R: 3}, Square{S: 4}}
	if got := CountCircles(ss); got != 2 {
		t.Fatalf("CountCircles = %d, want 2", got)
	}
}
