# go-for-systems-engineers

An executable reference on Go's semantics for an engineer fluent in Rust and
OCaml (and here also Swift and Kotlin). Not a tutorial: no syntax, no "hello
world". Each module isolates a place where Go's semantics, guarantees, or
runtime cost diverge from those languages, demonstrates it with compiled code,
measures it where a number is warranted, and states the theory precisely.

## Status and provenance

This is a living flush of notes, updated over time, not a finished course.
Modules are added, corrected, and expanded as the work continues; expect gaps
and revisions.

**The code is the truth.** The prose explains it, but the compiled, tested,
measured programs are what is authoritative here; when prose and code disagree,
the code wins, and the falsifiability harness below is what keeps the prose
honest against it.

The work is AI-assisted. The author, Danny Willems, focuses on the actual
gathering of knowledge — the semantics, the measurements, and their meaning —
which is the main use of AI here: everything repetitive and of low value to the
human learning process is delegated to the AI for productivity. Concretely, the
AI generates and sketches the repository under the author's direction, to the
shape the author wants, and automates the cumbersome, low-insight scaffolding:
the infrastructure and CI, the multi-language build, and the tooling setup
(formatters, linters, compilers, benchmark harness). The claims are then
verified over time by the author, mechanically via the harness rather than on
trust. Treat this as a learning tool, cross-checked against the linked official
sources, not as an authority.

## The one rule: falsifiability

Every empirical claim in a README is produced by a program in this repo and
injected between markers by [`tools/capture`](tools/capture). Nothing is
hand-typed: not a timing, not a size, not a memory behavior. `make docs`
regenerates the blocks; the `docs-fresh` CI job re-runs it and fails on any
diff. If a measurement contradicts a stated thesis, the thesis is wrong and gets
fixed. See [`modules/00-bootstrap`](modules/00-bootstrap) for the machinery.

Blocks that capture architecture-specific output (assembly, escape analysis,
timings) are marked `portable: false`: regenerated locally and committed, but
excluded from the `docs-fresh` gate because a CI amd64 runner legitimately
differs from the author's arm64. Benchmarks run only under a manual CI job;
their committed `benchstat` summaries are the source of truth.

## Five languages

Go is the subject; the others are the contrast. Each `modules/NN-*` has `go/`,
`rust/`, `ocaml/`, `swift/`, and `kotlin/` subtrees that build and test
independently. Toolchains driving this run:

<!-- BEGIN:output versions -->
```text
# compilers
go                     go version go1.26.5 darwin/arm64
rustc                  rustc 1.92.0 (ded5c06cf 2025-12-08)
cargo                  cargo 1.92.0 (344c4567c 2025-10-21)
ocaml                  The OCaml toplevel, version 5.4.0
swift                  Swift version 6.2.3
kotlinc                info: kotlinc-jvm 2.4.10 (JRE 26.0.1)
java                   openjdk version "26.0.1" 2026-04-21

# formatters
gofmt                  bundled with go
rustfmt                rustfmt 1.8.0-stable (ded5c06cf2 2025-12-08)
ocamlformat            0.28.1
swift-format           bundled with swift (6.2.3)
ktlint                 ktlint version 1.8.0

# linters / static analysis
go vet                 bundled with go
staticcheck            staticcheck 2026.1 (v0.7.0)
golangci-lint          golangci-lint has version 2.12.2 built with go1.26.5 from (unknown, modified: ?, mod sum: "h1:7+d1uY0bq1MU2UV3R5pW5Q7QWdcoq4naMRXM+gsJKrs=") on (unknown)
govulncheck            Go: go1.25.11
clippy                 clippy 0.1.92 (ded5c06cf2 2025-12-08)
ocaml -w               compiler warnings-as-errors (no separate linter is standard)
swift                  compiler diagnostics + swift-format lint
ktlint                 lint + format for Kotlin

# build / bench tooling
dune                   3.20.2
benchstat              golang.org/x/perf (installed)
hyperfine              hyperfine 1.20.0
```
<!-- END:output versions -->

## Usage

```
make setup          # print toolchain status
make build          # build all languages, all modules
make test           # test all
make test-race      # go test -race
make lint           # gofmt/vet/staticcheck/golangci-lint, clippy, ocamlformat,
                    # swift-format, ktlint, shellcheck
make docs           # regenerate captured README blocks
make docs-check     # fail if a portable block is stale (the docs-fresh gate)
make ci             # everything CI runs
make module M=03    # build+test+lint+docs-check one module
make bench          # how to run the (manual) benchmarks
```

Scope any target to one module with `M=NN`. CI (`.github/workflows/ci.yml`)
invokes these same targets, parameterized by version variables, so `make ci`
reproduces it locally.

## Curriculum

See [`PROGRESS.md`](PROGRESS.md) for status and [`VERDICT.md`](VERDICT.md) for
the cross-cutting synthesis (written last).

| #   | Module                                   | Core divergence                              |
| --- | ---------------------------------------- | -------------------------------------------- |
| 00  | Bootstrap                                | the falsifiability harness itself            |
| 01  | Interfaces & dynamic dispatch            | structural typing, itabs, existential types  |
| 02  | No sum types (+ exhaustiveness analyzer) | products vs coproducts; a `go/analysis` pass |
| 03  | Generics: type sets, GCShape             | partial monomorphization vs Rust/OCaml       |
| 04  | Errors                                   | `(T, error)` is a product, not a coproduct   |
| 05  | Memory, layout, allocations              | escape analysis, padding, GC knobs           |
| 06  | Memory model & data races                | happens-before; `Send`/`Sync`/`Sendable`     |
| 07  | Scheduler                                | GMP, preemption, `GOMAXPROCS` in cgroups     |
| 08  | API design & package boundaries          | encapsulation Go verifies vs `.mli`          |
| 09  | Reflection & code generation             | `reflect` cost vs `serde`/ppx                |
| 10  | Observability & performance              | pprof, trace, runtime/metrics                |
| 11  | Testing                                  | fuzzing, property tests, injected time       |
| 12  | Capstone                                 | one concurrent service in five languages     |
| 13  | Unsafe                                   | `unsafe.Pointer`, layout punning, FFI        |
| 14  | Type-system quirks (extension)           | what each type system can and cannot express |

## Theory stance

Modules frame divergences type-theoretically (products vs coproducts,
existential types for interfaces and `dyn`, System F and the higher-kinded gap,
GADTs as indexed families, the memory model as an axiomatic partial order) while
keeping every empirical claim tied to a runnable program. Concepts assumed:
ADTs, ownership/lifetimes, HM inference, functors and first-class modules,
async/executors, GC internals, memory models. Those are never re-explained; only
Go's position relative to them is.

## Honest limits

- The GitHub Actions workflow is authored to mirror the `make` targets but is
  not executed locally (no `act`); the `make` targets are the validated path.
- `docs-fresh` runs on macOS to have all five toolchains; a Linux-only mirror
  would need Swift-on-Linux and is not wired.
- Where a measurement cannot be made reliably (shared-runner timing, a
  non-reproducible behavior), the README says so rather than inventing a number.

## References

Every module ends with its own official-sources-first References section. The
language homes:

- Go: https://go.dev/doc/ (spec https://go.dev/ref/spec, memory model https://go.dev/ref/mem)
- Rust: https://doc.rust-lang.org/ (the Reference, the Book, the Rustonomicon)
- OCaml: https://ocaml.org/manual/
- Swift: https://docs.swift.org/swift-book/ and https://www.swift.org/documentation/
- Kotlin: https://kotlinlang.org/docs/home.html
