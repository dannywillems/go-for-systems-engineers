# 13 — Comparison: escape hatches

**Environment.** Go 1.26.5, Rust 1.92.0, OCaml 5.4.0, Swift 6.2.3, Kotlin
2.4.10 (OpenJDK 26). Go benchmark in [`go/bench.txt`](go/bench.txt).

## How each language lets you break the rules

| Language | Hatch | Still checked | Failure caught by | Raw pointers? |
| -------- | ----- | ------------- | ----------------- | ------------- |
| Go       | `unsafe.Pointer`, `uintptr` | type checker yes; GC yes (but not `uintptr`) | nothing (uintptr-held pointer is UB) | yes (`unsafe.Pointer`) |
| Rust     | `unsafe { }`: raw ptr deref, `transmute`, `*_unchecked` | borrow + type checker STILL run | Miri (dynamic UB detection) | yes (`*const`/`*mut`) |
| OCaml    | `Obj.magic`, `Bytes.unsafe_*`, `Bigarray` | nothing for `Obj.magic` | nothing (silent UB) | no (via C stubs / Bigarray) |
| Swift    | `withUnsafeBytes`, `Unsafe*Pointer`, `unsafeBitCast` | scoped to a closure; ARC still runs | AddressSanitizer / runtime traps | yes (scoped) |
| Kotlin   | `ByteBuffer`, FFM API (`java.lang.foreign`), (legacy `sun.misc.Unsafe`) | JVM verifier; GC yes | JVM (bounds-checked) / FFM confinement | no direct; off-heap via FFM |

## Reading

The spectrum runs from **least to most guarded**. `Obj.magic` in OCaml is the
sharpest knife: an unchecked cast with no verification and no tooling to catch a
mistake. Go's `unsafe` is next — the type checker and GC still run, so the ONE
way to get memory-unsafety is the `uintptr` laundering hazard, which the docs
enumerate precisely. Rust keeps the borrow checker on inside `unsafe` and, with
Miri, gives a way to *detect* the UB that raw pointers enable — the only language
here with a dedicated UB checker. Swift confines the hatch to scoped closures
with the runtime still active. The JVM is the most guarded: no raw pointers at
all, off-heap access mediated by a bounds-checked, confinement-checked FFM API.

The through-line: a language's safety is only as strong as how hard it makes the
escape, and how well it catches you when you take it. Go's position — type-safe
and GC-managed but with one explicit, documented lapse — is a deliberate middle,
consistent with the rest of these modules: a small set of guarantees, an escape
hatch when you need C-like control, and the honesty to tell you exactly where the
guarantee stops.

## References

See the module [`README`](README.md#references) for the per-language sources
(Go `unsafe`; Rustonomicon + Miri; OCaml `Obj`/`Bytes`/`Bigarray`; Swift unsafe
pointers; JVM FFM API + ByteBuffer).
