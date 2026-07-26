// Swift generics default to a WITNESS-TABLE strategy strikingly like Go's
// dictionary: one shared body receives a protocol-witness table (and value-
// witness table) at runtime, so a generic works on any conforming type without
// per-type code. Under whole-module / `@inlinable` optimization the compiler can
// SPECIALIZE (monomorphize) hot generics, so Swift spans Go's and Rust's
// strategies depending on optimization.

// region:sum:start

public func sum<T: AdditiveArithmetic>(_ xs: [T]) -> T {
  xs.reduce(.zero, +)
}

// region:sum:end
