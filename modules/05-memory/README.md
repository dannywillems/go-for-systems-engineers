# 05 — Memory, layout, allocations, and GC

**Thesis.** Go gives you value types with predictable layout (like Rust) but a
tracing garbage collector (unlike Rust, like OCaml/JVM). Three things bite an
engineer coming from either side: struct **padding** wastes memory silently,
**escape analysis** decides stack-vs-heap and is easy to defeat by accident, and
the **GC** trades throughput against memory (RSS) on a dial (`GOGC`). The
cross-language allocation numbers hold two surprises, measured below.

## Contents

- [Struct layout and padding](#struct-layout-and-padding)
- [Escape analysis](#escape-analysis)
- [Slice aliasing](#slice-aliasing)
- [GC: throughput vs RSS, and the cross-language surprise](#gc-throughput-vs-rss-and-the-cross-language-surprise)
- [References](#references)

## Struct layout and padding

Field order changes size. The same three fields are 24 or 16 bytes depending on
order; alignment inserts padding:

<!-- BEGIN:snippet go-padding -->
```go
// Padded orders fields worst-case: a bool, then an 8-byte int64, then a bool.
// Alignment forces padding, so the struct is 24 bytes even though its data is
// 10. fieldalignment would rewrite it to the Packed order below.
//
//nolint:govet // intentionally field-misaligned to measure padding
type Padded struct {
	A bool
	B int64
	C bool
}

// Packed puts the 8-byte field first, then the two bools pack together: 16
// bytes. Same fields, 33% less memory — multiplied across a large slice.
type Packed struct {
	B int64
	A bool
	C bool
}
```
<!-- END:snippet go-padding -->

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
sizeof(Padded) = 24 bytes
sizeof(Packed) = 16 bytes
saved per element = 8 bytes (7 MiB per 1M elements)
slice aliasing: orig was [1 2 3 4 5], after append(orig[:2], 99) it is [1 2 99 4 5]
```
<!-- END:output go-demo -->

The `fieldalignment` analyzer (in `golangci-lint`) flags the wasteful order;
across a large slice the 8 bytes/element is real memory. Rust's `repr(Rust)`
reorders fields automatically, so this footgun does not exist there; Go keeps
declaration order and leaves the ordering to you.

## Escape analysis

Whether a value lives on the stack or the heap is decided by escape analysis,
visible with `-gcflags=-m`:

<!-- BEGIN:snippet go-escape -->
```go
// Sink is an escape target for the boxing case below. It is exported so
// staticcheck does not flag it as write-only/unused.
var Sink any

// NoEscape: the array is used only locally, so escape analysis keeps it on the
// STACK (no allocation), like Rust's default.
func NoEscape() int {
	v := [4]int{1, 2, 3, 4}
	return v[0] + v[3]
}

// EscapesReturn: returning a pointer to a local forces the local to the HEAP —
// the value must outlive the frame.
func EscapesReturn() *int {
	x := 42
	return &x
}

// EscapesInterface: putting a concrete value into an interface boxes it on the
// HEAP, even though the value itself is tiny. A frequent, surprising source of
// allocation.
func EscapesInterface(n int) {
	Sink = n
}
```
<!-- END:snippet go-escape -->

<!-- BEGIN:output go-escape -->
```text
./escape.go:19:2: moved to heap: x
./escape.go:27:9: n escapes to heap
./escape.go:35:15: []int{...} escapes to heap
./escape.go:36:17: append escapes to heap
./escape.go:38:12: append does not escape
```
<!-- END:output go-escape -->

A local array stays on the stack; returning a pointer to a local, or boxing a
value into an `interface`, moves it to the heap. That interface-boxing line is
the surprising one: putting a tiny `int` into `any` allocates.

## Slice aliasing

The `go-demo` output above also shows the aliasing trap: a subslice shares the
parent's backing array, so `append` into it with spare capacity overwrites the
parent (`[1 2 3 4 5]` becomes `[1 2 99 4 5]`). The full-slice expression
`s[lo:hi:hi]` caps capacity to force a fresh array; `slices.Clone` copies.

## GC: throughput vs RSS, and the cross-language surprise

`cmd/gcdemo` runs 50M short-lived heap allocations and reports time, `NumGC`,
and `HeapSys`. Run under two `GOGC` values it shows the classic trade, then the
same workload across five languages (single run, optimized, non-deterministic —
regenerated manually by `scripts/mem-bench.sh`):

<!-- BEGIN:file measured -->
```text
machine: Apple M4 Pro (14 cores), macOS arm64
workload: 50,000,000 short-lived ~8-word heap objects, single run

# Go: GOGC throughput vs RSS trade (same workload, two GOGC values)
GOGC=100  alloc 50M: 489 ms  NumGC=891  HeapSys=7 MiB  (acc=1249999975000000)
GOGC=800  alloc 50M: 411 ms  NumGC=102  HeapSys=35 MiB  (acc=1249999975000000)

# Cross-language allocation time (optimized builds; see caveats)
Go     alloc 50M: 499 ms  NumGC=892  HeapSys=7 MiB  (acc=1249999975000000)
OCaml alloc 50M (minor GC): 432 ms (acc=1250000025000000)
Rust alloc 50M (Box, no GC): 566 ms (acc=1249999975000000)
Swift alloc 50M (class/ARC): 916 ms (acc=1249999975000000)
Kotlin alloc 50M (JVM GC): 293 ms (acc=1249999975000000)
```
<!-- END:file measured -->

Two results worth stating:

1. **`GOGC` is a throughput/RSS dial.** Raising it from 100 to 800 cuts GC
   cycles roughly 9x and speeds the workload up, at ~5x the `HeapSys`. There is
   no free lunch; there is a knob (and `GOMEMLIMIT` as a soft cap).
2. **Tracing GCs beat manual `malloc`/`free` here, and OCaml's minor GC beats
   Go's.** For this high-churn, short-lived workload the JVM and OCaml (bump-
   pointer allocation + generational collection of young garbage) are fastest,
   Go's collector is in the middle, and Rust's per-object `Box` (a `malloc` and
   `free` each) is *slower* — allocation-heavy short-lived workloads are exactly
   where a good GC wins. Swift's ARC is slowest, paying atomic refcount traffic
   per reference. This is workload-specific (long-lived or large objects flip
   the result), which is precisely why it is measured, not assumed.

See [`COMPARISON.md`](COMPARISON.md) for the five memory-management strategies
side by side and [`exercises/`](exercises).

## References

Official sources first, grouped by language.

### Go

- The Go Blog, "Getting to Go: memory management and the GC": https://go.dev/blog/ismmkeynote
- `runtime` (GOGC, GOMEMLIMIT): https://pkg.go.dev/runtime#hdr-Environment_Variables
- A Guide to the Go Garbage Collector: https://go.dev/doc/gc-guide
- Escape analysis / `-gcflags=-m`: https://go.dev/doc/faq#stack_or_heap
- `unsafe.Sizeof` and alignment: https://pkg.go.dev/unsafe#Sizeof
- `slices` package (`Clone`): https://pkg.go.dev/slices#Clone

### Rust

- The Rust Book, "What is Ownership?" (stack vs heap, no GC): https://doc.rust-lang.org/book/ch04-01-what-is-ownership.html

### OCaml

- OCaml Manual, "Understanding the garbage collector" (minor/major heaps): https://ocaml.org/docs/garbage-collection

### Swift

- Swift, Automatic Reference Counting (ARC): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/automaticreferencecounting/

### Kotlin (JVM)

- HotSpot GC tuning guide: https://docs.oracle.com/en/java/javase/21/gctuning/introduction-garbage-collection-tuning.html
