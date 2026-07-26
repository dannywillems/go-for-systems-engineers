# 10 — Comparison: five observability toolchains

**Environment.** Apple M4 Pro (14 cores), macOS arm64. Go 1.26.5, Rust 1.92.0,
OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). Measured numbers in
[`measured.txt`](measured.txt).

## Where measurement lives

The sharpest split is whether the measurement stack is IN the standard
library/toolchain or a separate ecosystem you assemble.

| Concern | Go | Rust | OCaml | Swift | Kotlin/JVM |
| ------- | -- | ---- | ----- | ----- | ---------- |
| Microbenchmark | `testing.B` (stdlib) | Criterion (crate) | manual / `core_bench` | XCTest `measure` | JMH (separate) |
| Statistical rigor | `benchstat` (x/perf) | Criterion (built-in) | `core_bench` | XCTest baselines | JMH (built-in) |
| CPU profile | `runtime/pprof` (stdlib) | `perf` + cargo-flamegraph | `perf` / landmarks | Instruments | JFR / async-profiler |
| Alloc profile | `pprof` heap + `AllocsPerRun` | custom `GlobalAlloc` | `Gc` stats / memtrace | Instruments (Allocations) | JFR / `ThreadMXBean` |
| Live/prod profile | `net/http/pprof` | `tracing` + exporters | — | MetricKit (device) | JFR (always-on), Micrometer |
| Execution trace | `runtime/trace` | `tracing` spans | Eio/landmarks | Instruments (Time Profiler) | JFR events |
| In-process alloc count | **yes** (`AllocsPerRun`) | **yes** (global alloc) | **yes** (`Gc.minor_words`) | **no** (Instruments) | **yes** (`ThreadMXBean`) |

Go and the JVM put the whole stack in the platform: `testing.B`+`pprof`+`trace`
for Go, JMH+JFR for the JVM. Rust benchmarking rigor is excellent but lives in
Criterion (a crate) and the profile comes from OS `perf`. Swift's is the most
IDE-coupled: Instruments is the tool, and there is no in-process allocation
counter in the standard library — you profile allocations by attaching a
profiler, not by asking the runtime.

## The in-process allocation count

Four of five let a program count its own allocations deterministically, which is
what makes the falsifiable claim in this module possible:

- **Go**: `testing.AllocsPerRun(n, f)` sets GOMAXPROCS=1, forces a GC, and
  reports the average mallocs per call. Deterministic. 63 vs 1 here.
- **Rust**: no built-in, but `#[global_allocator]` lets you wrap `System` and
  bump an `AtomicUsize` per `alloc`. Deterministic once you isolate the region
  (and single-threaded, since the counter is process-global). 7 vs 1.
- **OCaml**: `Gc.minor_words ()` is a running total of minor-heap words; the
  delta across an operation is its allocation. 2440 vs 160 words.
- **Kotlin/JVM**: `com.sun.management.ThreadMXBean.getThreadAllocatedBytes(id)`
  reports per-thread allocated bytes; the delta is the operation's allocation.
  82328 vs 1352 bytes.
- **Swift**: none in the standard library. The number in `measured.txt` is time,
  not allocations; allocation profiling means Instruments' Allocations
  instrument or `MallocStackLogging`.

The counts are the machine-independent part; the ns/op alongside them is not, so
the counts drive the gated red/green tests and the timings live only in the
measured file.

## A semantic fact the counts expose

The allocation counts are not just tooling trivia — they reveal a real
difference in the string types. Go's `+=` costs O(n) allocations because Go
strings are immutable, so each concatenation copies. Rust's `String`, OCaml's
`Buffer`, and Swift's copy-on-write `String` are growable buffers with amortized
doubling, so naive concatenation costs O(log n) reallocations, which is why
Rust's count is 7 (not 63) and Swift's `+=` is nearly as fast as
`reserveCapacity`. The JVM `String +=` is O(n) again (its `String` is also
immutable; `StringBuilder` is the growable one), which is why Kotlin's naive byte
count is the largest. Pre-sizing collapses all of them to a single allocation.

## Bottom line

For "did my optimization actually reduce allocations", Go, Rust, OCaml, and the
JVM give you an in-process, deterministic yes/no; Swift sends you to Instruments.
For "how fast is it", everyone needs repeated runs and a statistics tool
(`benchstat`, Criterion, JMH), because a single wall-clock number is noise. The
discipline this module encodes — gate on deterministic counts, report timings
separately as non-portable measurements — is exactly what the five toolchains'
shapes recommend.
