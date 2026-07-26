// Package gen demonstrates Go generics: type sets, GCShape stenciling, and the
// runtime dictionary — a hybrid between Rust/C++ full monomorphization and
// Java/OCaml uniform boxing.
//
// Go compiles ONE instantiation per "GC shape". Distinct value types (int,
// float64, string) each get their own stencil, so a generic over them is as
// fast as hand-written concrete code. All pointer types collapse to a SINGLE
// shape (go.shape.*uint8) sharing one stencil plus a runtime DICTIONARY that
// carries the concrete type's descriptors. See the objdump capture in the
// README for both facts. Not expressible: generic methods (Module 04) and
// higher-kinded types (reject-go-hkt).
package gen

// region:sum:start

// Number is a type set (a constraint listing the permitted underlying types).
// The ~ means "any type whose underlying type is this".
type Number interface {
	~int | ~int64 | ~float64
}

// Sum is generic over any Number. For value type args it is STENCILED to a
// dedicated instantiation, so it matches the concrete loop below.
func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

// region:sum:end

// SumIntConcrete is the hand-written monomorphic baseline.
func SumIntConcrete(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}

// Adder is an interface (dynamic dispatch), the pre-generics way to write "sum
// of anything addable". Every element is boxed and every Add is an itab call.
type Adder interface {
	AddTo(acc int) int
}

// Int is a value type implementing Adder.
type Int int

func (i Int) AddTo(acc int) int { return acc + int(i) }

// SumInterface sums via dynamic dispatch, for the cost comparison.
func SumInterface(xs []Adder) int {
	total := 0
	for _, x := range xs {
		total = x.AddTo(total)
	}
	return total
}
