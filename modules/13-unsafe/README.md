# 13 — Unsafe: raw memory, layout, and the FFI boundary

**Thesis.** Every one of these languages has an escape hatch out of its safety
guarantees, and they differ in how scoped, how checked, and how catchable the
resulting mistakes are. Go's `unsafe.Pointer` still type-checks and still runs
under the GC — but a `uintptr` is a plain number the GC does not track, which is
the one hazard that turns `unsafe` into use-after-free. The sanctioned zero-copy
conversions (`unsafe.String`/`unsafe.Slice`, Go 1.20+) keep the GC informed and
are measurably free.

## Contents

- [Zero-copy conversion, measured](#zero-copy-conversion-measured)
- [Layout without reflection](#layout-without-reflection)
- [The uintptr GC hazard](#the-uintptr-gc-hazard)
- [The same hatch in four languages](#the-same-hatch-in-four-languages)
- [References](#references)

## Zero-copy conversion, measured

`string(b)` copies; `unsafe.String(&b[0], len(b))` reinterprets in place:

<!-- BEGIN:snippet go-zerocopy -->
```go
// BytesToString reinterprets a []byte as a string WITHOUT copying: the string
// header points at the slice's backing array. Safe only if b is not mutated
// afterwards (strings are supposed to be immutable). string(b) copies; this
// does not.
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// StringToBytes reinterprets a string as a []byte WITHOUT copying. The result
// MUST NOT be written to (it aliases the string's immutable storage).
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
```
<!-- END:snippet go-zerocopy -->

<!-- BEGIN:file go-bench -->
```text
goos: darwin
goarch: arm64
pkg: github.com/dannywillems/go-for-systems-engineers/modules/13-unsafe/go
cpu: Apple M4 Pro
                  │     go      │
                  │   sec/op    │
StdStringCopy-14    11.26n ± 1%
ZeroCopyString-14   1.577n ± 0%
geomean             4.215n

                  │      go      │
                  │     B/op     │
StdStringCopy-14    80.00 ± 0%
ZeroCopyString-14   0.000 ± 0%
geomean                        ¹
¹ summaries must be >0 to compute geomean

                  │      go      │
                  │  allocs/op   │
StdStringCopy-14    1.000 ± 0%
ZeroCopyString-14   0.000 ± 0%
geomean                        ¹
¹ summaries must be >0 to compute geomean
```
<!-- END:file go-bench -->

The copy allocates 80 B and takes ~7x longer; the unsafe reinterpretation is one
pointer write and zero allocations. The catch is the invariant the type system
no longer enforces: the bytes must not be mutated while the string aliases them.

## Layout without reflection

`unsafe.Sizeof`/`Offsetof`/`Alignof` read a struct's layout at compile time, no
`reflect`:

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
Record: size=24 bytes, ID offset=8, align=8
BytesToString("hello") = "hello" (no copy)
StringToBytes("world") -> [119 111 114 108 100] ("world")
```
<!-- END:output go-demo -->

## The uintptr GC hazard

`unsafe.Pointer` is tracked by the GC; `uintptr` is not. Converting a pointer to
`uintptr`, storing it, and converting back is the classic bug: between the two
conversions the GC may free or (with a moving GC) relocate the object, leaving a
dangling number. The valid patterns are enumerated in the `unsafe` docs
(pointer arithmetic must happen in a single expression, `uintptr` never held
across a call). This is the one place Go's memory safety genuinely lapses, by
the programmer's explicit request.

## The same hatch in four languages

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
bytes_to_str(b"hello") = hello (unchecked, no copy)
f32_bits(1.0) = 0x3F800000
```
<!-- END:output rust-demo -->

<!-- BEGIN:output ocaml-demo -->
```text
$ dune exec bin/main.exe
bytes_to_string = hello (no copy)
Obj.magic launder 42 = 42
```
<!-- END:output ocaml-demo -->

<!-- BEGIN:output swift-demo -->
```text
floatBits(1.0) = 0x3F800000 (withUnsafeBytes)
bytesToString([104,105]) = hi
```
<!-- END:output swift-demo -->

<!-- BEGIN:output kotlin-demo -->
```text
floatBits(1.0f) = 0x3F800000 (ByteBuffer)
readIntLE([1,0,0,0]) = 1
```
<!-- END:output kotlin-demo -->

Rust's `unsafe` still runs the borrow checker; it only unlocks raw-pointer
deref, `transmute`, and `*_unchecked`, and Miri can detect the UB they enable.
OCaml's `Bytes.unsafe_to_string` is a sanctioned zero-copy, while `Obj.magic` is
the nuclear unchecked cast. Swift scopes it to `withUnsafeBytes`/`UnsafePointer`
closures. The JVM has no raw pointers at all: `sun.misc.Unsafe` is deprecated
and restricted, and the sanctioned modern path is the Foreign Function & Memory
API (`java.lang.foreign`); `ByteBuffer` is the portable structured-access tool.
See [`COMPARISON.md`](COMPARISON.md).

## References

Official sources first, grouped by language.

### Go

- `unsafe` package (the valid `Pointer`/`uintptr` patterns): https://pkg.go.dev/unsafe
- `unsafe.String` / `unsafe.Slice` (Go 1.20): https://pkg.go.dev/unsafe#String
- Go Memory Model (why uintptr is not tracked): https://go.dev/ref/mem

### Rust

- The Rustonomicon (unsafe Rust): https://doc.rust-lang.org/nomicon/
- `std::mem::transmute`: https://doc.rust-lang.org/std/mem/fn.transmute.html
- Miri (UB detection): https://github.com/rust-lang/miri

### OCaml

- `Bytes.unsafe_to_string`: https://ocaml.org/manual/5.4/api/Bytes.html
- `Obj` module (the unchecked escape hatch): https://ocaml.org/manual/5.4/api/Obj.html
- `Bigarray`: https://ocaml.org/manual/5.4/api/Bigarray.html

### Swift

- `withUnsafeBytes` / `UnsafeRawPointer`: https://developer.apple.com/documentation/swift/unsaferawpointer

### Kotlin (JVM)

- Foreign Function & Memory API (JEP 454): https://openjdk.org/jeps/454
- `java.nio.ByteBuffer`: https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/nio/ByteBuffer.html
