package errs

// Result is a hand-rolled Result monad. It is a LAWFUL monad (the property
// tests check left identity, right identity, and associativity), and yet it is
// unusable in idiomatic Go, for a precise reason: Map and AndThen change the
// type parameter (T -> U), Go methods may not introduce a type parameter, so
// they must be FREE FUNCTIONS. That destroys method chaining. See
// reject-genericmethod for the compile error a generic method would be.

// region:result:start

type Result[T any] struct {
	val T
	err error
}

func Ok[T any](v T) Result[T]      { return Result[T]{val: v} }
func Err[T any](e error) Result[T] { return Result[T]{err: e} }

func (r Result[T]) Unwrap() (T, error) { return r.val, r.err }

// Map and AndThen are free functions, NOT methods, because they introduce a new
// type parameter U. This is what forces `AndThen(Map(r, f), g)` instead of the
// `r.map(f).and_then(g)` a Rust or OCaml engineer expects.
func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	return Ok(f(r.val))
}

func AndThen[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	return f(r.val)
}

// region:result:end

// equalIntResult compares two Result[int] by observable content, for the laws.
func equalIntResult(a, b Result[int]) bool {
	return a.err == b.err && a.val == b.val
}
