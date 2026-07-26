# 02 — The absence of sum types

**Thesis.** Go has products (structs, `A × B`) and a subtyping-based existential
(interfaces), but no coproduct `A + B` with a compiler-checked eliminator. The
idiom that fakes one — a *sealed interface* (an interface with an unexported
marker method, so variants are closed to one package) — closes the variant set
but leaves `switch x.(type)` unchecked for totality. A missing case compiles and
silently returns a wrong answer. This module builds the missing check as a
`go/analysis` pass, and shows the other four languages' compilers doing it for
free.

## The sealed interface and its hole

`Expr` is sealed by the unexported `exprNode()` marker; the variant set
`{Lit, Add, Mul, Neg}` is closed:

<!-- BEGIN:snippet go-variants -->
```go
type Lit struct{ V int }

type Add struct{ L, R Expr }

type Mul struct{ L, R Expr }

type Neg struct{ X Expr }

func (Lit) exprNode() {}
func (Add) exprNode() {}
func (Mul) exprNode() {}
func (Neg) exprNode() {}
```
<!-- END:snippet go-variants -->

An incomplete switch over `Expr` is legal Go. It does not panic; it silently
returns the zero value or falls through, which is worse:

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
Eval((2*3) + -4) = 2
incomplete.Name(Red)  = "red"
incomplete.Name(Blue) = "?"  (silently wrong: Blue unhandled)
```
<!-- END:output go-demo -->

`incomplete.Name(Blue)` returns `"?"` because the `Blue` case is missing and
nothing forces it to exist.

## The analyzer restores the check

[`go/exhaustive`](go/exhaustive) is a `go/analysis` pass: for any interface
marked `//sumtype:decl`, it collects the package's implementers and reports a
type switch that omits a variant and has no `default`. It is verified by
`analysistest` (in `make test`) and runnable via `go vet -vettool`. Run on the
incomplete package it produces the diagnostic the compiler withholds:

<!-- BEGIN:output go-analyzer -->
```text
examples/incomplete/incomplete.go:19:2: non-exhaustive type switch on Color: missing cases Blue (add them or a default clause)
exit status 3
```
<!-- END:output go-analyzer -->

This is a real, if partial, fix: it depends on the `//sumtype:decl` convention
and only sees switches whose static type is the sealed interface. The point is
that in Go, exhaustiveness is an opt-in *tool*, not a language guarantee.

## The other four check it in the compiler

Rust, OCaml, Swift, and Kotlin all make a non-exhaustive match a compile error.
Each of these does not build:

<!-- BEGIN:output rust-reject -->
```text
error[E0004]: non-exhaustive patterns: `Color::Blue` not covered
```
<!-- END:output rust-reject -->

<!-- BEGIN:output ocaml-reject -->
```text
Error (warning 8 [partial-match]): this pattern-matching is not exhaustive.
```
<!-- END:output ocaml-reject -->

<!-- BEGIN:output swift-reject -->
```text
reject.swift:11:3: error: switch must be exhaustive
```
<!-- END:output swift-reject -->

<!-- BEGIN:output kotlin-reject -->
```text
reject.kt:13:5: error: 'when' expression must be exhaustive. Add the 'Blue' branch or an 'else' branch.
```
<!-- END:output kotlin-reject -->

## Option is the smallest coproduct

`Option<T> = None + Some T` is the one-bit sum, and its treatment is the same
story. Rust's `Option<&T>` is niche-optimized to the size of `&T` (the null
pointer is the `None` tag):

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
eval((2*3) + -4) = 2
size_of::<&u8> = 8
size_of::<Option<&u8>> = 8
size_of::<Option<Box<u8>>> = 8
size_of::<Option<u8>> = 2
```
<!-- END:output rust-demo -->

Go has no `Option`; it uses `(T, bool)` comma-ok, a nil `*T`, or the zero value,
none of which the compiler forces you to check. See
[`COMPARISON.md`](COMPARISON.md) for the five definitions side by side and the
`*T` / `Option<T>` / `T?` representation table, and [`exercises/`](exercises).

## References

Official sources first, grouped by language.

### Go

- Go spec, Type switches: https://go.dev/ref/spec#Type_switches
- `go/analysis` framework: https://pkg.go.dev/golang.org/x/tools/go/analysis
- `analysistest`: https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest
- Prior art: BurntSushi/go-sumtype (the `//sumtype:decl` convention): https://github.com/BurntSushi/go-sumtype

### Rust

- The Rust Book, Enums and pattern matching: https://doc.rust-lang.org/book/ch06-00-enums.html
- Rust Reference, `match` expressions (exhaustiveness): https://doc.rust-lang.org/reference/expressions/match-expr.html

### OCaml

- OCaml Manual, variants and warning 8 (partial-match): https://ocaml.org/manual/5.4/coreexamples.html

### Swift

- Swift, Enumerations: https://docs.swift.org/swift-book/documentation/the-swift-programming-language/enumerations/

### Kotlin

- Kotlin, Sealed classes and `when`: https://kotlinlang.org/docs/sealed-classes.html
