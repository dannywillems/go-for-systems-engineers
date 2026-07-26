# Progress

Status of each module. Updated after each module lands. Languages per module:
Go (subject), Rust, OCaml, Swift, Kotlin.

Toolchains: Go 1.26.5 (via GOTOOLCHAIN), Rust 1.92.0 (edition 2024), OCaml 5.4.0
+ dune 3.20.2, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). Machine for benchmarks:
Apple M4 Pro (14 cores), macOS arm64.

## Legend

- done: builds, tests, lints, docs regenerate and pass `make module M=NN`.
- partial: some languages or sections pending.
- pending: not started.

## Modules

| #   | Module                        | Status  | Notes                                             |
| --- | ----------------------------- | ------- | ------------------------------------------------- |
| 00  | Bootstrap                     | done    | capture engine, Makefile, CI, 5-lang build, bench harness |
| 01  | Interfaces & dispatch         | done    | itab/existentials, devirt (-m), dispatch bench, Rust orphan-rule reject |
| 02  | No sum types + analyzer       | done    | go/analysis exhaustiveness pass (+analysistest, CI), 4-language reject captures, Option/niche |
| 03  | Generics / GCShape            | pending |                                                   |
| 04  | Errors                        | pending |                                                   |
| 05  | Memory, layout, allocations   | pending |                                                   |
| 06  | Memory model & data races     | pending |                                                   |
| 07  | Scheduler                     | pending |                                                   |
| 08  | API design & boundaries       | pending |                                                   |
| 09  | Reflection & codegen          | pending |                                                   |
| 10  | Observability & performance   | pending |                                                   |
| 11  | Testing                       | pending |                                                   |
| 12  | Capstone                      | pending |                                                   |
| 13  | Unsafe                        | pending |                                                   |
| 14  | Type-system quirks            | pending |                                                   |

## Build order

Priority after 00: the semantic-divergence core (01, 02, 03, 04, 06), then the
concurrency/GC/runtime measurement modules (05, 07, 10) which carry the
cross-language benchmark weight, then 08/09/11, the two extension modules
(13, 14), and the capstone (12) last.

## Toolchain notes

- Go 1.26.5 is fetched by `GOTOOLCHAIN=auto` from a `go 1.26.5` directive in
  every `go.mod`, even though the installed driver may be older.
- OCaml uses a dedicated opam switch (`OCAML_SWITCH=5.4.0`); CI overrides the
  switch via `OPAM="opam exec --"`.
- Kotlin is built with `kotlinc` directly (no Gradle) to keep the repo
  self-contained; coroutine-dependent modules (06/07) will need a classpath for
  `kotlinx-coroutines` and that is noted where it lands.
- Swift builds with SwiftPM; Swift 6 strict concurrency (`Sendable`) is a
  first-class comparison point in Module 06.
