# 08 — API design & package boundaries

**Thesis.** All five languages enforce encapsulation at COMPILE time — an access
violation is a compiler error, not a lint or a convention — but they differ in
two axes: the UNIT of encapsulation (Go: the package, visibility by identifier
case; Rust: the module tree with graduated `pub` levels; OCaml: the `.mli`
signature; Swift: file/module with a five-level ladder; Kotlin: the module with
four levels) and how much they can HIDE. The falsifiable spine is one reject per
language: a program that reaches across a boundary and whose exact compiler error
is captured here, proving the boundary is real.

## Contents

- [Go: the package, visibility by case](#go-the-package-visibility-by-case)
- [The five boundaries, each rejected by its compiler](#the-five-boundaries-each-rejected-by-its-compiler)
- [Hiding the representation: opaque struct vs abstract type](#hiding-the-representation-opaque-struct-vs-abstract-type)
- [Exercises](#exercises)
- [References](#references)

## Go: the package, visibility by case

Go's encapsulation unit is the package, and the rule is mechanical: a
Capitalized identifier is exported, a lowercase one is package-private. An
exported struct whose fields are all unexported is Go's abstract type — other
packages hold it and call its methods but cannot read, write, or construct its
representation:

<!-- BEGIN:snippet go-opaque -->
```go
// Account is opaque: every field is unexported, so no other package can read or
// write balance, nor build an Account with a struct literal. The only way to get
// one is Open, and the only way to change it is the methods.
type Account struct {
	balance int64
}

// Open is the sole constructor: it enforces the invariant (balance >= 0) that
// the hidden field guarantees from then on.
func Open(initial int64) (*Account, error) {
	if initial < 0 {
		return nil, ErrOverdraft
	}
	return &Account{balance: initial}, nil
}

// Deposit and Withdraw are the only mutators; Balance is the only reader.
func (a *Account) Deposit(amount int64) { a.balance += amount }

func (a *Account) Withdraw(amount int64) error {
	if amount > a.balance {
		return ErrOverdraft
	}
	a.balance -= amount
	return nil
}

func (a *Account) Balance() int64 { return a.balance }
```
<!-- END:snippet go-opaque -->

Go has a second, coarser boundary: a directory named `internal/` may be imported
only by code rooted at its parent. Reaching for it from anywhere else is
rejected by the compiler, which is the reject captured below.

## The five boundaries, each rejected by its compiler

Each `reject-*` directory holds a program that violates a boundary and does not
compile. The exact errors are injected here, so the claim "the boundary is
compiler-enforced" is falsifiable:

**Go** — `internal/` is importable only within its parent subtree:

<!-- BEGIN:output go-reject -->
```text
	other/other.go:8:8: use of internal package rejectinternal/lib/internal/store not allowed
```
<!-- END:output go-reject -->

**Rust** — a struct field is private by default; `pub` on the struct does not
expose it:

<!-- BEGIN:output rust-reject -->
```text
error[E0616]: field `balance` of struct `Account` is private
```
<!-- END:output rust-reject -->

**OCaml** — the `.mli` makes the type abstract, so its representation is invisible
outside the module:

<!-- BEGIN:output ocaml-reject -->
```text
Error: Unbound record field Store.n
```
<!-- END:output ocaml-reject -->

**Swift** — `private` confines a member to its enclosing declaration:

<!-- BEGIN:output swift-reject -->
```text
reject.swift:11:9: error: 'balance' is inaccessible due to 'private' protection level
```
<!-- END:output swift-reject -->

**Kotlin** — `private` confines a member to its class:

<!-- BEGIN:output kotlin-reject -->
```text
reject.kt:12:15: error: cannot access 'var balance: Long': it is private in 'Account'.
```
<!-- END:output kotlin-reject -->

## Hiding the representation: opaque struct vs abstract type

There is a real expressiveness gradient in what "hidden" means. Go and Swift and
Kotlin hide by making the FIELDS inaccessible: the type name is exported, and the
representation happens to be a struct/class the caller cannot read into. Rust is
finer-grained: each item carries its own `pub` / `pub(crate)` / `pub(super)`
level, so you choose the exact radius of every field and function.

OCaml's `.mli` is the strongest and the cleanest: `type t` in the signature is
genuinely ABSTRACT — the client sees a name with no representation at all, cannot
tell it is a record, cannot build one with a literal, and cannot pattern-match
it. The implementation's concrete type is checked against the signature but never
escapes it:

<!-- BEGIN:snippet ocaml-mli -->
```ocaml
type t
(** an account; representation hidden *)

exception Overdraft

val open_ : int -> t
(** [open_ initial] raises Overdraft if initial < 0 *)

val deposit : t -> int -> t
val withdraw : t -> int -> t
val balance : t -> int
```
<!-- END:snippet ocaml-mli -->

This is signature ascription, the same mechanism ML functors are built on. Go
approximates it with an unexported struct or an interface, but the type name is
always exported alongside a concrete (if opaque) struct; OCaml can export the
name while the representation is formally unknowable. See
[`COMPARISON.md`](COMPARISON.md) for the full ladder.

## Exercises

[`exercises/go`](exercises/go) is red until you build an opaque, always-valid
`Ratio` type whose invariant (reduced form, non-zero denominator) cannot be
broken from outside the package. [`solutions/go`](solutions/go) is the verified
corrigé:

```
make exercises M=08   # red
make solutions M=08   # green
```

## References

Official sources first, grouped by language.

### Go

- Effective Go, "Names" (exported = Capitalized): https://go.dev/doc/effective_go#names
- The Go spec, "Exported identifiers": https://go.dev/ref/spec#Exported_identifiers
- `go` command docs, "Internal Directories": https://pkg.go.dev/cmd/go#hdr-Internal_Directories
- Organizing a Go module: https://go.dev/doc/modules/layout

### Rust

- The Rust Reference, "Visibility and privacy": https://doc.rust-lang.org/reference/visibility-and-privacy.html
- The Book, "Controlling scope and privacy with pub": https://doc.rust-lang.org/book/ch07-03-paths-for-referring-to-an-item-in-the-module-tree.html

### OCaml

- OCaml Manual, "Signatures" (module types, abstract types): https://ocaml.org/manual/5.4/moduleexamples.html#s%3Asignature
- OCaml Manual, "The module system": https://ocaml.org/manual/5.4/moduleexamples.html

### Swift

- The Swift Programming Language, "Access Control": https://docs.swift.org/swift-book/documentation/the-swift-programming-language/accesscontrol/

### Kotlin (JVM)

- Kotlin docs, "Visibility modifiers": https://kotlinlang.org/docs/visibility-modifiers.html
