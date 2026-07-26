// Swift has BOTH a typed `Result<Success, Failure>` (a coproduct, with map /
// flatMap chaining) and `throws`/`try` for propagation. Errors are values;
// `try` is the propagation operator, and an unhandled `throws` is a compile
// error at the call site (you must `try`, `try?`, or `try!`).

// region:result:start

public enum CalcError: Error { case tooBig }

/// throws-based: `try` propagates, the happy path stays linear.
public func chain(_ x: Int) throws -> Int {
  let v = x * 2
  if v > 100 { throw CalcError.tooBig }
  return v + 1
}

/// The same, as a chained `Result` with flatMap (map + flatMap chain, like
/// Rust's map/and_then).
public func chainResult(_ x: Int) -> Result<Int, CalcError> {
  Result.success(x)
    .map { $0 * 2 }
    .flatMap { $0 > 100 ? .failure(.tooBig) : .success($0 + 1) }
}

// region:result:end
