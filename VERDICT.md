# Verdict

Cross-cutting synthesis across the modules built so far. Every claim points to
the module that measured it; nothing here is asserted without a captured
program behind it. This file evolves as modules land (08–12 still pending).

## Where Go diverges from a Rust/OCaml intuition, and the cost

- **No coproducts (sum types).** `(T, error)` is a product, not a sum
  ([04](modules/04-errors)); the "sealed interface + type switch" idiom has no
  compiler exhaustiveness check ([02](modules/02-sum-types)). Cost: a class of
  "forgot a case" and "used the value on error" bugs that are silent at runtime.
  Recovered partially by an external `go/analysis` pass, not the type system.
- **No generic methods, no higher-kinded types.** A `Result` monad is lawful but
  unusable because `Map`/`AndThen` cannot be methods ([04](modules/04-errors));
  `Functor`/`Monad` abstractions do not typecheck ([03](modules/03-generics),
  [14](modules/14-type-system)). Cost: no do-notation, no method chaining, more
  free-function nesting.
- **Structural interfaces, implicit existentials.** No orphan rule, no
  coherence to protect ([01](modules/01-interfaces)); the itab appears wherever
  a value flows into an interface. Cost: a two-word value, an indirect call, and
  lost inlining — the interface path was ~3.5x the direct path in the module's
  benchmark.
- **Managed memory with one explicit lapse.** Type-safe and GC-managed, but a
  `uintptr` is untracked, so `unsafe` can produce use-after-free by request
  ([13](modules/13-unsafe)). No compile-time data-race prevention: races are
  undefined behavior caught only dynamically by `-race`
  ([06](modules/06-memory-model)), where Rust and Swift 6 reject them at compile
  time.

## The runtime trade-offs, in numbers (M4 Pro)

- **Dispatch** ([01](modules/01-interfaces)): direct 747ns, interface 2.6µs,
  devirtualized 1.34µs — the cost is lost inlining, not the indirect jump.
- **Generics** ([03](modules/03-generics)): generic (239ns) ≈ concrete (240ns)
  for value types (GCShape stenciling), ≪ interface (1.05µs). Pointer types
  collapse to one shared stencil + a dictionary.
- **GC** ([05](modules/05-memory)): `GOGC` is a throughput/RSS dial (100→800:
  ~9x fewer collections, ~5x RSS, faster). Counter-intuitively the tracing GCs
  (JVM, OCaml) beat Rust's per-object `malloc`/`free` for high-churn short-lived
  allocation, and OCaml's minor GC beat Go's.
- **Synchronization** ([06](modules/06-memory-model)): atomic 51ns vs mutex
  107ns under contention.
- **Scheduler** ([07](modules/07-scheduler)): task-latency p50 swings from 31µs
  (GOMAXPROCS workers) to 28ms (workers=2); oversubscription *helps* because
  goroutines are cheap. Pure CPU parallel throughput is within ~15% across all
  five languages — the scheduler model is not the bottleneck there.

## What each language made easy vs painful (so far)

- **Go**: easy — cheap concurrency, one visible error idiom, a strong external
  tooling and dynamic-checker culture, fast compiles, small type system to
  learn. Painful — the missing coproducts/generic-methods/HKT, and correctness
  under concurrency living in tests + review rather than the compiler.
- **Rust**: easy — compile-time data-race and memory safety, zero-cost generics,
  typestate. Painful — the borrow checker's learning curve, monomorphization
  compile-time cost (its small binaries here were a runtime-floor artifact, not
  a monomorphization win).
- **OCaml**: easy — the richest type-level modelling (GADTs, functors,
  first-class modules), a fast minor GC. Painful — thinner lint/security
  tooling, a young multicore story (domains 2022).
- **Swift**: easy — value types + ARC + Swift-6 `Sendable` compile-time race
  safety, expressive protocols. Painful — ARC was slowest under allocation
  churn; strict-concurrency adoption cost.
- **Kotlin/JVM**: easy — mature GC (fastest under allocation churn here), rich
  ecosystem, coroutines. Painful — type erasure limits, unchecked exceptions,
  no static race prevention.

## The axes a choice actually turns on

Stated as measured trade-offs, not advocacy:

1. **Compile-time vs test-time correctness.** Rust/Swift-6 move data-race and
   exhaustiveness guarantees into the type checker; Go/Kotlin rely on `-race`,
   analyzers, and discipline. Worth it in proportion to how costly a
   production concurrency/logic bug is for the team.
2. **Type-system size vs expressiveness.** Go's small system is fast to learn
   and compile and hard to over-abstract; the ML-lineage systems model more but
   ask more. The gaps (02–04, 14) are the measured price of Go's choice.
3. **Allocation profile, not "GC vs no-GC".** For short-lived high-churn work a
   tracing GC can beat manual free (05); for large/long-lived data or when heap
   allocation can be avoided, Rust's model wins. Measure the workload.
4. **Concurrency shape.** Go's cheap preemptive goroutines make oversubscription
   safe and blocking code non-blocking under the hood (07); the cooperative
   async runtimes need explicit blocking-offload. For embarrassingly parallel
   CPU work, all five are within ~15%.

No module found a universal winner, and this file will not invent one. The point
of the repository is to make each trade concrete and reproducible.
