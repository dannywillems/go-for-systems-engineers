package gen

// region:shape:start

// Identity is generic over any type. The instantiations forced by Shapes below
// expose GCShape stenciling in the compiled assembly (see the README): the
// compiler emits ONE stencil per GC shape, not per type. For a value type the
// shape is essentially its UNDERLYING type, so int and float64 differ but
// Celsius (underlying float64) shares float64's shape; every pointer type
// collapses to the one shared shape go.shape.*uint8 (which relies on a runtime
// dictionary).
func Identity[T any](x T) T { return x }

type Cat struct{ Legs int }

type Dog struct{ Tail bool }

// Celsius has underlying type float64, so it carries no shape of its own: it
// shares go.shape.float64 with float64. Instantiating it adds no shape line.
type Celsius float64

// Shapes forces the compiler to emit the instantiations the README inspects.
func Shapes() {
	_ = Identity[int](0)
	_ = Identity[int64](0)
	_ = Identity[float64](0)
	_ = Identity[Celsius](0) // collapses into go.shape.float64
	_ = Identity[string]("")
	_ = Identity[*Cat](nil)
	_ = Identity[*Dog](nil)
}

// region:shape:end
