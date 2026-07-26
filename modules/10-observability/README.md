# 10 — Observability & performance

**Thesis.** Go ships a first-class measurement stack in the standard library and
toolchain, so profiling and benchmarking are not third-party concerns:
`testing.B` + `-benchmem` for microbenchmarks, `runtime/pprof` + `go tool pprof`
for CPU/heap/block/mutex profiles, `runtime/trace` for the execution tracer, and
the versioned `runtime/metrics` API for programmatic metrics. The falsifiable
spine of this module is the allocation **count**: it is deterministic
(`testing.AllocsPerRun` forces a GC and counts mallocs), so an optimization's
allocation reduction is a *fact* a program proves, whereas wall-clock time is
noisy and belongs in a measured file read through `benchstat` for significance.

## Contents

- [The measured claim: allocations are deterministic](#the-measured-claim-allocations-are-deterministic)
- [The two builders](#the-two-builders)
- [Profiling: runtime/pprof](#profiling-runtimepprof)
- [The stable metrics API: runtime/metrics](#the-stable-metrics-api-runtimemetrics)
- [Counting allocations in the other four languages](#counting-allocations-in-the-other-four-languages)
- [Measured results](#measured-results)
- [Exercises](#exercises)
- [References](#references)

## The measured claim: allocations are deterministic

Two implementations of the same function build a 64-part string: one with the
`+=` operator, one with a pre-sized `strings.Builder`. The number of allocations
each performs is not a matter of opinion — a program prints it, and the count is
machine-independent:

<!-- BEGIN:output go-allocs -->
```text
$ go run ./cmd/allocs
ConcatPlus  (64 parts): 63 allocs/op
BuilderGrow (64 parts): 1 allocs/op
reduction: 63x fewer allocations
```
<!-- END:output go-allocs -->

63 versus 1. Because Go strings are immutable, every `+=` allocates a fresh
backing array and copies everything so far (quadratic copying, O(n)
allocations); the pre-sized `Builder` allocates its buffer once and `String()`
returns it without a copy. `TestBuilderAllocatesLess` turns this into a
red/green gate (`grow == 1 && grow < plus`), so the claim cannot silently rot.

## The two builders

<!-- BEGIN:snippet go-builder -->
```go
// ConcatPlus builds a string with the += operator in a loop. Strings are
// immutable, so each += allocates a NEW backing array and copies everything so
// far: quadratic copying and O(n) allocations.
func ConcatPlus(parts []string) string {
	s := ""
	for _, p := range parts {
		s += p
	}
	return s
}

// BuilderGrow builds the same string with strings.Builder, pre-sizing the
// backing array once with Grow so the whole build costs a SINGLE allocation.
func BuilderGrow(parts []string) string {
	total := 0
	for _, p := range parts {
		total += len(p)
	}
	var b strings.Builder
	b.Grow(total)
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}
```
<!-- END:snippet go-builder -->

## Profiling: runtime/pprof

A CPU profile is a standard-library artifact. `runtime/pprof` writes a
pprof-format profile that `go tool pprof` analyzes offline (`top`, `list`,
`web`, a flame graph). The `measured.txt` block below shows the profiler
correctly attributing ~99% of CPU to the hot loop.

<!-- BEGIN:snippet go-pprof -->
```go
// CPUProfile runs work while a CPU profile is being collected, writing the
// pprof-format profile to w. The profile is analyzed offline with
// `go tool pprof` (top, list, web) — profiling is a standard-library concern,
// not a third-party one.
func CPUProfile(w io.Writer, work func()) error {
	if err := pprof.StartCPUProfile(w); err != nil {
		return err
	}
	defer pprof.StopCPUProfile()
	work()
	return nil
}
```
<!-- END:snippet go-pprof -->

The same package also exposes `net/http/pprof` (live profiles from a running
server), block and mutex profiles (`runtime.SetBlockProfileRate`,
`SetMutexProfileFraction`), and `runtime/trace` for the execution tracer, which
visualizes goroutine scheduling, GC pauses, and syscalls on a timeline.

## The stable metrics API: runtime/metrics

`runtime.MemStats` is legacy; `runtime/metrics` is the versioned, self-describing
replacement. Every metric has a documented name (`/gc/heap/allocs:bytes`,
`/sched/goroutines:goroutines`) and a kind, so it is safe to consume
programmatically:

<!-- BEGIN:snippet go-metrics -->
```go
// ReadMetric reads one sample from the stable runtime/metrics API. Unlike the
// legacy runtime.MemStats, this API is versioned and self-describing: every
// metric has a documented name like "/gc/heap/allocs:bytes" and a kind.
func ReadMetric(name string) uint64 {
	sample := []metrics.Sample{{Name: name}}
	metrics.Read(sample)
	if sample[0].Value.Kind() == metrics.KindUint64 {
		return sample[0].Value.Uint64()
	}
	return 0
}
```
<!-- END:snippet go-metrics -->

## Counting allocations in the other four languages

The falsifiable technique — count allocations in-process, deterministically —
ports to three of the four; the fourth is the interesting exception:

- **Rust** has no `AllocsPerRun`, but it lets you install a `#[global_allocator]`,
  so a counting wrapper over the system allocator gives exact counts:

  <!-- BEGIN:snippet rust-counter -->
  <!-- END:snippet rust-counter -->

- **OCaml** reads `Gc.minor_words ()` before and after; the difference is exactly
  what the operation allocated.
- **Kotlin/JVM** reads `com.sun.management.ThreadMXBean.getThreadAllocatedBytes`,
  the bytes a thread has allocated.
- **Swift** is the exception: its standard library has no in-process allocation
  counter. The idiomatic path is `os_signpost` + Instruments — you bracket a
  region so the profiler shows its duration and allocations on a timeline, an
  *external* tool rather than an in-process number:

  <!-- BEGIN:snippet swift-signpost -->
  <!-- END:snippet swift-signpost -->

A second cross-language fact falls out of the counts: Go's `+=` is O(n)
allocations because its strings are immutable, but Rust's `String`, OCaml's
`Buffer`, Kotlin's `StringBuilder`, and Swift's copy-on-write `String` are
growable, so naive concatenation there costs only O(log n) reallocations (or,
for the JVM `String +=`, O(n) again — the JVM string is immutable too). See
[`COMPARISON.md`](COMPARISON.md).

## Measured results

Regenerated locally by [`scripts/obs-bench.sh`](../../scripts/obs-bench.sh)
(non-portable: allocation *counts* are deterministic but ns/op, profile
percentages, and metric values are machine-specific, so this block is skipped by
the docs-freshness gate).

<!-- BEGIN:file measured -->
```text
machine: Apple M4 Pro (14 cores), macOS arm64
workload: build a 64-part string two ways (naive concat vs pre-sized),
1,000,000 iterations, optimized builds.

# Allocation cost per build (counts are deterministic; ns/op is not)
Rust concat_plus: 251 ns/op (7 allocs)  with_cap: 122 ns/op (1 allocs)
OCaml concat_caret: 874 ns/op (2440 words)  buffer_build: 345 ns/op (160 words)
Kotlin concatPlus: 782 ns/op (82328 B)  builder: 237 ns/op (1352 B)
Swift concatPlus: 1977 ns/op  reserve: 1803 ns/op

# Go microbenchmark (go test -bench -benchmem): allocs/op deterministic
BenchmarkConcatPlus-14     	  514694	      2330 ns/op	   19592 B/op	      63 allocs/op
BenchmarkBuilderGrow-14    	 5931955	       200.2 ns/op	     576 B/op	       1 allocs/op

# Go CPU profile: hottest function (go tool pprof -top)
     2.61s 98.86% 98.86%      2.63s 99.62%  observability.hotLoop

# Go runtime/metrics readout (names are a stable API; values are not)
/gc/heap/allocs:bytes              = 4755656
/gc/heap/objects:objects           = 530
/sched/goroutines:goroutines       = 20
/memory/classes/total:bytes        = 17647880
```
<!-- END:file measured -->

Reading it: allocation counts match across the deterministic tools (Go 63/1,
Rust 7/1, OCaml 2440/160 words, Kotlin 82328/1352 bytes), and the ns/op tracks
the allocation cost — the pre-sized build is several times faster everywhere,
except Swift, whose copy-on-write `+=` is already close to `reserveCapacity`, so
pre-sizing buys little. The `benchstat` tool (golang.org/x/perf) turns repeated
`-benchmem` runs into a mean with a confidence interval; use it rather than
eyeballing a single ns/op, which is noise.

## Exercises

[`exercises/go`](exercises/go) is red until you implement an allocation-free
join and a percentile over a benchmark's samples. [`solutions/go`](solutions/go)
is the verified corrigé:

```
make exercises M=10   # red
make solutions M=10   # green
```

## References

Official sources first, grouped by language.

### Go

- `testing` (`B`, `AllocsPerRun`, `B.Loop`, `ReportAllocs`): https://pkg.go.dev/testing
- Profiling Go programs (the pprof blog post): https://go.dev/blog/pprof
- `runtime/pprof`: https://pkg.go.dev/runtime/pprof
- `net/http/pprof`: https://pkg.go.dev/net/http/pprof
- `runtime/metrics`: https://pkg.go.dev/runtime/metrics
- `runtime/trace` + the execution tracer: https://pkg.go.dev/runtime/trace
- `benchstat` (golang.org/x/perf): https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
- Diagnostics guide: https://go.dev/doc/diagnostics

### Rust

- `std::alloc::GlobalAlloc` (custom global allocator): https://doc.rust-lang.org/std/alloc/trait.GlobalAlloc.html
- Criterion.rs (statistics-driven benchmarking): https://bheisler.github.io/criterion.rs/book/
- `tracing` (structured, async-aware instrumentation): https://docs.rs/tracing/
- cargo-flamegraph: https://github.com/flamegraph-rs/flamegraph

### OCaml

- `Gc` (minor_words, allocated_bytes, stat): https://ocaml.org/manual/5.4/api/Gc.html
- Memtrace (allocation profiler): https://github.com/janestreet/memtrace
- Landmarks (a profiling library): https://github.com/LexiFi/landmarks

### Swift

- `os_signpost` / Instruments (points of interest): https://developer.apple.com/documentation/os/logging/recording_performance_data
- swift-metrics: https://github.com/apple/swift-metrics
- XCTest performance tests (`measure`): https://developer.apple.com/documentation/xctest/performance_tests

### Kotlin (JVM)

- `com.sun.management.ThreadMXBean.getThreadAllocatedBytes`: https://docs.oracle.com/en/java/javase/21/docs/api/jdk.management/com/sun/management/ThreadMXBean.html
- JMH (the JVM microbenchmark harness): https://github.com/openjdk/jmh
- JDK Flight Recorder (JFR): https://docs.oracle.com/en/java/javase/21/jfapi/
- async-profiler: https://github.com/async-profiler/async-profiler
