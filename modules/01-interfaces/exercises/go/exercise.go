// Package exercises: Module 01 reader tasks. RED until you implement the stubs.
// CI does not run these. Compare with ../../solutions/go.
package exercises

import "math"

type Shape interface{ Area() float64 }

type Circle struct{ R float64 }

func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

type Square struct{ S float64 }

func (s Square) Area() float64 { return s.S * s.S }

// TODO(reader): return the shape with the largest Area. On an empty slice
// return nil. This exercises dynamic dispatch through the itab.
func MaxAreaShape(ss []Shape) Shape {
	return nil // replace me
}

// TODO(reader): count how many elements are (dynamically) Circles, using a type
// switch or a comma-ok type assertion. This exercises the runtime type tag in
// the interface value.
func CountCircles(ss []Shape) int {
	return -1 // replace me
}
