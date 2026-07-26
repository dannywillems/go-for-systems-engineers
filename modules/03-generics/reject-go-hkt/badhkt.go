// Package badhkt DOES NOT COMPILE. Go has no higher-kinded types: a type
// parameter cannot itself be a type constructor, so you cannot abstract over
// `F[_]` the way `Functor`/`Monad` require. `Fmap` below tries to apply the type
// parameter F to a type argument A, which Go rejects ("F is not a generic
// type"). This is why a generic `map` over an arbitrary container, or a Monad
// abstraction, is inexpressible in Go — the point where Rust (GATs) and OCaml
// (functors / higher-kinded via modules) go further.
package badhkt

func Fmap[F any, A any, B any](x F[A], f func(A) B) F[B] {
	var z F[B]
	_ = x
	_ = f
	return z
}
