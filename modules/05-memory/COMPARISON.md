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

## Header layout across the five languages

The Go sizes measured in the README (`go-layout`) have direct analogues
elsewhere. Two ideas recur: the **fat pointer** (a pointer paired with metadata
in one 2-word value — Go's slice/interface, Rust's `&str`/`&dyn`, Swift's
existential) and the **per-object header** (a hidden word or two the runtime
prepends to every heap object — present in OCaml and the JVM, absent for value
types in Go/Rust/Swift). Sizes are 64-bit; 1 word = 8 bytes.

| Language | Value struct | Growable array | String / view | Dynamic / existential value | Per-heap-object header |
| -------- | ------------ | -------------- | ------------- | --------------------------- | ---------------------- |
| Go       | fields inline + padding (you order them); no header | `[]T` = 3-word header `(ptr,len,cap)`, shares its backing array | `string` = 2-word `(ptr,len)`, immutable | `interface` = 2-word `(itab, data)`; boxes a non-pointer payload | none on values; only GC/runtime metadata |
| Rust     | fields inline, compiler-reordered, no header | `Vec<T>` = 3 words `(ptr,len,cap)`, **owns** the buffer | `&str` / `&[T]` = 2-word fat pointer `(ptr,len)` | `&dyn Trait` / `Box<dyn>` = 2-word fat pointer `(data, vtable)`; enums are inline with **niche** packing | none (no GC, no object header) |
| OCaml    | uniform: a 1-word block header + 1 word per field, the whole thing boxed | boxed block (array) | boxed block, bytes packed (not 1 word/char) | first-class module `(module S)` packs an existential (boxed) | **1 word on every heap block** (size + tag byte + GC color); `int`/`char`/`bool` are unboxed **immediates** (63-bit, low bit = 1) |
| Swift    | value types inline, no header (ARC only for class fields) | `Array` = 1-word pointer to a heap buffer (copy-on-write) | `String` = 16 bytes; ≤15 UTF-8 bytes stored inline (SSO), else a heap buffer | `any P` = existential container: a 3-word inline value buffer + value-witness-table ptr + protocol-witness-table ptr (≈5 words); large values box | classes: `isa` pointer + ARC refcount |
| Kotlin/JVM | primitives unboxed inline; every object lives on the heap | `ArrayList` = an object wrapping a `T[]` | `String` = an object over a packed `byte[]` | an interface reference is one (compressed) pointer; dispatch goes through the klass vtable | **every object**: mark word + klass pointer ≈ 12–16 bytes (`+4` for an array length); references 4 bytes (compressed oops) or 8 |

The reads that matter for a systems engineer:

- **Go and Rust agree on the fat-pointer shapes** — `[]T`/`Vec<T>` are both a
  3-word `(ptr,len,cap)`; `string`/`&str` are both 2-word `(ptr,len)`;
  `interface`/`&dyn` are both 2-word `(data, table)`. The difference is
  ownership (Go's slice is a shared view with a GC behind it; Rust's `Vec` owns
  and frees its buffer) and that Rust puts **sum types inline** with niche
  packing where Go would use a boxed interface.
- **OCaml is the uniform-representation extreme**: every value is one word (an
  immediate or a pointer), and every heap object pays a one-word header. That
  uniformity is exactly what lets polymorphic code run without specialization
  (Module 03), and exactly what costs the boxing seen in the allocation
  benchmark — the two are the same design decision.
- **Swift** keeps value types header-free like Go/Rust, but its **existential
  container** is heavier than a 2-word interface (a 3-word inline buffer plus
  witness-table pointers), which is why `some P`/generics are preferred over
  `any P` on hot paths.
- **The JVM taxes every object** with a 12–16 byte header, so the difference
  between a `long` (8 bytes, unboxed) and a boxed `Long` (header + 8) is large
  and per-object — the reason `Int?` boxing and object churn dominate JVM memory
  behavior, and why escape analysis + scalar replacement matter so much there.

## References

Official sources first; see the module [`README`](README.md#references) for the
full per-language list (Go GC guide + GOGC, Rust ownership, OCaml GC, Swift ARC,
HotSpot GC tuning).
