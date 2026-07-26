package shapes

import (
	"math"
	"testing"
	"unsafe"
)

func TestStructuralSatisfaction(t *testing.T) {
	// Compile-time assertions that both value and pointer satisfy Shape (Area
	// has a value receiver, so both method sets include it).
	var _ Shape = Circle{}
	var _ Shape = (*Circle)(nil)
	var _ Shape = Square{}

	var s Shape = Circle{R: 2}
	if got, want := s.Area(), 4*math.Pi; math.Abs(got-want) > 1e-9 {
		t.Errorf("Area = %v, want %v", got, want)
	}
}

func TestInterfaceValueIsTwoWords(t *testing.T) {
	var s Shape
	if got := unsafe.Sizeof(s); got != 16 {
		t.Errorf("sizeof(interface) = %d, want 16 (two 8-byte words)", got)
	}
}

func TestSumAgree(t *testing.T) {
	cs := []Circle{{1}, {2}, {3}}
	ss := []Shape{Circle{1}, Circle{2}, Circle{3}}
	if math.Abs(SumDirect(cs)-SumInterface(ss)) > 1e-9 {
		t.Error("direct and interface sums disagree")
	}
	if math.Abs(SumDevirt(ss)-SumInterface(ss)) > 1e-9 {
		t.Error("devirt and interface sums disagree")
	}
}

// --- Benchmarks: direct vs interface vs devirtualized dispatch ---------------

const n = 1024

var (
	benchCircles []Circle
	benchShapes  []Shape
	sink         float64
)

func init() {
	benchCircles = make([]Circle, n)
	benchShapes = make([]Shape, n)
	for i := range n {
		c := Circle{R: float64(i%7) + 1}
		benchCircles[i] = c
		benchShapes[i] = c
	}
}

func BenchmarkDirect(b *testing.B) {
	for b.Loop() {
		sink = SumDirect(benchCircles)
	}
}

func BenchmarkInterface(b *testing.B) {
	for b.Loop() {
		sink = SumInterface(benchShapes)
	}
}

func BenchmarkDevirt(b *testing.B) {
	for b.Loop() {
		sink = SumDevirt(benchShapes)
	}
}
