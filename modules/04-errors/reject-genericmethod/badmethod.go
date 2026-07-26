// Package badmethod DOES NOT COMPILE. It shows the concrete reason a Result
// monad is unusable in Go: Map changes the type parameter (T -> U), which would
// require a type parameter ON THE METHOD, and Go forbids that.
package badmethod

type Result[T any] struct{ v T }

// Methods cannot have type parameters. This is the language restriction that
// forces Map/AndThen to be free functions and kills method chaining.
func (r Result[T]) Map[U any](f func(T) U) Result[U] {
	return Result[U]{v: f(r.v)}
}
