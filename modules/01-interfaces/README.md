# 01 — Interfaces, dynamic dispatch, and cost

**Thesis.** A Go interface value is an existential package `∃X. (X, vtable_X)`
represented as two words `(itab, data)`. Satisfaction is *structural*: the
coercion `X <: Shape` needs no declaration and is witnessed by an `itab` the
linker synthesizes where the coercion happens. This buys retroactive,
decentralized polymorphism at the cost of a two-word value, an indirect call,
and lost inlining — measured below.

## Structural satisfaction

No `implements` clause exists anywhere; `Circle` is a `Shape` by structure:

<!-- BEGIN:snippet go-iface -->
```go
// Shape is satisfied structurally: any type with an Area() float64 method is a
// Shape, with no "implements" clause anywhere. This is the width-subtyping
// coercion X <: Shape done by structure, not by nominal declaration.
type Shape interface {
	Area() float64
}

type Circle struct{ R float64 }

func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

type Square struct{ S float64 }

func (s Square) Area() float64 { return s.S * s.S }
```
<!-- END:snippet go-iface -->

The interface value is two words; a bare pointer is one. The last two lines are
the typed-nil trap (an interface holding a nil `*T` is itself non-nil, because
its `itab` is set) — pursued in Module 04:

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
a Circle used as a Shape has area 3.1416
sizeof(interface value) = 16 bytes
sizeof(*Circle)         = 8 bytes
nil interface == nil:            true
interface holding (*T)(nil) == nil: false
```
<!-- END:output go-demo -->

## Dispatch cost

Three paths over 1024 shapes: `SumDirect` (concrete slice, statically bound and
inlined), `SumInterface` (existential slice, dispatched through the `itab`), and
`SumDevirt` (type-asserted back to concrete). The compiler's own report shows
`SumDirect` and the asserted branch of `SumDevirt` inlining `Circle.Area`, and
the demo call site being *devirtualized*; `SumInterface` is absent, i.e. not
inlined:

<!-- BEGIN:output go-devirt -->
```text
./shapes.go:42:18: inlining call to Circle.Area
./shapes.go:65:19: inlining call to Circle.Area
cmd/demo/main.go:16:63: devirtualizing s.Area to shapes.Circle
cmd/demo/main.go:16:63: inlining call to shapes.Circle.Area
```
<!-- END:output go-devirt -->

Measured (`benchstat`, `-count=12`):

<!-- BEGIN:file go-bench -->
```text
goos: darwin
goarch: arm64
pkg: github.com/dannywillems/go-for-systems-engineers/modules/01-interfaces/go
cpu: Apple M4 Pro
             │     go      │
             │   sec/op    │
Direct-14      747.5n ± 0%
Interface-14   2.605µ ± 0%
Devirt-14      1.340µ ± 0%
geomean        1.377µ

             │      go      │
             │     B/op     │
Direct-14      0.000 ± 0%
Interface-14   0.000 ± 0%
Devirt-14      0.000 ± 0%
geomean                   ¹
¹ summaries must be >0 to compute geomean

             │      go      │
             │  allocs/op   │
Direct-14      0.000 ± 0%
Interface-14   0.000 ± 0%
Devirt-14      0.000 ± 0%
geomean                   ¹
¹ summaries must be >0 to compute geomean
```
<!-- END:file go-bench -->

The interface path costs several times the direct path here. The dominant cost
is **lost inlining**, not the indirect jump: `Area` is a single multiply, so
inlining it away is most of the win, and a modern branch predictor makes the
indirect call itself cheap. With a heavier method body the *relative* overhead
shrinks. `SumDevirt` recovers roughly half by restoring inlining on the common
type. All three allocate nothing.

## Why Go has no orphan rule

Rust's coherence forbids implementing a foreign trait for a foreign type. This
crate does not compile:

<!-- BEGIN:output rust-orphan-error -->
```text
   Compiling reject-orphan v0.1.0 (./reject-orphan)
error[E0117]: only traits defined in the current crate can be implemented for types defined outside of the crate
   = note: for more information see https://doc.rust-lang.org/reference/items/implementations.html#orphan-rules
```
<!-- END:output rust-orphan-error -->

Go has no such rule: a method is declared syntactically with its receiver type,
and interface satisfaction is computed structurally at each use site, so there
is no global "one impl per type" invariant to protect. The trade-off is real —
Rust guarantees trait dispatch is unambiguous program-wide; Go cannot make the
same guarantee for methods promoted through struct embedding. See
[`COMPARISON.md`](COMPARISON.md) for the five-language axis (nominal vs
structural, static vs dynamic, existential representation) and
[`exercises/`](exercises).

## References

Official sources first, grouped by language.

### Go

- Go spec, Interface types: https://go.dev/ref/spec#Interface_types
- Russ Cox, "Go Data Structures: Interfaces" (the itab layout): https://research.swtch.com/interfaces
- The Go Blog, "The Laws of Reflection": https://go.dev/blog/laws-of-reflection

### Rust

- Rust Reference, trait objects (`dyn`): https://doc.rust-lang.org/reference/types/trait-object.html
- Rust Reference, orphan / coherence rules: https://doc.rust-lang.org/reference/items/implementations.html#orphan-rules

### OCaml

- OCaml Manual, Objects: https://ocaml.org/manual/5.4/objectexamples.html
- OCaml Manual, First-class modules: https://ocaml.org/manual/5.4/firstclassmodules.html

### Swift

- Swift, Protocols (existential `any` vs opaque `some`): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/protocols/

### Kotlin

- Kotlin, Interfaces: https://kotlinlang.org/docs/interfaces.html
