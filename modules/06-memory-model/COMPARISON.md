# 06 — Comparison: five memory models

**Environment.** Apple M4 Pro (14 cores), macOS arm64. Go 1.26.5, Rust 1.92.0,
OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). Benchmark in
[`go/bench.txt`](go/bench.txt).

## Where data-race safety lives

| Language | Parallelism unit | Shared-state safety | Race is... | Safe primitive |
| -------- | ---------------- | ------------------- | ---------- | -------------- |
| Go       | goroutine        | runtime + `-race` (dynamic) | undefined behavior, sometimes detected | channel / `sync.Mutex` / `sync/atomic` |
| Rust     | `thread` / async | `Send`/`Sync` in the type system | **compile error** | `Arc<Mutex<_>>` / `Atomic*` |
| OCaml 5  | domain           | runtime (no static check) | undefined behavior | `Atomic` / `Mutex` |
| Swift 6  | task / actor     | `Sendable` + actor isolation | **compile error** | `actor` / `Sendable` types |
| Kotlin   | thread / coroutine | JVM/JMM (no static check) | undefined behavior | `Atomic*` / `synchronized` / `@Volatile` |

The split is stark: **Rust and Swift 6 make a data race unrepresentable**
(the racy program does not compile), while **Go, OCaml 5, and Kotlin/JVM make it
representable but undefined**, caught only dynamically (`-race`) or not at all.
That is the single most consequential axis for a concurrent system: a
compile-time guarantee versus a test-time detector plus discipline.

## The synchronized counter, five ways

<!-- BEGIN:snippet go-counters -->
```go
// RacyInc increments a shared int from many goroutines with NO synchronization.
// This is a data race: the result is nondeterministic AND the program has
// undefined behavior. -race flags it (see the captured report); the compiler
// does not. Rust and Swift 6 reject the analogous program at compile time.
func RacyInc() int {
	var counter int // shared, unguarded
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perG {
				counter++ // read-modify-write with no happens-before edge
			}
		}()
	}
	wg.Wait()
	return counter
}

// AtomicInc uses sync/atomic, which establishes the happens-before edges the
// memory model requires. Always returns exactly want.
func AtomicInc() int64 {
	var counter atomic.Int64
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perG {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
	return counter.Load()
}

// MutexInc uses a mutex; Lock/Unlock also establish happens-before.
func MutexInc() int {
	var (
		mu      sync.Mutex
		counter int
	)
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perG {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return counter
}
```
<!-- END:snippet go-counters -->

<!-- BEGIN:snippet rust-atomic -->
```rust
/// The synchronized counter: `Arc<AtomicU64>` is `Send + Sync`, so sharing it
/// across threads compiles, and the result is always deterministic.
pub fn atomic_count(threads: usize, per: u64) -> u64 {
    let counter = Arc::new(AtomicU64::new(0));
    let handles: Vec<_> = (0..threads)
        .map(|_| {
            let c = Arc::clone(&counter);
            thread::spawn(move || {
                for _ in 0..per {
                    c.fetch_add(1, Ordering::Relaxed);
                }
            })
        })
        .collect();
    for h in handles {
        h.join().unwrap();
    }
    counter.load(Ordering::SeqCst)
}
```
<!-- END:snippet rust-atomic -->

<!-- BEGIN:snippet ocaml-atomic -->
```ocaml
let atomic_count threads per =
  let counter = Atomic.make 0 in
  let domains =
    List.init threads (fun _ ->
        Domain.spawn (fun () ->
            for _ = 1 to per do
              Atomic.incr counter
            done))
  in
  List.iter Domain.join domains;
  Atomic.get counter
```
<!-- END:snippet ocaml-atomic -->

<!-- BEGIN:snippet swift-actor -->
```swift
public actor Counter {
  private var n = 0
  public init() {}
  public func inc() { n += 1 }
  public func value() -> Int { n }
}

/// Increments an actor from `threads` concurrent child tasks. The actor's
/// isolation makes every `inc()` mutually exclusive without an explicit lock.
public func count(_ threads: Int, _ per: Int) async -> Int {
  let c = Counter()
  await withTaskGroup(of: Void.self) { group in
    for _ in 0..<threads {
      group.addTask {
        for _ in 0..<per { await c.inc() }
      }
    }
  }
  return await c.value()
}
```
<!-- END:snippet swift-actor -->

<!-- BEGIN:snippet kotlin-atomic -->
```kotlin
fun atomicCount(
    threads: Int,
    per: Int,
): Int {
    val counter = AtomicInteger(0)
    val ts =
        (1..threads).map {
            Thread {
                repeat(per) { counter.incrementAndGet() }
            }
        }
    ts.forEach { it.start() }
    ts.forEach { it.join() }
    return counter.get()
}
```
<!-- END:snippet kotlin-atomic -->

## Outputs

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
atomic_count(8, 100000) = 800000 (correct)
```
<!-- END:output rust-demo -->

<!-- BEGIN:output ocaml-demo -->
```text
$ dune exec bin/main.exe
atomic_count 8 100000 = 800000 (correct)
```
<!-- END:output ocaml-demo -->

<!-- BEGIN:output swift-demo -->
```text
actor count(8, 10000) = 80000 (correct)
```
<!-- END:output swift-demo -->

<!-- BEGIN:output kotlin-demo -->
```text
atomicCount(8, 100000) = 800000 (correct)
```
<!-- END:output kotlin-demo -->

## Atomic vs mutex under contention (Go)

Both are correct; the cost differs. Measured with `b.RunParallel`
(`-count=10`, `benchstat`):

<!-- BEGIN:file go-bench -->
```text
goos: darwin
goarch: arm64
pkg: github.com/dannywillems/go-for-systems-engineers/modules/06-memory-model/go
cpu: Apple M4 Pro
          │     go      │
          │   sec/op    │
Atomic-14   51.26n ± 5%
Mutex-14    107.1n ± 1%
geomean     74.09n

          │      go      │
          │     B/op     │
Atomic-14   0.000 ± 0%
Mutex-14    0.000 ± 0%
geomean                ¹
¹ summaries must be >0 to compute geomean

          │      go      │
          │  allocs/op   │
Atomic-14   0.000 ± 0%
Mutex-14    0.000 ± 0%
geomean                ¹
¹ summaries must be >0 to compute geomean
```
<!-- END:file go-bench -->

On this machine the atomic counter is roughly half the cost of the mutex under
full contention, with zero allocations for both. This gap is workload-dependent
(a mutex protecting more than one word, or lower contention, changes it), which
is exactly why the number is measured here rather than asserted.

## Reading

Go's model is deliberately permissive: cheap goroutines, a small set of
synchronization primitives, and a dynamic detector rather than a type-system
proof. The cost is that correctness under concurrency is a property of the tests
and the reviewer, not the compiler — the opposite of Rust and Swift 6, which pay
in type-system complexity to make the guarantee static. OCaml 5's domains and
Kotlin/JVM sit with Go on this axis (runtime primitives, no static race check),
though the JVM has decades of concurrent-library maturity behind it.

## References

Official sources first. See the module [`README.md`](README.md#references) for
the full list (Go Memory Model, race detector, `sync/atomic`; Rust `Send`/`Sync`
and the Rustonomicon; OCaml domains and `Atomic`; Swift `Sendable`; the Java
Memory Model / JLS ch. 17).
