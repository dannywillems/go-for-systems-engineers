# 14b — Propositions as types: what each system can and cannot state

A type system is a **logic**, and a program is a **proof** (the Curry–Howard
correspondence: a value of type `A` is evidence for the proposition `A`). So the
question "what can this type system express?" is the question "what can this
logic *state and prove*?" — and its **limitations** are the propositions it
cannot phrase at all. This document catalogs the propositions type theorists
care about, from propositional logic up to full dependent types, says which
well-known problem each solves, and gives the verdict for all five languages
plus a reference column for the dependently-typed languages that reach the
frontier.

Verdict legend: **✓** first-class / native · **~** encodable with a named
workaround, or partial · **✗** not expressible in the type system (a runtime
check or discipline is the only recourse). Every version-sensitive cell below
was checked against the pinned toolchains (Go 1.26.5, Rust 1.92.0, OCaml 5.4.0,
Swift 6.2.3, Kotlin 2.4.10); the boldest ones are backed by compiled code in
this module (see "Verified in this module").

## Contents

- [I. Propositional logic — the simply-typed core](#i-propositional-logic--the-simply-typed-core)
- [II. Quantifiers — System F and Fω](#ii-quantifiers--system-f-and-f)
- [III. Subtyping, variance, and extensibility](#iii-subtyping-variance-and-extensibility)
- [IV. Indexed types — toward dependency](#iv-indexed-types--toward-dependency)
- [V. Substructural types and effects](#v-substructural-types-and-effects)
- [VI. The dependent frontier — where all five stop](#vi-the-dependent-frontier--where-all-five-stop)
- [The master matrix](#the-master-matrix)
- [Where each ceiling is](#where-each-ceiling-is)
- [Verified in this module](#verified-in-this-module)
- [References](#references)

## I. Propositional logic — the simply-typed core

The simply-typed lambda calculus corresponds to the implicational fragment of
intuitionistic propositional logic; add products and sums and you have the whole
`{∧, ∨, →, ⊤, ⊥}` fragment. Every mainstream language covers most of it — the one
real gap is Go's missing coproduct.

- **Conjunction `A ∧ B`** = a **product** (struct / tuple / record). Evidence for
  both at once. _Useful for:_ bundling related data. _All five._
- **Implication `A → B`** = a **function**. Evidence transformer. _All five._
- **Truth `⊤`** = the **unit** type (`()` / `struct{}` / `Unit` / `Void` in
  Swift). The trivially-provable proposition. _All five._
- **Disjunction `A ∨ B`** = a **coproduct / sum** (tagged union). Evidence for one
  of two, and you must handle both. _Useful for:_ ASTs, state machines,
  error-or-value, "one of N shapes." _Go **✗** (Module 02; a sealed interface +
  analyzer only); Rust/OCaml/Swift/Kotlin ✓._
- **Totality of case analysis (exhaustiveness)** = the sum's **eliminator is
  total**; a forgotten case is a type error. _Useful for:_ never silently
  mishandling a variant. _Go **~** (the `go/analysis` pass of Module 02); others
  ✓ (the compiler proves it)._
- **Falsehood `⊥`** = an **uninhabited / bottom type** (`!` Rust, `Never` Swift,
  `Nothing` Kotlin, `type void = |` OCaml). No value exists, so a function
  returning it never returns; it is a subtype of every type, typing `exit`,
  `panic`, and the impossible branch. _Useful for:_ total functions, typed
  unreachability, "this can't happen" that the checker believes. _Go **✗** (no
  bottom type — a panicking function still has an ordinary return type); Rust,
  OCaml, Swift, Kotlin ✓._
- **Negation `¬A = A → ⊥`** = a function into the bottom type; combined with `⊥`,
  lets you write refutations (`absurd : void → 'a`). _Same support as `⊥`._

## II. Quantifiers — System F and Fω

Adding quantification over types is the jump to second-order logic. `∀` is
polymorphism; `∃` is abstraction; quantifying over *type constructors* is Fω, the
higher-kinded frontier where Go, Swift, and Kotlin stop.

- **Universal `∀X. …`** = **parametric polymorphism** (generics). One
  implementation for all element types. _Useful for:_ containers, generic
  algorithms. _All five (Module 03)._
- **Parametricity / free theorems** — from a polymorphic *type alone* you get a
  theorem for free: the only total, pure inhabitant of `∀X. X → X` is the
  identity; `∀X. List X → List X` can only permute/drop/duplicate, never invent
  an element (Reynolds; Wadler, "Theorems for free"). _Useful for:_ reasoning and
  security ("this function provably cannot look at your data"). _All five **~**:
  the theorem holds only for the reflection-free fragment — Go's `reflect`, Rust's
  `Any`/`TypeId`, OCaml's `Obj.magic`, Swift's `Mirror`, and JVM reflection each
  break it._
- **Existential `∃X. …`** = **abstract data type / interface**: "there is some
  type with these operations; you may not see which." _Useful for:_ plugins,
  representation hiding (Module 08). _Go ✓ (interfaces, implicit); Rust ✓
  (`dyn`/`impl`); OCaml ✓ (first-class modules, explicit); Swift ✓ (`any`/`some`);
  Kotlin ✓ (interfaces)._
- **Rank-N polymorphism** — a `∀` to the left of an arrow: a function that takes a
  *polymorphic* function as an argument and uses it at several types. _Useful
  for:_ the ST monad's `runST`, capability-safe callbacks. _OCaml ✓ (polymorphic
  record fields / first-class modules); Rust **~** (only for lifetimes, via
  `for<'a>` HRTB, never for types); Go, Swift, Kotlin ✗._
- **Higher-kinded types `∀F:*→*. …`** = quantifying over a **type constructor**
  (System Fω). The prerequisite for a *generic* `Functor`/`Monad`/`Traversable`
  and for do-notation that works over any container. _Useful for:_ writing `map`,
  `traverse`, `sequence` once for all containers. _Go ✗, Swift ✗, Kotlin ✗
  (a type parameter cannot be applied like `F[A]` — Module 03's HKT reject); Rust
  **~** (GATs approximate it); OCaml **~** (functors give it at the *module*
  level, or the defunctionalized `('a, 'f) app` encoding). Haskell/Scala ✓._
- **Type classes with coherence** = **ad-hoc polymorphism** made principled: a
  dictionary of operations selected by type, with a global-uniqueness guarantee.
  _Useful for:_ `Ord`, `Num`, `Serialize` without passing dictionaries by hand.
  _Rust ✓ (traits + the orphan rule = coherence, Module 01); Swift **~**
  (protocols, but retroactive conformance is allowed); OCaml **~** (modules /
  functors, explicit, no coherence); Kotlin **~** (interfaces, plus context
  parameters); Go ✗ (methods are structural, no lawful class, no dictionary)._
- **Associated types / type families** = a **type-level function** attached to an
  instance (`trait Iterator { type Item; }`). _Useful for:_ an iterator's element
  type, a collection's key type. _Rust ✓, Swift ✓ (associated types), OCaml ✓
  (abstract type members of a module signature); Kotlin ✗, Go ✗._

## III. Subtyping, variance, and extensibility

- **Subtyping `A <: B`** = coercion without conversion. _Nominal_ (class
  hierarchies) or _structural_ (shape). _Swift ✓ / Kotlin ✓ (class + interface,
  nominal); OCaml ✓ (objects and polymorphic variants, structural); Go **~**
  (interface satisfaction is a restricted structural subtyping — Module 01 — but
  there is no subtyping *between* named types); Rust **~** (only lifetimes are
  subtyped)._
- **Variance (covariance / contravariance)** — when does `F<Cat> <: F<Animal>`?
  The classic rule: a function is contravariant in its argument, covariant in its
  result. _Useful for:_ read-only covariant collections, safe callbacks.
  _Declaration-site: OCaml ✓ (`+'a`/`-'a`), Kotlin ✓ (`out`/`in`). Rust: variance
  is **inferred** structurally (not declared, not use-site). Go: generics are
  **invariant** (✗). Swift: ✗ for generics._
- **Bounded / F-bounded quantification** — `∀X <: C. …`, and the recursive
  `X : C<X>` ("comparable to my own type"). _Useful for:_ "you can sort a slice
  only if its elements are ordered." _All five ✓ for bounds (Go `constraints`,
  Rust/Swift/Kotlin trait/protocol/upper bounds, OCaml module constraints);
  F-bounded ✓ in Rust/Swift/Kotlin/Go, **~** in OCaml (recursive module types)._
- **Row polymorphism** = "any record/object with *at least* these fields." _Useful
  for:_ extensible records, config that accepts supersets. _OCaml ✓ (object types
  and polymorphic variants carry a row variable); Go **~** (structural interfaces
  are a method-only row); Rust, Swift, Kotlin ✗._
- **Polymorphic variants (open / extensible sums)** = sum types whose variant set
  is not closed, reusable across types. _Useful for:_ growing an AST without
  editing a central `enum`. _OCaml ✓ (`` [`A | `B ] ``) only; the others have only
  closed sums._
- **The expression problem** — extend both the data variants *and* the operations
  of a datatype without recompiling existing code, keeping static type safety
  (Wadler). _Useful for:_ extensible compilers/interpreters. _OCaml ✓ (polymorphic
  variants; or tagless-final); Rust **~**, Swift **~**, Kotlin **~**, Go **~**
  (each solves one direction cheaply — new ops *or* new variants — and pays for
  the other; tagless-final via traits/protocols/interfaces recovers both with
  boilerplate)._

## IV. Indexed types — toward dependency

Here types start to depend on *other types* or on *values*, the on-ramp to
dependent types. This is where OCaml (GADTs) and Rust/Swift (const/value
generics) pull ahead, and where the truly dependent propositions begin to appear
just out of reach.

- **GADTs (type-indexed constructors)** — a constructor whose *return type* refines
  the type index, so a single `eval : 'a expr → 'a` gives an `int` for an
  `int expr` and a `bool` for a `bool expr`, checked statically (Module 14's
  showcase). _Useful for:_ well-typed interpreters, type-safe serializers,
  tagless-final. _OCaml ✓; Rust **~** (trait-object / visitor encodings, not
  per-variant refinement); Go, Swift, Kotlin ✗._
- **Propositional type equality `a ≡ b`** — a *value* that is evidence two types
  are equal (`type (_,_) eq = Refl : ('a,'a) eq`), letting you safely coerce
  `a → b` with **no** `unsafe`/`Obj.magic`. The identity-type fragment of a
  dependent theory, at the type level. _Useful for:_ safe dynamic casts,
  heterogeneous maps keyed by type, GADT plumbing. _OCaml ✓ (verified below);
  Rust, Swift, Kotlin, Go ✗ (you cannot construct a first-class witness that two
  types are the same)._
- **Type-level naturals / const (value) generics** — a *number* in the type,
  `[T; N]` / `Vector<N, T>`, with the compiler doing arithmetic on sizes. _Useful
  for:_ fixed-size crypto buffers, matrix dimensions. _Rust ✓ (const generics),
  Swift ✓ (value generics, Swift 6.2); OCaml **~** (Peano encoded in a GADT);
  Kotlin ✗, Go ✗._
- **Length-indexed vectors (safe `head`/`zip`)** — the length is in the type, so
  `zip` on mismatched lengths is a *compile* error and `head` of a non-empty
  vector needs no runtime check. _Useful for:_ eliminating a class of
  out-of-bounds and length-mismatch bugs. _Rust ✓ (`[T; N]` + const generics —
  verified reject below); Swift **~** (value generics, maturing); OCaml **~**
  (GADT Peano); Kotlin, Go ✗._
- **Refinement types `{x : T | P x}`** — a type carrying a *logical predicate*
  (`Pos = {n | n > 0}`, `NonEmpty`, `Sorted`). _Useful for:_ validated input,
  non-empty lists, positive amounts — checked by the type checker, not at runtime.
  _None of the five ✗ (all approximate with a newtype + smart constructor, whose
  invariant is a **runtime** check the type system does not know). Native in
  Liquid Haskell, F*, and refinement-typed dialects._
- **Units of measure** — `3.0<m/s>`, with dimensional arithmetic checked. _Useful
  for:_ physics/finance safety (the Mars Climate Orbiter was lost to a
  units bug). _All five **~** via phantom types (Module 14) + newtypes; native
  only in F#._

## V. Substructural types and effects

Structural logic lets you use a hypothesis any number of times. *Substructural*
logics restrict that — and that restriction is exactly what models resources
(you can spend a coin once) and effects (this code may throw). This is Rust's and
Swift's home turf; Go and Kotlin sit it out.

- **Affine types (use *at most* once)** = Rust's **ownership / move semantics**:
  after a move the old binding is gone, so a value is consumed at most once.
  _Useful for:_ files/sockets/locks that must not be used after close, no
  double-free, no use-after-move. _Rust ✓; Swift ✓ (`~Copyable` noncopyable
  types, verified); OCaml ✗ upstream (**~** with the OxCaml mode extension);
  Kotlin, Go ✗._
- **Linear types (use *exactly* once)** — stronger: the resource must also not be
  *dropped* silently. _Useful for:_ protocols that must be completed, must-consume
  handles. _Rust **~** (`#[must_use]` + affine approximates), Swift **~**
  (`consuming` + `~Copyable`); native in Haskell (LinearTypes) and Idris 2; OCaml,
  Kotlin, Go ✗._
- **Regions / lifetimes (no dangling reference)** — a type-level scope proving a
  borrow does not outlive its owner. _Useful for:_ memory safety without a GC.
  _Rust ✓ (`'a` lifetimes); Swift **~** (`~Escapable` + lifetime dependencies,
  emerging); OCaml, Kotlin, Go ✗ (a tracing GC removes the *need* but offers no
  type-level region)._
- **Thread-safety in the type (`Send`/`Sync`/`Sendable`)** — data-race freedom
  proved by typing which values may cross threads (Module 06). _Rust ✓, Swift ✓
  (`Sendable`); OCaml, Kotlin, Go ✗ (races are a runtime concern)._
- **Typestate / protocol state machines** — the object's *state* lives in its
  type, so an illegal transition (read before open) does not typecheck (Module
  14). _Useful for:_ builders that must set required fields, connections opened
  before use. _Rust ✓, Swift ✓, OCaml ✓ (phantom); Go **~**, Kotlin **~**
  (phantom, more limited)._
- **Session types** — a *channel's* protocol as a type ("send Int, then receive
  Bool, then close"), with deviation a compile error. _Useful for:_ communication
  protocol conformance. _All five **~** (library encodings exist in Rust/OCaml/
  Scala); native in session-typed research languages._
- **Effect systems / typed effects** — the type records which effects a function
  may perform (throw, block, do IO, allocate). _Useful for:_ knowing what a call
  can do, taming `async` "function coloring." _Swift **~/✓** (`throws(E)` typed
  throws + `async` + `rethrows` — a real if narrow effect system, verified);
  Kotlin **~** (`suspend` colors async); Rust **~** (`Result`/`async` encode two
  effects, no effect polymorphism); OCaml **~** (effect *handlers* in OCaml 5, but
  the effects are **untyped** — absent from signatures); Go ✗. Full algebraic
  effect systems: Koka, Eff, Frank._
- **Immutability in the type** — `const`/`readonly`/`let` that the checker
  enforces. _Rust ✓ (`&` vs `&mut`, bindings immutable by default), OCaml ✓
  (immutable by default), Swift ✓ (`let`, `borrowing`); Kotlin **~** (`val` binds
  a reference but does not deep-freeze); Go **~** (`const` for scalars only, no
  immutable struct/slice)._

## VI. The dependent frontier — where all five stop

Full **dependent types** let a type depend on a *value*, unifying types and
propositions completely. None of the five languages reach this; the reference
column is Agda, Idris 2, Rocq (Coq), Lean 4, and F*.

- **Dependent function `Π(x:A). B(x)`** — the return *type* depends on the argument
  *value*: `replicate : (n : Nat) → a → Vec n a`, or a `printf` whose argument
  types are computed from the format value. _All five ✗._
- **Dependent pair `Σ(x:A). B(x)`** — a value bundled with a *proof* about it: `(n
  : Nat, Vec n a)`, or "an index together with evidence it is in bounds." _All
  five ✗._
- **The identity type + `J` (equality with transport over value families)** —
  full propositional equality of *values*, the substrate of equational proofs.
  _All five ✗_ (OCaml's `(_,_) eq` is only the *type*-equality fragment).
- **Termination / totality checking** — the compiler proves recursion terminates,
  so a proof cannot cheat by looping. _All five ✗._
- **Universe hierarchy (`Type : Type₁ : …`) and proof-carrying code / program
  extraction** — types as first-class values, and extracting a verified program
  from its proof. _All five ✗._
- **(Bonus) Turing-complete type-level computation** — even without dependency,
  some checkers compute arbitrarily at the type level: _Rust ✓ (trait resolution
  is Turing-complete), OCaml **~** (GADTs + functors compute a great deal); Swift,
  Kotlin, Go ✗. (C++ templates are the classic ✓.)_

## The master matrix

`✓` native · `~` encodable / partial (mechanism in the entry above) · `✗` not in
the type system. "Dep." = the dependently-typed reference languages (Agda /
Idris 2 / Rocq / Lean 4 / F*).

| # | Proposition (type) | Go | Rust | OCaml | Swift | Kotlin | Dep. |
|---|--------------------|----|------|-------|-------|--------|------|
| 1 | Product `A ∧ B` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 2 | Function `A → B` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 3 | Unit `⊤` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 4 | Sum / coproduct `A ∨ B` | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 5 | Exhaustiveness (total eliminator) | ~ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 6 | Bottom `⊥` / `Never` | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 7 | Negation `¬A = A → ⊥` | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 8 | Universal `∀X` (generics) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 9 | Parametricity / free theorems | ~ | ~ | ~ | ~ | ~ | ✓ |
| 10 | Existential `∃X` (ADT / interface) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 11 | Rank-N polymorphism | ✗ | ~ | ✓ | ✗ | ✗ | ✓ |
| 12 | Higher-kinded types `∀F:*→*` | ✗ | ~ | ~ | ✗ | ✗ | ✓ |
| 13 | Type classes + coherence | ✗ | ✓ | ~ | ~ | ~ | ✓ |
| 14 | Associated types / type families | ✗ | ✓ | ✓ | ✓ | ✗ | ✓ |
| 15 | Subtyping `A <: B` | ~ | ~ | ✓ | ✓ | ✓ | ✓ |
| 16 | Declaration-site variance | ✗ | ✗ | ✓ | ✗ | ✓ | ✓ |
| 17 | Bounded / F-bounded quantification | ✓ | ✓ | ~ | ✓ | ✓ | ✓ |
| 18 | Row polymorphism | ~ | ✗ | ✓ | ✗ | ✗ | ✓ |
| 19 | Polymorphic (open) variants | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ |
| 20 | Expression problem (both axes) | ~ | ~ | ✓ | ~ | ~ | ✓ |
| 21 | GADTs (type-indexed constructors) | ✗ | ~ | ✓ | ✗ | ✗ | ✓ |
| 22 | Propositional TYPE equality `a ≡ b` | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ |
| 23 | Type-level naturals / value generics | ✗ | ✓ | ~ | ✓ | ✗ | ✓ |
| 24 | Length-indexed vectors (safe zip) | ✗ | ✓ | ~ | ~ | ✗ | ✓ |
| 25 | Refinement types `{x \| P x}` | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| 26 | Units of measure | ~ | ~ | ~ | ~ | ~ | ✓ |
| 27 | Affine (use ≤ once) | ✗ | ✓ | ✗ | ✓ | ✗ | ✓ |
| 28 | Linear (use = once) | ✗ | ~ | ✗ | ~ | ✗ | ✓ |
| 29 | Regions / lifetimes | ✗ | ✓ | ✗ | ~ | ✗ | ✓ |
| 30 | Thread-safety (`Send`/`Sendable`) | ✗ | ✓ | ✗ | ✓ | ✗ | ✓ |
| 31 | Typestate / protocol state | ~ | ✓ | ✓ | ✓ | ~ | ✓ |
| 32 | Session types | ~ | ~ | ~ | ~ | ~ | ✓ |
| 33 | Effect system / typed effects | ✗ | ~ | ~ | ~ | ~ | ✓ |
| 34 | Immutability in the type | ~ | ✓ | ✓ | ✓ | ~ | ✓ |
| 35 | Phantom types | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 36 | Dependent function `Π(x:A).B(x)` | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| 37 | Dependent pair `Σ(x:A).B(x)` | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| 38 | Identity type + `J` / value proofs | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| 39 | Termination / totality checking | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| 40 | Proof-carrying code / extraction | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |

Rough `✓` counts (native only): **Go 9**, **Kotlin 13**, **Swift 20**,
**Rust 21**, **OCaml 22** — and the dependent column, 40. The counts are a crude
proxy (a `~` often matters more than a `✓`), but the shape is real: Go sits
lowest by deliberate design, OCaml and Rust reach highest along *different* axes
(OCaml toward logic and indexing, Rust toward resources and lifetimes), Swift
tracks Rust on substructural features while keeping OO subtyping, and Kotlin is
the pragmatic middle. None of the five crosses into dependency.

## Where each ceiling is

- **Go** stops at the propositional core minus the coproduct: products,
  functions, `∀`/`∃`, phantom types, bounded generics — and no sum type, no
  bottom, no variance, no HKT, no GADTs, no substructural anything. This is the
  deliberate "small language" trade of every module here; the missing pieces are
  recovered (if at all) by external tools and runtime discipline.
- **Rust** climbs the substructural axis (affine ownership, lifetimes,
  `Send`/`Sync`) and has coherent type classes, const generics, and a
  Turing-complete type checker — but has no GADTs, no rank-N over types, no true
  HKT, and no dependent types.
- **OCaml** climbs the logical/indexing axis: GADTs, propositional type equality,
  polymorphic variants, row polymorphism, first-class-module existentials,
  declaration-site variance, rank-N — the richest here for type-level *modelling*
  — but has no ownership/affinity (upstream) and no dependency.
- **Swift** mirrors much of Rust's substructural set (`~Copyable`, `Sendable`,
  emerging lifetimes) and adds value generics and typed throws, on top of nominal
  subtyping — but no HKT, no GADTs, no variance.
- **Kotlin** offers declaration-site variance, a bottom type, and sealed
  exhaustive sums, but is capped by the JVM (erasure, no value types) and reaches
  none of the substructural or indexed features.
- **The frontier all five miss** is the dependent one: `Π`/`Σ`, the identity type
  with value-level proofs, refinement types, termination checking, and
  proof-carrying code — the province of Agda, Idris 2, Rocq, Lean 4, and F*.
  The lesson is the standard expressiveness/tractability trade: the more
  propositions a type system can state, the harder (eventually undecidable) type
  inference and checking become, which is precisely why production languages stop
  short of full dependency.

## Verified in this module

The claims that are easy to get wrong are backed by compiled code, checked by the
docs-freshness gate:

- **GADTs** (row 21): OCaml's `eval : 'a expr → 'a` — the `ocaml-gadt` snippet and
  the `ocaml-demo` output in [`README.md`](README.md).
- **Propositional type equality** (row 22): OCaml's `(_,_) eq = Refl` and the
  `Obj.magic`-free cast — the `ocaml-eq` snippet in the README; a test asserts the
  cast round-trips.
- **Type-level naturals / length-indexed vectors** (rows 23–24): Rust const
  generics — the `rust-lenvec` snippet (a `zip` requiring equal `N`) and the
  `rust-len-reject` output, where a mismatched-length call is rejected with
  `error[E0308]: expected an array with a size of 3, found one with a size of 2`.
- **Phantom types, typestate, variance, and the four compile-rejects** (rows 16,
  31, 35) are the original showcases of this module.

The rest of the matrix is a survey; where a cell says `~`, the entry above names
the exact workaround (a functor, a newtype + smart constructor, a phantom type, a
library), so nothing is hand-waved.

## References

Foundational:

- Howard, "The formulae-as-types notion of construction" (Curry–Howard, 1980).
- Wadler, "Propositions as Types" (2015): https://homepages.inf.ed.ac.uk/wadler/papers/propositions-as-types/propositions-as-types.pdf
- Girard, *Proofs and Types* (System F / Fω): https://www.paultaylor.eu/stable/prot.pdf
- Reynolds, "Types, Abstraction and Parametric Polymorphism" (1983); Wadler,
  "Theorems for Free!" (1989): https://homepages.inf.ed.ac.uk/wadler/papers/free/free.pdf
- Pierce, *Types and Programming Languages* (TAPL), and Harper, *Practical
  Foundations for Programming Languages* (PFPL).
- Wadler, "The Expression Problem" (1998): https://homepages.inf.ed.ac.uk/wadler/papers/expression/expression.txt
- Yallop & White, "Lightweight Higher-Kinded Polymorphism" (the OCaml `app`
  encoding, 2014): https://www.cl.cam.ac.uk/~jdy22/papers/lightweight-higher-kinded-polymorphism.pdf

Per language:

- Rust: const generics + GATs (Reference): https://doc.rust-lang.org/reference/items/generics.html ; ownership/lifetimes: https://doc.rust-lang.org/book/ch04-00-understanding-ownership.html
- OCaml: GADTs and first-class modules: https://ocaml.org/manual/5.4/gadts.html ; polymorphic variants: https://ocaml.org/manual/5.4/polyvariant.html ; effect handlers (untyped): https://ocaml.org/manual/5.4/effects.html
- Swift: noncopyable `~Copyable` (SE-0390), typed throws (SE-0413), value generics
  (SE-0452), `~Escapable` (SE-0446): https://github.com/apple/swift-evolution
- Kotlin: generics/variance and `Nothing`: https://kotlinlang.org/docs/generics.html

Dependent-typed languages: Agda (https://agda.readthedocs.io), Idris 2
(https://idris2.readthedocs.io), Rocq/Coq (https://rocq-prover.org), Lean 4
(https://leanprover.github.io), F* (https://www.fstar-lang.org).
