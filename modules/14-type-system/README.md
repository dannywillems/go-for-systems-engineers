# 14 — Type-system quirks and abstractions

**Thesis.** Each type system draws the line between "expressible" and "not"
somewhere different. This module shows one showcase construct per language,
compiling, and next to each a program that is *rejected* — because a type system
is defined as much by what it refuses as by what it accepts. Go's small system
still hides a trick (phantom type parameters); OCaml's GADTs express a typed
evaluator the others cannot match without heavy machinery.

## Contents

- [Go: phantom type parameters](#go-phantom-type-parameters)
- [Rust: typestate](#rust-typestate)
- [OCaml: GADTs](#ocaml-gadts)
- [Swift and Kotlin](#swift-and-kotlin)
- [What each cannot express](#what-each-cannot-express)
- [References](#references)

## Go: phantom type parameters

A type parameter used only at the type level makes `ID[User]` and `ID[Order]`
distinct types though both are `int64`:

<!-- BEGIN:snippet go-phantom -->
```go
// ID[Tag] is a PHANTOM-typed identifier: the Tag type parameter appears only in
// the type, never in a field, so it costs nothing at runtime (an int64) yet
// makes ID[User] and ID[Order] distinct types the compiler refuses to mix.
type ID[Tag any] int64

type User struct{}

type Order struct{}

type UserID = ID[User]

type OrderID = ID[Order]

// LookupUser accepts only a UserID; passing an OrderID is a compile error, even
// though both are int64. See reject-go for the rejection.
func LookupUser(id UserID) string {
	return fmt.Sprintf("user #%d", int64(id))
}
```
<!-- END:snippet go-phantom -->

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
user #7
UserID(7) and OrderID(7) share underlying int64: true
```
<!-- END:output go-demo -->

Mixing them is a compile error:

<!-- BEGIN:output go-reject -->
```text
./badphantom.go:18:20: cannot use OrderID(7) (constant 7 of int64 type OrderID) as UserID value in argument to LookupUser
```
<!-- END:output go-reject -->

## Rust: typestate

The object's state lives in its type, so an invalid transition does not
typecheck (the method is absent for that state):

<!-- BEGIN:snippet rust-typestate -->
```rust
pub struct Open;
pub struct Closed;

pub struct Door<State> {
    _state: PhantomData<State>,
}

impl Door<Closed> {
    pub fn closed() -> Self {
        Door {
            _state: PhantomData,
        }
    }
    pub fn open(self) -> Door<Open> {
        Door {
            _state: PhantomData,
        }
    }
}

impl Door<Open> {
    pub fn close(self) -> Door<Closed> {
        Door {
            _state: PhantomData,
        }
    }
}
```
<!-- END:snippet rust-typestate -->

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
Door: closed -> open -> close typechecks (typestate)
```
<!-- END:output rust-demo -->

<!-- BEGIN:output rust-reject -->
```text
error[E0599]: no method named `open` found for struct `Door<Open>` in the current scope
```
<!-- END:output rust-reject -->

## OCaml: GADTs

The showcase. The constructor's return type is indexed by the value's type, so a
single `eval : 'a expr -> 'a` returns the right type per constructor — no tags,
no runtime type test, and the `If` branches are forced to agree statically:

<!-- BEGIN:snippet ocaml-gadt -->
```ocaml
type _ expr =
  | Int : int -> int expr
  | Bool : bool -> bool expr
  | Add : int expr * int expr -> int expr
  | If : bool expr * 'a expr * 'a expr -> 'a expr
  | Eq : int expr * int expr -> bool expr

(* [eval : 'a expr -> 'a]. The return type varies with the constructor, tracked
   by the GADT index; the [If] branches are forced to agree, statically. *)
let rec eval : type a. a expr -> a = function
  | Int n -> n
  | Bool b -> b
  | Add (x, y) -> eval x + eval y
  | If (c, t, e) -> if eval c then eval t else eval e
  | Eq (x, y) -> eval x = eval y
```
<!-- END:snippet ocaml-gadt -->

<!-- BEGIN:output ocaml-demo -->
```text
$ dune exec bin/main.exe
eval (if 2=2 then 1+3 else 0) = 4
eval (3 = 4) = false
```
<!-- END:output ocaml-demo -->

The same `eval` returns an `int` for an `int expr` and a `bool` for a `bool
expr`. Expressing this in Go is impossible; in Rust it needs a visitor or
trait-object machinery.

## Swift and Kotlin

Swift expresses the same typestate via phantom generics + constrained
extensions; Kotlin's showcase is declaration-site variance (`out`/`in`):

<!-- BEGIN:output swift-demo -->
```text
Door<Closed>.open().close() typechecks (Swift typestate)
```
<!-- END:output swift-demo -->

<!-- BEGIN:output swift-reject -->
```text
reject.swift:15:16: error: referencing instance method 'open()' on 'Door' requires the types 'Open' and 'Closed' be equivalent
```
<!-- END:output swift-reject -->

<!-- BEGIN:output kotlin-demo -->
```text
firstName(catProducer) = cat (covariance via out)
```
<!-- END:output kotlin-demo -->

<!-- BEGIN:output kotlin-reject -->
```text
reject.kt:6:23: error: type parameter 'T' is declared as 'out' but occurs in 'in' position in type 'T (of interface Bad<out T>)'.
```
<!-- END:output kotlin-reject -->

## What each cannot express

- **Go**: no higher-kinded types (Module 03 reject), no generic methods
  (Module 04 reject), no GADTs, no declaration-site variance. Phantom type
  parameters and structural interfaces are the ceiling.
- **Rust**: very expressive (GATs, const generics, typestate) but no true HKT
  without workarounds; specialization is unstable.
- **OCaml**: GADTs, polymorphic variants, first-class modules (explicit
  existentials), and functors (higher-kinded over modules) — the richest here
  for type-level modelling.
- **Swift**: associated types, `some`/`any`, conditional conformance; no HKT.
- **Kotlin**: variance and reified generics, but JVM erasure limits runtime type
  reflection.

See [`COMPARISON.md`](COMPARISON.md) for the expressiveness matrix.

## References

Official sources first, grouped by language.

### Go

- Go spec, Type parameters: https://go.dev/ref/spec#Type_parameter_declarations

### Rust

- The Rustonomicon, PhantomData: https://doc.rust-lang.org/nomicon/phantom-data.html
- The Embedded Rust Book, Typestate programming: https://docs.rust-embedded.org/book/static-guarantees/typestate-programming.html

### OCaml

- OCaml Manual, Generalized algebraic datatypes (GADTs): https://ocaml.org/manual/5.4/gadts.html
- OCaml Manual, Polymorphic variants: https://ocaml.org/manual/5.4/polyvariant.html

### Swift

- Swift, Generics (associated types, constraints): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/generics/
- Swift, Opaque and boxed protocol types (`some`/`any`): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/opaquetypes/

### Kotlin

- Kotlin, Generics: variance (`in`/`out`): https://kotlinlang.org/docs/generics.html#variance
