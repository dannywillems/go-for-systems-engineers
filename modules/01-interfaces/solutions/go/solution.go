// Package solutions is the corrigé for Module 01. Excluded from the default
// build; run via `make solutions`.
package solutions

import "math"

type Shape interface{ Area() float64 }

type Circle struct{ R float64 }

func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

type Square struct{ S float64 }

func (s Square) Area() float64 { return s.S * s.S }

func MaxAreaShape(ss []Shape) Shape {
	var best Shape
	for _, s := range ss {
		if best == nil || s.Area() > best.Area() {
			best = s
		}
	}
	return best
}

func CountCircles(ss []Shape) int {
	n := 0
	for _, s := range ss {
		if _, ok := s.(Circle); ok {
			n++
		}
	}
	return n
}
