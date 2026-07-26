package gen

import "testing"

func TestSumsAgree(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	if Sum(ints) != 15 || SumIntConcrete(ints) != 15 {
		t.Fatal("value sums disagree")
	}
	adders := []Adder{Int(1), Int(2), Int(3), Int(4), Int(5)}
	if SumInterface(adders) != 15 {
		t.Fatal("interface sum wrong")
	}
}

const n = 1024

var (
	benchInts   []int
	benchAdders []Adder
	sink        int
)

func init() {
	benchInts = make([]int, n)
	benchAdders = make([]Adder, n)
	for i := range n {
		benchInts[i] = i
		benchAdders[i] = Int(i)
	}
}

// Concrete: hand-written monomorphic loop (the baseline).
func BenchmarkConcrete(b *testing.B) {
	for b.Loop() {
		sink = SumIntConcrete(benchInts)
	}
}

// Generic: Sum[int] is STENCILED for int, so it should match Concrete.
func BenchmarkGeneric(b *testing.B) {
	for b.Loop() {
		sink = Sum(benchInts)
	}
}

// Interface: boxed elements + itab dispatch, the pre-generics cost.
func BenchmarkInterface(b *testing.B) {
	for b.Loop() {
		sink = SumInterface(benchAdders)
	}
}
