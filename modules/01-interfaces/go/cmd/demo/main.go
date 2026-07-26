// Command demo prints the deterministic, architecture-independent facts about
// interface representation. stdout is injected into the module README.
package main

import (
	"fmt"
	"unsafe"

	shapes "github.com/dannywillems/go-for-systems-engineers/modules/01-interfaces/go"
)

func main() {
	// Structural satisfaction: no "implements" clause exists, yet a Circle is a
	// Shape. The coercion is by structure.
	var s shapes.Shape = shapes.Circle{R: 1}
	fmt.Printf("a Circle used as a Shape has area %.4f\n", s.Area())

	// An interface value is two words (itab pointer, data pointer). A bare
	// pointer is one word. This is the runtime cost of the existential package.
	var iface shapes.Shape
	var ptr *shapes.Circle
	fmt.Printf("sizeof(interface value) = %d bytes\n", unsafe.Sizeof(iface))
	fmt.Printf("sizeof(*Circle)         = %d bytes\n", unsafe.Sizeof(ptr))

	// A nil interface has a nil itab AND nil data; an interface holding a nil
	// pointer has a non-nil itab. This is the "typed nil" trap (Module 04). The
	// value is built behind a function boundary so the concrete type is not
	// visible at the comparison (both to mirror real code and to keep the
	// behavior a runtime fact, not a compile-time-known constant).
	// The value is stored in a slice so the comparison is a genuine runtime
	// check, not a compile-time-known constant.
	box := typedNilBox()
	fmt.Printf("nil interface == nil:            %v\n", iface == nil)
	fmt.Printf("interface holding (*T)(nil) == nil: %v\n", box[0] == nil)
}

type nilShape struct{}

func (*nilShape) Area() float64 { return 0 }

// typedNilBox returns a one-element slice whose element is a Shape wrapping a
// nil *nilShape. That element is non-nil: it carries a non-nil itab and a nil
// data word, the essence of the typed-nil trap.
func typedNilBox() []shapes.Shape {
	var p *nilShape
	return []shapes.Shape{p}
}
