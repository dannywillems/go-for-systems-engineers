# 04 — Comparison: error handling across five languages

**Environment.** Apple M4 Pro, macOS arm64. Go 1.26.5, Rust 1.92.0, OCaml
5.4.0, Swift 6.2.3, Kotlin 2.4.10.

## The shape of "can fail"

| Language | "can fail" type      | Propagation      | Compiler forces handling? | Chaining |
| -------- | -------------------- | ---------------- | ------------------------- | -------- |
| Go       | `(T, error)` product | manual `if err != nil` | no (`errcheck` is external) | no (free functions only) |
| Rust     | `Result<T, E>` sum   | `?`              | yes (`#[must_use]` + `match`) | `.map().and_then()` |
| OCaml    | `('a, 'e) result` sum | `let*` (user-defined bind) | via a warning; enforced under `-warn-error` | `let*` do-notation |
| Swift    | `throws` / `Result`  | `try`            | yes (`try`/`try?`/`try!` required at call site) | `Result.map/flatMap` |
| Kotlin   | unchecked exceptions / `Result` | throw / `runCatching` | no (unchecked) | `Result.map/fold` |

Two axes separate them. **Product vs sum**: only Go models failure as a product,
so only Go has the "value AND error together" trap and the typed-nil trap. The
other four use a coproduct, where "error" excludes "value" by construction.
**Forced handling**: Rust and Swift make ignoring the error channel a compile
error; OCaml does under `-warn-error`; Go and Kotlin do not, so both rely on an
external tool (`errcheck`) or discipline.

## The same computation, five ways

<!-- BEGIN:snippet go-result -->
```go
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
```
<!-- END:snippet go-result -->

<!-- BEGIN:snippet rust-result -->
```rust
/// The SAME computation as the Go demo, but written as a left-to-right chain
/// (`map` then `and_then`), which Go's free-function Result cannot express.
pub fn chain(x: i32) -> Result<i32, String> {
    Ok(x).map(|v| v * 2).and_then(|v| {
        if v > 100 {
            Err("too big".into())
        } else {
            Ok(v + 1)
        }
    })
}

/// `?` propagates the Err early; the happy path stays at the left margin.
pub fn use_question(x: i32) -> Result<i32, String> {
    let doubled = chain(x)?;
    Ok(doubled + 100)
}
```
<!-- END:snippet rust-result -->

<!-- BEGIN:snippet ocaml-result -->
```ocaml
let ( let* ) = Result.bind

(* The same computation as the Go and Rust demos, in let*-notation. *)
let chain x =
  let* v = Ok (x * 2) in
  if v > 100 then Error "too big" else Ok (v + 1)

(* let* threads the error through automatically; the success path is linear. *)
let use_bind x =
  let* doubled = chain x in
  Ok (doubled + 100)
```
<!-- END:snippet ocaml-result -->

<!-- BEGIN:snippet swift-result -->
```swift
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
```
<!-- END:snippet swift-result -->

<!-- BEGIN:snippet kotlin-result -->
```kotlin
class CalcException(
    msg: String,
) : Exception(msg)

// Exception-based: nothing at the call site forces handling.
fun chain(x: Int): Int {
    val v = x * 2
    if (v > 100) throw CalcException("too big")
    return v + 1
}

// Value-based: runCatching turns exceptions into a Result you can map/fold.
fun chainResult(x: Int): Result<Int> = runCatching { chain(x) }
```
<!-- END:snippet kotlin-result -->

## Outputs (all compute the same 7)

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
Ok(3).map(*2).and_then(+1) = Ok(7)  (method chaining)
chain(60) = Err("too big")
use_question(3) with ? = Ok(107)
```
<!-- END:output rust-demo -->

<!-- BEGIN:output ocaml-demo -->
```text
$ dune exec bin/main.exe
chain 3 = Ok 7  (let* notation)
chain 60 = Error too big
use_bind 3 = Ok 107
```
<!-- END:output ocaml-demo -->

<!-- BEGIN:output swift-demo -->
```text
try? chain(3) = 7
chainResult(3) = success(7)
chainResult(60) = failure(Errs.CalcError.tooBig)
```
<!-- END:output swift-demo -->

<!-- BEGIN:output kotlin-demo -->
```text
chain(3) = 7
chainResult(3) = Success(7)
chainResult(60).isFailure = true
```
<!-- END:output kotlin-demo -->

## Reading

Go's `(T, error)` product is the deliberate consequence of not having sum types
(Module 02) plus not having generic methods (Module 03): with no coproduct it
cannot express `Result`, and even the hand-rolled encoding cannot chain. What Go
offers instead is uniform, explicit, in-line error handling with `%w` wrapping,
`errors.Is`/`As` for matching a sentinel or extracting a typed error, and
`errors.Join` for a DAG of errors — all shown returning `true` in `go-demo`
above. The trade is ergonomics and a compiler guarantee for a smaller type
system and a single, visible control-flow idiom; the costs are the two silent
traps (product confusion, typed nil) that this module makes reproducible.
