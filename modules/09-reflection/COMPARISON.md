# 09 — Comparison: five ways to be generic over a type

**Environment.** Go 1.26.5, Rust 1.92.0, OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10.
Every demo output and reject error is injected into the README from the code.

## The taxonomy

| Language | Mechanism | When | Runtime type info? | Works on any type? | Error timing |
| -------- | --------- | ---- | ------------------ | ------------------ | ------------ |
| Go       | reflection (`reflect`, `encoding/json`) | run time | yes (carried in the value) | yes, no declaration needed | run time (or `go vet`) |
| Rust     | `#[derive]` procedural macro | compile time | no | only types you derive for | compile time |
| Swift    | protocol conformance synthesis | compile time | no | only types you declare conformance for | compile time |
| Kotlin   | `data class` generated members | compile time | no (JVM reflection exists but is separate) | only `data class`es | compile time |
| OCaml    | hand-written, or a ppx macro | compile time | no | only types you write/derive for | compile time |

The single dividing line: Go carries the type descriptor in the value and
inspects it at run time; the other four burn a specialized implementation for
each type into the binary at compile time and have no runtime type descriptor to
walk.

## What each side buys and loses

**Reflection (Go)** buys maximum flexibility: `Describe` and `json.Marshal` work
on a type defined in another package, or one that did not exist when the code was
written, with zero ceremony. It loses compile-time safety (a bad tag or an
unhandled shape surfaces at run time, or silently does nothing) and pays a
per-call cost, even though `encoding/json` amortizes the reflection walk with a
cached per-type encoder.

**Compile-time codegen (the other four)** buys speed (the specialized path is
generated, no walk, no indirection) and safety (a type that cannot be handled is
a compile error — see all four rejects). It loses openness: the operation exists
only for types you asked for, and adding it to a foreign type needs a newtype
wrapper (Rust), an extension (Swift), or a ppx annotation (OCaml). It also needs
the build step that generates the code.

## OCaml is the instructive outlier

OCaml sits at neither pole in its standard library: no runtime reflection (you
genuinely cannot walk a value's type at run time — types are erased) and no
built-in derive. The everyday answer is a ppx (`[@@deriving show, eq, yojson]`),
which is compile-time codegen supplied by a syntactic preprocessor rather than
the compiler itself. This repo hand-writes the operations to stay
dependency-free, which is the honest third option and a reminder that "generic
over a type" is not free anywhere — someone or something writes the per-type
code; the only question is who (you, a macro, or the runtime) and when.

## The proc-macro / synthesis family

Rust `derive`, Swift synthesis, Kotlin `data class`, and OCaml ppx are one idea
in four dresses: a compile-time program that reads a type's structure and emits
code. Rust and OCaml expose it as user-programmable macros (you can write your
own `derive` / ppx); Swift and Kotlin bake a fixed menu into the compiler
(`Equatable`/`Hashable`/`Codable`; `equals`/`hashCode`/`toString`/`copy`). Go
deliberately has no such compile-time metaprogramming in the language — its
answer to "generate code" is `go generate`, an external tool that writes `.go`
files you commit, keeping the language small at the cost of a separate step.

## Bottom line

For a closed set of your own types on a hot path, compile-time codegen wins on
both speed and safety. For open-ended, whenever-you-need-it genericity over
arbitrary types (a debugger, a generic logger, a config loader that accepts any
struct), reflection is the tool, and its cost is a constant factor that
`encoding/json`'s caching keeps modest. Go's choice — reflection plus `go
generate`, and no macro system — is the same minimalism-over-power trade seen
across these modules: fewer language features, and a static analyzer (`go vet`)
to recover the safety the dynamic mechanism gave up.
