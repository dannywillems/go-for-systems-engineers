// Package exercises: Module 03 reader tasks. RED until you implement the stubs.
// Note these MUST be free functions: Map changes the type parameter (T -> U),
// and Go methods cannot have type parameters (Module 04).
package exercises

// TODO(reader): implement a generic map. Return a new slice with f applied to
// each element. Preallocate the result to len(xs).
func Map[T, U any](xs []T, f func(T) U) []U {
	return nil // replace me
}

// TODO(reader): implement a generic filter. Return a new slice of the elements
// for which pred returns true, preserving order.
func Filter[T any](xs []T, pred func(T) bool) []T {
	return nil // replace me
}
