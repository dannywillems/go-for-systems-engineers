# 05 — Comparison: five memory-management strategies

**Environment.** Apple M4 Pro (14 cores), macOS arm64. Go 1.26.5, Rust 1.92.0,
OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). The measured allocation
timings are in [`measured.txt`](measured.txt), injected into the module
[`README`](README.md#gc-throughput-vs-rss-and-the-cross-language-surprise).

## The strategy axis

| Language | Reclamation | Allocation | Layout control | Cost model |
| -------- | ----------- | ---------- | -------------- | ---------- |
| Go       | tracing GC (non-generational, concurrent mark) | size-classed heap | declaration order (padding is yours) | GC pauses; `GOGC`/`GOMEMLIMIT` dial |
| Rust     | none (ownership/RAII) | allocator `malloc`/`free` per object | auto field reorder (`repr(Rust)`) | deterministic free; no pause, no amortization |
| OCaml    | generational tracing GC | bump-pointer minor heap | boxed by default; unboxed floats/records opt-in | very cheap young allocation + collection |
| Swift    | ARC (reference counting) | allocator per object | value types stack/inline; classes heap | atomic refcount traffic; cycle leaks |
| Kotlin   | JVM tracing GC (generational, TLAB) | thread-local bump (TLAB) | erased/boxed; JIT scalar replacement | GC pauses; JIT can elide allocation |

## What the numbers said (see `measured.txt`)

For 50M short-lived ~8-word objects, the tracing GCs (Kotlin/JVM, OCaml) were
the *fastest*, Go's collector was in the middle, Rust's per-object `Box` was
*slower*, and Swift's ARC was slowest. The intuition "no GC is always faster"
is wrong for allocation-heavy, short-lived workloads: bump-pointer allocation
plus generational collection of young garbage beats a `malloc`/`free` pair per
object, and beats atomic refcounting per reference. This is not a verdict on the
languages — long-lived data, large objects, or low allocation rates reverse it,
and Rust's model wins decisively when you can avoid heap allocation entirely
(which its ownership system makes easy and the others do not). It is a
measured caution against assuming the GC is the bottleneck.

The other axis is **control**: Rust and Go expose value types with predictable
layout (Rust even reorders fields for you); OCaml and the JVM box by default and
recover unboxing through specific features (`[@@unboxed]`, JIT scalar
replacement). The `GOGC` result in `measured.txt` is the one knob every Go
service operator eventually turns: throughput and tail latency versus RSS.

## References

Official sources first; see the module [`README`](README.md#references) for the
full per-language list (Go GC guide + GOGC, Rust ownership, OCaml GC, Swift ARC,
HotSpot GC tuning).
