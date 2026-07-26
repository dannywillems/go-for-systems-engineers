# 00 — Comparison: toolchain and static-analysis landscape

**Measurement environment.** Apple M4 Pro (14 cores), macOS (Darwin 25.5.0,
arm64). Toolchains: Go 1.26.5 (fetched via `GOTOOLCHAIN`; the installed driver
may be older), Rust 1.92.0 (edition 2024), OCaml 5.4.0 + dune 3.20.2. CI mirrors
these via the version variables in `.github/workflows/ci.yml`. Exact versions
this run drove:

<!-- BEGIN:output toolchain-versions -->
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
<!-- END:output toolchain-versions -->

## The three fixtures, side by side

Extracted from source, never retyped:

<!-- BEGIN:snippet go-sum -->
```go
// Sum returns 1 + 2 + ... + n. The value is identical on every 64-bit target
// and in every language, which makes it a clean cross-toolchain fixture.
func Sum(n int) int {
	total := 0
	for i := 1; i <= n; i++ {
		total += i
	}
	return total
}

// WordSizeBytes is sizeof(int) on this target: 8 on any 64-bit platform,
// so it is stable across the author's arm64 and a CI amd64 runner.
func WordSizeBytes() int {
	var x int
	return int(unsafe.Sizeof(x))
}
```
<!-- END:snippet go-sum -->

<!-- BEGIN:snippet rust-sum -->
```rust
/// `sum(n)` returns 1 + 2 + ... + n. The result is identical on every 64-bit
/// target and in every language, which makes it a clean cross-toolchain fixture.
pub fn sum(n: u64) -> u64 {
    (1..=n).sum()
}

/// Size of a pointer-width integer in bytes: 8 on any 64-bit platform.
pub fn word_size_bytes() -> usize {
    std::mem::size_of::<usize>()
}
```
<!-- END:snippet rust-sum -->

<!-- BEGIN:snippet ocaml-sum -->
```ocaml
(** [sum n] returns 1 + 2 + ... + n. Identical on every 64-bit target and in
    every language, which makes it a clean cross-toolchain fixture. *)
let sum n =
  let total = ref 0 in
  for i = 1 to n do
    total := !total + i
  done;
  !total

(** Native word size in bytes: 8 on any 64-bit platform. Note OCaml's [int] is
    63-bit (one tag bit), but the machine word it lives in is still 8 bytes. *)
let word_size_bytes () = Sys.word_size / 8
```
<!-- END:snippet ocaml-sum -->

<!-- BEGIN:snippet swift-sum -->
```swift
/// `sum(n)` returns 1 + 2 + ... + n. Swift's `Int` is 64-bit on a 64-bit
/// target, so the value fits and matches the Go/Rust/OCaml/Kotlin sides.
public func sum(_ n: Int) -> Int {
  if n < 1 { return 0 }
  var total = 0
  for i in 1...n { total += i }
  return total
}

/// Size of the native word in bytes: 8 on any 64-bit platform.
public func wordSizeBytes() -> Int {
  MemoryLayout<Int>.size
}
```
<!-- END:snippet swift-sum -->

<!-- BEGIN:snippet kotlin-sum -->
```kotlin
/**
 * Returns 1 + 2 + ... + n. On the JVM `Int` is 32-bit, which overflows here, so
 * the accumulator is a 64-bit `Long`. This 32-bit-Int fact is itself a
 * representational divergence from Go/Rust/Swift (64-bit Int) worth noting.
 */
fun sum(n: Int): Long {
    var total = 0L
    for (i in 1..n) total += i
    return total
}

/**
 * Native word size in bytes: 8 on a 64-bit JVM. There is no `sizeof` on the
 * JVM, so this reads the data model the JVM runs under ("64" -> 8 bytes).
 */
fun wordSizeBytes(): Int = System.getProperty("sun.arch.data.model").toInt() / 8
```
<!-- END:snippet kotlin-sum -->

## Integer models: same value, five representations

The fixture returns one number, but the type that carries it differs, which
already teaches something about each runtime's memory model:

| Language    | `Int` width | Boxing / tagging                                            |
| ----------- | ----------- | ---------------------------------------------------------- |
| Go          | 64-bit      | unboxed machine word                                       |
| Rust        | 64-bit (`u64` here) | unboxed; width is explicit in the type              |
| Swift       | 64-bit      | unboxed value type                                         |
| OCaml       | 63-bit      | 1 tag bit distinguishes ints from pointers (no heap box)   |
| Kotlin/JVM  | 32-bit      | primitive `int`; needs `Long` (64-bit) to hold this sum    |

OCaml's tag bit is why `int` is 63-bit: the GC must tell an immediate integer
from a pointer in a uniform machine word, so the low bit is reserved. Kotlin's
32-bit `Int` is a JVM inheritance; the accumulator overflows unless promoted to
`Long`. These are not trivia: they surface again in Module 05 (layout) and
Module 06 (what a torn read of a wide value can even mean).

## Static-analysis maturity: five philosophies

The ecosystems put correctness checking in very different places.

| Concern            | Go                                    | Rust                          | OCaml 5                          | Swift                          | Kotlin                         |
| ------------------ | ------------------------------------- | ----------------------------- | -------------------------------- | ------------------------------ | ------------------------------ |
| Formatter          | `gofmt` (canonical, bundled)          | `rustfmt`                     | `ocamlformat`                    | `swift format` (bundled)       | `ktlint`                       |
| Core vet           | `go vet` (bundled)                    | `rustc` lints                 | compiler warnings                | compiler diagnostics           | `kotlinc` warnings             |
| Deep static linter | `staticcheck`, `golangci-lint` (~50)  | `clippy` (~600 lints)         | none standard (type system)      | `swift format lint`; SwiftLint | `ktlint`; `detekt`             |
| Vuln scanner       | `govulncheck` (reachability)          | `cargo audit`/`deny`          | `opam audit` (less mature)       | limited                        | JVM: OWASP dep-check, Snyk     |
| Data-race checker  | `-race` (dynamic, TSan)               | *static*: `Send`/`Sync`       | domains new; dynamic tooling nascent | *static*: `Sendable` (Swift 6) | JVM: dynamic; no static guarantee |
| Memory management  | tracing GC                            | none (ownership/RAII)         | tracing GC (gen., bump minor)    | ARC (ref counting)             | tracing GC (JVM)               |
| Benchmark harness  | `testing.B`+`benchstat` (bundled)     | `criterion`                   | `core_bench`                     | XCTest `measure` / package-benchmark | JMH (JVM)              |

The theoretical reading: **Rust and Swift 6 push the most concurrency invariants
into the type system** (`Send`/`Sync`; `Sendable`), so a class of data-race bugs
is *rejected at compile time*. **Go and Kotlin/JVM push little** and rely on
dynamic detection (`-race`; JVM tooling) and discipline. OCaml sits between. On
memory, the five split three ways: no GC (Rust), reference counting (Swift ARC,
with cycle leaks as the cost), and tracing GC (Go, OCaml, JVM). That memory axis
is what Modules 05-07 measure directly, in numbers rather than adjectives.

## Verifying the harness itself

`make docs` regenerates every block above; `make docs-check` (the `docs-fresh`
gate) fails if any portable block drifts. Run `make module M=00` to build, test,
lint, and check-docs this module in isolation.
