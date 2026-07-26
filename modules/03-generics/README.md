# 03 — Generics: type sets, GCShape, dictionaries

**Thesis.** Go generics are a **hybrid** between Rust/C++ full monomorphization
and Java/OCaml uniform boxing. The compiler emits one instantiation per *GC
shape*: each distinct value type gets its own stencil (as fast as hand-written
code), while **all pointer types collapse to a single shape**
(`go.shape.*uint8`) that shares one stencil plus a runtime **dictionary**. Two
things are inexpressible: generic methods (Module 04) and higher-kinded types.

## Type sets and stenciling

A constraint is a *type set*; `Sum` is generic over it:

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

For a value type, `Sum[int]` is stenciled to a dedicated instantiation, so it
matches the hand-written concrete loop; the interface (boxed) path is far
slower. Measured (`-count=12`, `benchstat`):

<!-- BEGIN:file go-bench -->
```text
goos: darwin
goarch: arm64
pkg: github.com/dannywillems/go-for-systems-engineers/modules/03-generics/go
cpu: Apple M4 Pro
             │     go      │
             │   sec/op    │
Concrete-14    240.3n ± 1%
Generic-14     239.3n ± 1%
Interface-14   1.053µ ± 0%
geomean        392.7n

             │      go      │
             │     B/op     │
Concrete-14    0.000 ± 0%
Generic-14     0.000 ± 0%
Interface-14   0.000 ± 0%
geomean                   ¹
¹ summaries must be >0 to compute geomean

             │      go      │
             │  allocs/op   │
Concrete-14    0.000 ± 0%
Generic-14     0.000 ± 0%
Interface-14   0.000 ± 0%
geomean                   ¹
¹ summaries must be >0 to compute geomean
```
<!-- END:file go-bench -->

`Generic` tracks `Concrete` to within noise; `Interface` is several times
slower. Generics buy interface-like abstraction at concrete-like speed — for
value types.

## GCShape: the pointer collapse

The mechanism is visible in the compiled assembly. `Identity` is instantiated at
six types; the emitted shapes are:

<!-- BEGIN:snippet go-shape -->
```go
// Identity is generic over any type. The instantiations forced by Shapes below
// expose GCShape stenciling in the compiled assembly (see the README): each
// value type gets its OWN shape, while every pointer type collapses to the one
// shared shape go.shape.*uint8 (which then relies on a runtime dictionary).
func Identity[T any](x T) T { return x }

type Cat struct{ Legs int }

type Dog struct{ Tail bool }

// Shapes forces the compiler to emit the instantiations the README inspects.
func Shapes() {
	_ = Identity[int](0)
	_ = Identity[int64](0)
	_ = Identity[float64](0)
	_ = Identity[string]("")
	_ = Identity[*Cat](nil)
	_ = Identity[*Dog](nil)
}
```
<!-- END:snippet go-shape -->

<!-- BEGIN:output go-shapes -->
```text
gen.Identity[go.shape.*uint8]
gen.Identity[go.shape.float64]
gen.Identity[go.shape.int]
gen.Identity[go.shape.int64]
gen.Identity[go.shape.string]
```
<!-- END:output go-shapes -->

Six instantiations, **five** shapes: `*Cat` and `*Dog` both became
`go.shape.*uint8`. Distinct value types keep distinct stencils; all pointers
share one, and the concrete type is recovered at runtime from a per-call
dictionary. This is why a generic over pointer types is not fully monomorphized:
it pays a dictionary indirection Rust does not.

## Higher-kinded types are inexpressible

A type parameter cannot itself be a type constructor, so `Functor`/`Monad`
abstractions do not typecheck:

<!-- BEGIN:output go-hkt-reject -->
```text
./badhkt.go:10:34: invalid operation: F[A] (F is not a generic type)
```
<!-- END:output go-hkt-reject -->

## The monomorphization counter-cost, measured honestly

Rust monomorphizes: each `sum::<T>` is its own machine code. The textbook
counter-cost is bigger binaries and slower compiles. But measure it and the
naive expectation is wrong at this scale:

<!-- BEGIN:output go-binsize -->
```text
Go demo binary:            2509714 bytes
```
<!-- END:output go-binsize -->

<!-- BEGIN:output rust-binsize -->
```text
Rust demo binary (release): 373408 bytes
```
<!-- END:output rust-binsize -->

The Go binary is several times *larger* than the Rust one, because Go statically
links a runtime, GC, and scheduler (a multi-megabyte floor) while a small Rust
release binary with LTO does not. Monomorphization's cost is real, but it shows
up as *growth* with the number of instantiations and in compile time — not as
baseline size for a small program. Stating "monomorphization makes binaries big"
without measuring would have been exactly backwards here.

See [`COMPARISON.md`](COMPARISON.md) for the five strategies (Go GCShape+dict,
Rust monomorphization, OCaml functors + boxed parametric polymorphism, Swift
witness tables, Kotlin/JVM erasure + `reified`), and [`exercises/`](exercises).

## References

Official sources first.

- Go spec, Type parameters and constraints: https://go.dev/ref/spec#Type_parameter_declarations
- The Go Blog, "An Introduction To Generics": https://go.dev/blog/intro-generics
- Go generics implementation (GCShape, dictionaries) design: https://go.googlesource.com/proposal/+/refs/heads/master/design/generics-implementation-dictionaries-go1.18.md
- PlanetScale, "Generics can make your Go code slower" (dictionaries measured): https://planetscale.com/blog/generics-can-make-your-go-code-slower
- The Rust Reference, generic parameters / monomorphization: https://doc.rust-lang.org/reference/items/generics.html
- OCaml Manual, Functors: https://ocaml.org/manual/5.4/moduleexamples.html#s:functors
- Swift generics (The Swift Programming Language): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/generics/
- Kotlin, Generics and type erasure: https://kotlinlang.org/docs/generics.html
