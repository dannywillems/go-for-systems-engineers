# 03 — Comparison: five generics strategies

**Environment.** Apple M4 Pro, macOS arm64. Go 1.26.5, Rust 1.92.0 (edition
2024), OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). Go benchmark in
[`go/bench.txt`](go/bench.txt).

## The strategy axis

| Language | Strategy | Per-type code? | Runtime cost | Reification |
| -------- | -------- | -------------- | ------------ | ----------- |
| Go       | GCShape stenciling + dictionary | per value shape; ONE per all pointers | dictionary for pointer shapes | none (dict has type info) |
| Rust     | full monomorphization | one per instantiation | zero (fully specialized) | at compile time |
| OCaml    | parametric polymorphism (boxed) + functors | none (shared, boxed) | uniform boxing | none |
| Swift    | witness tables (like Go's dict); can specialize | shared by default | witness table; specialized under `-O`/`@inlinable` | at runtime via metadata |
| Kotlin   | JVM type erasure | none (erased) | boxing | only via `inline`+`reified` |

Three clusters. **Rust** (and specialized Swift) pays code size and compile time
to get zero-overhead specialization. **OCaml parametric polymorphism** and
**Kotlin erasure** share one boxed implementation, paying uniform boxing.
**Go and default Swift** sit in the middle: value types are specialized (Go
stencils, Swift can), pointer/reference types share a body plus a runtime
dictionary/witness table. Go's `go.shape.*uint8` collapse and Swift's witness
table are nearly the same idea.

## The five sums

<!-- BEGIN:snippet go-sum -->
```go
// Number is a type set (a constraint listing the permitted underlying types).
// The ~ means "any type whose underlying type is this".
type Number interface {
	~int | ~int64 | ~float64
}

// Sum is generic over any Number. For value type args it is STENCILED to a
// dedicated instantiation, so it matches the concrete loop below.
func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}
```
<!-- END:snippet go-sum -->

<!-- BEGIN:snippet rust-sum -->
```rust
/// Generic over any addable, copyable type with a zero. Monomorphized per T.
pub fn sum<T>(xs: &[T]) -> T
where
    T: Add<Output = T> + Copy + Default,
{
    xs.iter().fold(T::default(), |acc, &x| acc + x)
}
```
<!-- END:snippet rust-sum -->

<!-- BEGIN:snippet ocaml-functor -->
```ocaml
module type ADD = sig
  type t

  val zero : t
  val add : t -> t -> t
end

module Sum (M : ADD) = struct
  let sum xs = List.fold_left M.add M.zero xs
end

module IntSum = Sum (struct
  type t = int

  let zero = 0
  let add = ( + )
end)

module FloatSum = Sum (struct
  type t = float

  let zero = 0.0
  let add = ( +. )
end)
```
<!-- END:snippet ocaml-functor -->

<!-- BEGIN:snippet swift-sum -->
```swift
public func sum<T: AdditiveArithmetic>(_ xs: [T]) -> T {
  xs.reduce(.zero, +)
}
```
<!-- END:snippet swift-sum -->

<!-- BEGIN:snippet kotlin-generic -->
```kotlin
fun <T> sum(
    xs: List<T>,
    zero: T,
    add: (T, T) -> T,
): T = xs.fold(zero, add)

// reified: the type T survives because the function is inlined at each call.
inline fun <reified T> typeName(): String = T::class.simpleName ?: "?"
```
<!-- END:snippet kotlin-generic -->

## Outputs

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
sum::<i64>  = 15
sum::<f64>  = 7
sum::<u32>  = 6
```
<!-- END:output rust-demo -->

<!-- BEGIN:output ocaml-demo -->
```text
$ dune exec bin/main.exe
IntSum.sum [1;2;3;4;5] = 15
FloatSum.sum [1.5;2.5;3.] = 7
```
<!-- END:output ocaml-demo -->

<!-- BEGIN:output swift-demo -->
```text
sum<Int>    = 15
sum<Double> = 7.0
```
<!-- END:output swift-demo -->

<!-- BEGIN:output kotlin-demo -->
```text
sum(ints)   = 15
sum(doubles)= 7.0
reified typeName<Int>() = Int
```
<!-- END:output kotlin-demo -->

## What each cannot express

- **Go**: no generic methods (Module 04) and no higher-kinded types (the
  README's HKT reject), so `Functor`/`Monad`-style abstractions are out.
- **Rust**: expresses HKT-adjacent patterns via GATs, and generic methods, at
  the cost of the monomorphization footprint.
- **OCaml**: functors give higher-kinded abstraction over modules; parametric
  polymorphism is first-class but boxed.
- **Swift**: rich generics with associated types and `some`/`any`, witness-table
  based.
- **Kotlin**: erasure means you cannot inspect `T` at runtime except through
  `reified` (which requires `inline`), and no HKT.

## Reading

Go's GCShape design is a deliberate middle path: most of monomorphization's
speed for value types, without emitting a stencil per pointer type (which keeps
compile time and binary growth in check), at the price of a dictionary
indirection on the pointer path and the two expressiveness gaps above. Whether
that trade is right depends on whether the generic code is hot on value types
(where Go matches Rust) or on pointer types (where the dictionary shows up).
