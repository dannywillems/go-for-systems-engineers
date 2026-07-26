package gen

// region:shape:start

// Identity is generic over any type. The instantiations forced by Shapes below
// expose GCShape stenciling in the compiled assembly (see the README): each
// value type gets its OWN shape, while every pointer type collapses to the one
// shared shape go.shape.*uint8 (which then relies on a runtime dictionary).
func Identity[T any](x T) T { return x }

type Cat struct{ Legs int }

type Dog struct{ Tail bool }

// Shapes forces the compiler to emit the instantiations the README inspects.
func Shapes() {
	_ = Identity[int](0)
	_ = Identity[int64](0)
	_ = Identity[float64](0)
	_ = Identity[string]("")
	_ = Identity[*Cat](nil)
	_ = Identity[*Dog](nil)
}

// region:shape:end
