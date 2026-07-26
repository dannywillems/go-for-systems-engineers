// Package shapes demonstrates Go interfaces as existential packages, structural
// satisfaction, itab-mediated dynamic dispatch, and the cost of each.
//
// Type-theoretically, a value of interface type is an existential:
//
//	Shape  ~=  exists X. (X, { Area: X -> float64 })
//
// The runtime representation is a two-word pair (itab, data): the itab is the
// witness pairing the concrete type X with the interface's method set, built by
// the linker. Structural satisfaction means the coercion X <: Shape needs no
// declaration at the definition site; the itab is synthesized where the
// coercion happens.
package shapes

import "math"

// region:iface:start

// Shape is satisfied structurally: any type with an Area() float64 method is a
// Shape, with no "implements" clause anywhere. This is the width-subtyping
// coercion X <: Shape done by structure, not by nominal declaration.
type Shape interface {
	Area() float64
}

type Circle struct{ R float64 }

func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

type Square struct{ S float64 }

func (s Square) Area() float64 { return s.S * s.S }

// region:iface:end

// SumDirect sums over a concrete slice: the call c.Area() is statically bound
// and inlinable. This is the universally-quantified path resolved at compile
// time (no existential is ever formed).
func SumDirect(cs []Circle) float64 {
	var total float64
	for _, c := range cs {
		total += c.Area()
	}
	return total
}

// SumInterface sums over an existential slice: each s.Area() dispatches through
// the itab of s's concrete type. The compiler cannot inline the callee because
// the concrete type is not known at this site.
func SumInterface(ss []Shape) float64 {
	var total float64
	for _, s := range ss {
		total += s.Area()
	}
	return total
}

// SumDevirt takes an existential but immediately pins the concrete type in a
// local, which the compiler can devirtualize back to a static call. See the
// `-gcflags=-m` capture in the README.
func SumDevirt(ss []Shape) float64 {
	var total float64
	for _, s := range ss {
		if c, ok := s.(Circle); ok {
			total += c.Area()
		} else {
			total += s.Area()
		}
	}
	return total
}
