// Swift enums are coproducts and `switch` must be exhaustive; a missing case is
// a compile error. `indirect` allows the recursive payload. See reject-swift.

// region:enum:start

public indirect enum Expr {
  case lit(Int)
  case add(Expr, Expr)
  case mul(Expr, Expr)
  case neg(Expr)
}

public func eval(_ e: Expr) -> Int {
  switch e {
  case .lit(let v): return v
  case .add(let l, let r): return eval(l) + eval(r)
  case .mul(let l, let r): return eval(l) * eval(r)
  case .neg(let x): return -eval(x)
  }
}

// region:enum:end

// Optional<T> is itself an enum (`.none | .some(T)`); Swift needs no nullable
// hack, and pointer-like payloads are niche-packed to the pointer's own size.
