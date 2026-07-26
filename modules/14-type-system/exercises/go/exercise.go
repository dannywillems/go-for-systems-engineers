// Package exercises: Module 14 reader tasks. RED until you implement the stubs.
package exercises

// Distinct newtypes prevent mixing Celsius and Fahrenheit at compile time.
type Celsius float64

type Fahrenheit float64

// TODO(reader): convert Celsius to Fahrenheit: F = C*9/5 + 32.
func CToF(c Celsius) Fahrenheit {
	return 0 // replace me
}

// Stack is a generic LIFO. Methods on a generic TYPE are allowed (they add no
// new type parameter); only a method introducing its OWN type parameter is not.
type Stack[T any] struct {
	items []T
}

// TODO(reader): push v onto the stack.
func (s *Stack[T]) Push(v T) {
	// replace me
}

// TODO(reader): pop the top value; ok is false if the stack is empty.
func (s *Stack[T]) Pop() (v T, ok bool) {
	return v, false // replace me
}

// TODO(reader): return the number of items.
func (s *Stack[T]) Len() int {
	return -1 // replace me
}
