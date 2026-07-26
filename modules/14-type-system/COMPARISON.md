# 14 — Comparison: type-system expressiveness

**Environment.** Go 1.26.5, Rust 1.92.0, OCaml 5.4.0, Swift 6.2.3, Kotlin
2.4.10. Every "yes/no" below is backed by compiling code or a captured rejection
in this module and Modules 01–04.

## Expressiveness matrix

| Feature | Go | Rust | OCaml | Swift | Kotlin |
| ------- | -- | ---- | ----- | ----- | ------ |
| Parametric generics | yes | yes | yes | yes | yes |
| Generic methods (new type param on a method) | **no** (M04) | yes | yes | yes | yes |
| Higher-kinded types (`F[_]`) | **no** (M03) | no (workarounds) | yes (functors) | no | no |
| Sum types / exhaustive match | **no** (M02, analyzer) | yes | yes | yes | yes |
| GADTs (type-indexed constructors) | no | partial | **yes** | no | no |
| Declaration-site variance | no | no (use-site via bounds) | yes (`+`/`-`) | no | **yes** (`in`/`out`) |
| Phantom types | yes (type params) | yes (`PhantomData`) | yes | yes | yes |
| Typestate (compile-time state machine) | limited | **yes** | yes | yes | limited |
| Existentials | interfaces (implicit) | `dyn`/`impl` | first-class modules (explicit) | `any`/`some` | interfaces |
| Runtime type reflection | `reflect` | limited (`Any`) | limited (`Obj`) | mirrors | limited (erasure) |

## Reading

The systems cluster by lineage. **OCaml** (ML lineage) is the richest for
type-level modelling here: GADTs, polymorphic variants, first-class modules, and
functors give higher-kinded abstraction the others lack. **Rust** trades some of
that (no true HKT) for zero-cost guarantees and typestate, and is the most
expressive of the systems languages. **Swift** and **Kotlin** are pragmatic
object-functional hybrids — Swift leans on protocols + `some`/`any`, Kotlin on
variance + sealed hierarchies, both limited by their runtime (Swift's witness
tables, the JVM's erasure). **Go** is deliberately the smallest: nominal types,
structural interfaces, parametric generics, and no more — its designers traded
expressiveness for a system that is fast to learn and fast to compile, and the
gaps (Modules 02–04) are the measured cost of that trade.

The through-line of the whole repository: Go's type system asks less of the
programmer and guarantees less, compensating with tooling (the exhaustiveness
analyzer, `-race`, `errcheck`) where the other languages use the type checker.
Which trade is right depends on the team and the workload — the point of these
modules is to make the trade concrete rather than tribal.

## References

See the module [`README`](README.md#references) for the per-language sources
(Go spec; Rustonomicon PhantomData + typestate; OCaml GADTs + polymorphic
variants; Swift generics + opaque types; Kotlin variance).
