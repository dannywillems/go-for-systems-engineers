# 11 — Comparison: five testing toolchains

**Environment.** Go 1.26.5, Rust 1.92.0, OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10.
The Go demo, the caught-bug output, and every snippet are injected from the code.

## The four techniques across the five languages

| Technique | Go | Rust | OCaml | Swift | Kotlin/JVM |
| --------- | -- | ---- | ----- | ----- | ---------- |
| Unit / table | `testing` + `t.Run` (stdlib) | `#[test]` (built in) | Alcotest / `assert` | Swift Testing `@Test` | kotlin.test / JUnit |
| Table-driven | slice of cases + `t.Run` | a loop of `assert_eq!` | a list + `List.iter` | `@Test(arguments:)` (parameterized) | `@ParameterizedTest` / a loop |
| Property-based | `testing/quick` (stdlib) | proptest / quickcheck (crate) | QCheck | SwiftCheck | Kotest property |
| Fuzzing | `go test -fuzz` (stdlib, native) | cargo-fuzz (libFuzzer) | crowbar / afl | libFuzzer via SwiftPM | Jazzer / JQF |
| Virtual time | `testing/synctest` (stdlib) | `tokio::time::pause`/`advance` | manual clock / Eio mock | `Clock` protocol injection | `runTest` (kotlinx-coroutines-test) |

The standout is the "stdlib" column for Go: property testing, fuzzing, AND
virtual time are all in the standard library and the `go test` command. In the
other four, property testing and fuzzing are third-party (proptest, QCheck,
SwiftCheck, Kotest; cargo-fuzz, crowbar, Jazzer). That is the same
batteries-included stance seen across these modules — fewer external choices, one
blessed way — with the same trade: the third-party tools (proptest's shrinking,
Kotest's generators) are often richer than the stdlib equivalent
(`testing/quick` has no shrinking).

## Property testing is the load-bearing idea

The repo hand-rolls a property loop in Rust/OCaml/Swift/Kotlin (idempotence over
an LCG-generated corpus) precisely to show the idea is not framework-specific:
state an invariant, generate many inputs, assert. What the real frameworks add
on top is worth naming:

- **Generation**: typed, composable generators (proptest `Strategy`, Kotest
  `Arb`, QCheck `Gen`) versus `testing/quick`'s reflection-driven defaults.
- **Shrinking**: on failure, proptest/quickcheck/Kotest reduce the counterexample
  to a minimal one ("" or "a a", not a 4KB blob). `testing/quick` does not shrink
  — a real gap when a failure is a huge random string.
- **Reproducibility**: a printed seed to replay the exact failing run.

For the invariant checks in this module, the hand-rolled loop is enough; for
production property testing, the shrinking alone justifies the dependency.

## Virtual time is the newest and the most valuable

Every one of the five now has a way to make time-dependent tests deterministic,
and it is the highest-leverage testing feature of the last few years because it
kills a whole category of flaky, slow, sleep-based tests:

- **Go** `testing/synctest` fakes time for an entire goroutine bubble, with no
  change to the code under test — the limiter calls `time.Now`/`time.Sleep`
  normally. This is the least invasive of the five.
- **Rust** `tokio::time::pause`/`advance` fakes time for the tokio runtime, but
  only for `tokio::time` (not `std::thread::sleep`), so the code must be async on
  tokio.
- **Kotlin** `runTest` gives coroutines a virtual scheduler; `delay` is skipped.
  Same async-only caveat.
- **Swift** injects a custom `Clock` — explicit, requires the code to be written
  against the `Clock` protocol rather than a global now.
- **OCaml** has no standard virtual-time facility; you inject a clock or mock an
  Eio time source manually.

Go's is uniquely non-invasive because it intercepts the real `time` package
inside the bubble; the others require the code under test to already route time
through the runtime or an injected abstraction.

## Bottom line

For "does this pure function hold its invariant", every language's property
tooling is fine and the idea ports in a dozen lines. For "does this concurrent,
time-dependent code behave", virtual time is the difference between a fast
deterministic suite and a slow flaky one, and Go's `testing/synctest` is the
least intrusive way in — no dependency, no clock injection, no restructuring.
That, plus native fuzzing and property testing in the standard toolchain, makes
Go's out-of-the-box test story the broadest of the five, at the cost of the
richer generation and shrinking the dedicated third-party frameworks provide.
