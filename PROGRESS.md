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
| 03  | Generics / GCShape            | done    | GCShape objdump (pointer collapse to *uint8), stencil bench (generic≈concrete≪interface), HKT reject, honest Go-vs-Rust binary size, functors/witness/erasure |
| 04  | Errors                        | done    | (T,error) product-not-sum (io.EOF), typed nil, Result monad laws + unusable, errcheck + generic-method rejects |
| 05  | Memory, layout, allocations   | done    | padding/fieldalignment, escape analysis (-m), slice aliasing, GOGC throughput-vs-RSS, 5-language alloc bench (OCaml minor GC beats Go; GCs beat malloc here) |
| 06  | Memory model & data races     | done    | -race capture, dynamic-detector limits, goleak, errgroup; Rust Send/Sync + Swift 6 Sendable rejects; OCaml domains/Atomic; atomic-vs-mutex bench; References |
| 07  | Scheduler                     | done    | GMP + async preemption (100k goroutines), p50/p99 latency vs worker count, GOMAXPROCS/cgroup (Go 1.25 fix), 5-language CPU throughput (all within ~15%) |
| 08  | API design & boundaries       | done    | compile-verified encapsulation: opaque type in 5 langs + 5 compile-rejects (Go internal/, Rust E0616, OCaml abstract .mli, Swift private, Kotlin private); representation-hiding ladder (OCaml .mli abstract type strongest); Ratio-invariant exercise |
| 09  | Reflection & codegen          | done    | runtime reflection (Go reflect/encoding/json) vs compile-time codegen (Rust derive, Swift synthesis, Kotlin data class, OCaml ppx/manual); 5 deterministic demos + 4 compile-rejects; error-timing is the real divergence (go vet catches bad tag; others reject at compile); measured reflection cost (~6x/3x allocs, amortized by cached per-type encoder) |
| 10  | Observability & performance   | done    | deterministic alloc COUNT as the falsifiable spine (AllocsPerRun: 63 vs 1); pprof (hotLoop 99.6%) + runtime/metrics + benchstat; in-process alloc counters in Go/Rust(global alloc)/OCaml(minor_words)/JVM(ThreadMXBean), Swift the exception (Instruments/os_signpost) |
| 11  | Testing                       | done    | table + property (testing/quick) + native fuzz + testing/synctest virtual time (Go 1.25, rate-limiter tested deterministically); reject-buggy shows the fuzz invariants CATCHING a planted bug; hand-rolled property loops in the other 4; Dedup exercise |
| 12  | Capstone                      | done    | concurrent bounded cache (mutex+eviction+semaphore backpressure+ctx shutdown) in 5 langs; throughput+p50/p99/p999 bench + migration matrix (binary/LOC/cold-compile); actor trades throughput for a tight tail |
| 13  | Unsafe                        | done    | unsafe.String/Slice zero-copy (1.6ns/0 allocs vs 11ns/80B), Offsetof layout, uintptr GC hazard; Rust/Miri, OCaml Obj.magic, Swift withUnsafeBytes, JVM FFM/ByteBuffer |
| 14  | Type-system quirks            | done    | Go phantom type params, Rust typestate, OCaml GADT eval + type-equality witness, Rust const-generic length-indexing (+reject), Swift typestate, Kotlin variance; 5 compile-rejects; PROPOSITIONS.md = 40-proposition catalogue (Curry-Howard, propositional core -> quantifiers -> indexed -> substructural -> dependent frontier) with a per-proposition 5-language +/-/x matrix (every version-sensitive cell compiler-checked) |

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
