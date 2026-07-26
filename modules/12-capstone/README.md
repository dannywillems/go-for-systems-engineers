# 12 — Capstone: a concurrent bounded cache

**Thesis.** The five runtime stories of Modules 05–07 (shared mutable state
under a lock, bounded memory via eviction, backpressure via a semaphore,
graceful shutdown via cancellation) compose into one small program. Built the
same way in all five languages, the concurrent cache is the substance a
migration decision actually turns on — request plumbing is not — and its
measured throughput, tail latency, binary size, and lines of code are the
comparison, not adjectives about the languages.

The HTTP transport is intentionally omitted: it would drag in a web framework
per language (and a Gradle/ktor/NIO stack for Kotlin) without changing the
concurrency semantics under test. Everything here is dependency-free and
directly comparable.

## Contents

- [The workload](#the-workload)
- [The reference implementation (Go)](#the-reference-implementation-go)
- [Measured results](#measured-results)
- [Reading the numbers](#reading-the-numbers)
- [The migration matrix](#the-migration-matrix)
- [Exercises](#exercises)
- [References](#references)

## The workload

Every implementation runs the identical shape: 64 worker threads/goroutines/
tasks each issue 10,000 `Get`s (640,000 total) against a 256-entry bounded
cache over a backend whose `fetch` sleeps 100 µs. Keys are drawn from an LCG
over exactly `capacity` distinct values, so the hot set equals the cache size
and the cache absorbs ~99.9% of the load once warm. A `maxInflight` semaphore
caps concurrent backend fetches (backpressure). Each `Get`'s latency is
recorded; the harness reports throughput (kops/s), p50/p99/p999, and the
fraction of gets that reached the backend.

Correctness is a separate red/green test per language (`cache_test.go`, the
Rust `#[test]`, the OCaml/Swift/Kotlin test binaries): under concurrent load the
cache must return the correct value for every key, never exceed capacity, and
actually reduce backend calls.

## The reference implementation (Go)

The Go cache is the reference the other four mirror. The two bounded resources
are explicit: the map is guarded by a `sync.Mutex` and capped by eviction, and
concurrent backend fetches are capped by a buffered channel used as a counting
semaphore. The `select` on `ctx.Done()` alongside the semaphore send is what
makes shutdown graceful — a cancelled context unblocks a waiting `Get` instead
of deadlocking it:

<!-- BEGIN:snippet go-cache -->
```go
// Cache is a concurrent, bounded key->value cache over a Backend. It caps the
// number of entries (evicting an arbitrary one on overflow) and the number of
// concurrent backend fetches (a semaphore = backpressure).
type Cache struct {
	entries  map[int]int
	backend  *Backend
	sem      chan struct{} // backpressure: bounds concurrent fetches
	mu       sync.Mutex
	capacity int
}

func NewCache(capacity, maxInflight int, b *Backend) *Cache {
	return &Cache{
		entries:  make(map[int]int, capacity),
		backend:  b,
		sem:      make(chan struct{}, maxInflight),
		capacity: capacity,
	}
}

// Get returns the cached value or fetches it from the backend, respecting the
// context (graceful shutdown) and the in-flight backpressure limit.
func (c *Cache) Get(ctx context.Context, key int) (int, error) {
	c.mu.Lock()
	if v, ok := c.entries[key]; ok {
		c.mu.Unlock()
		return v, nil
	}
	c.mu.Unlock()

	// Backpressure: block until a fetch slot is free, or the context is done.
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	v := c.backend.fetch(key)

	c.mu.Lock()
	if len(c.entries) >= c.capacity {
		for k := range c.entries { // evict one arbitrary entry
			delete(c.entries, k)
			break
		}
	}
	c.entries[key] = v
	c.mu.Unlock()
	return v, nil
}
```
<!-- END:snippet go-cache -->

The other four express the same design with their native primitives: Rust a
`Mutex<HashMap>` plus a hand-rolled `Condvar` semaphore (std has no counting
semaphore); OCaml a `Mutex`-guarded `Hashtbl` plus `Semaphore.Counting` across
domains; Kotlin a `synchronized` map plus `java.util.concurrent.Semaphore` on
platform threads; and Swift an `actor`, where the mutual exclusion is the actor
itself and the `await` at the backend fetch suspends it so fetches overlap while
map access stays serialized. See [`COMPARISON.md`](COMPARISON.md) for the
side-by-side.

## Measured results

Regenerated locally by [`scripts/capstone-bench.sh`](../../scripts/capstone-bench.sh)
(non-portable: timings are machine-specific and this block is skipped by the
docs-freshness gate).

<!-- BEGIN:file measured -->
```text
machine: Apple M4 Pro (14 cores), macOS arm64
workload: 640k concurrent Get / 64 workers over a 256-entry bounded
cache; backend fetch sleeps 100us; hot set == capacity so the cache
absorbs ~99.9% of load. HTTP transport omitted (dependency-free).

# Throughput and latency (single run, optimized)
Go     640k gets/64w: 6677 kops/s  p50=200ns p99=298µs p999=953.3µs  backend=0.1% of gets
Rust   640k gets/64w: 7122 kops/s  p50=42ns p99=346.125µs p999=948.958µs  backend=0.1% of gets
OCaml  640k gets/64w: 5413 kops/s  p50=0.0us p99=266.1us p999=1241.9us  backend=0.1% of gets
Swift  640k gets/64w: 2794 kops/s  p50=21.0us p99=41.2us p999=59.9us  backend=0.0% of gets
Kotlin 640k gets/64w: 7357 kops/s  p50=0.1us p99=374.0us p999=622.7us  backend=0.1% of gets

# Migration matrix
lang     binary         core-LOC   cold-compile-s
Go       2580770        83         1.5
Rust     495136         100        0.4
OCaml    1501304        43         0.4
Swift    91000          28         2.0
Kotlin   5614009        31         2.2
```
<!-- END:file measured -->

## Reading the numbers

- **All five saturate on the same wall.** Throughput lands in a ~2.7x band
  (Go/Rust/Kotlin ~7k kops/s at the top, Swift lowest) even though the
  concurrency models differ completely. The workload is memory- and
  lock-bound, not language-bound; once the cache is warm the backend sleep is
  off the hot path, so the map+lock throughput dominates. This is the same
  lesson as Module 07: for this class of work the runtime choice is a ~2x
  factor, not a 10x one.
- **The actor trades throughput for a tight tail.** Swift's `actor` serializes
  all cache access through one executor, which caps throughput — but its p999 is
  an order of magnitude _tighter_ than the mutex-based implementations
  (tens of µs vs ~1 ms). A single serialization point has no lock-convoy tail;
  the mutex implementations show a fat p999 from threads piling up on the lock
  under contention. Whether you want the throughput or the predictable tail is a
  real design axis, and the numbers make it concrete rather than folkloric.
- **Backend load collapses to ~0.1%.** All five drive the backend to the same
  ~0.1% of gets: the eviction policy and hot-set sizing, not the language,
  determine hit rate. The cache does its job identically everywhere.

## The migration matrix

The `# Migration matrix` block in `measured.txt` reports the operational costs
a migration weighs alongside throughput: stripped binary size, core-logic lines
of code, and a cold (clean) compile time. Read them as order-of-magnitude — the
LOC counts reflect comment density as much as language density, and cold-compile
times are noisy on a shared machine.

The shape that survives the noise: Swift and Rust produce the smallest binaries
(the Swift binary dynamically links its runtime; Go and Kotlin bundle theirs, so
Go is a few MB and the Kotlin fat jar is the largest). Rust and OCaml compile
this crate fastest from clean; Swift and Kotlin are slowest. None of these is
decisive on its own; the point of the capstone is that they are _measured
together_ against the one program whose semantics the migration actually
depends on.

## Exercises

[`exercises/go`](exercises/go) is red until you implement a context-aware
bounded semaphore and the get-or-fetch path. [`solutions/go`](solutions/go) is
the verified corrigé. Run with `-race`:

```
make exercises M=12   # red
make solutions M=12   # green
```

## References

Official sources first, grouped by language.

### Go

- `sync` (Mutex): https://pkg.go.dev/sync
- `context` (cancellation, deadlines): https://pkg.go.dev/context
- Effective Go, "Channels" (channel-as-semaphore): https://go.dev/doc/effective_go#channels
- `sync/atomic`: https://pkg.go.dev/sync/atomic

### Rust

- `std::sync::{Mutex, Condvar}`: https://doc.rust-lang.org/std/sync/
- `std::sync::atomic`: https://doc.rust-lang.org/std/sync/atomic/
- The Rustonomicon, "Atomics": https://doc.rust-lang.org/nomicon/atomics.html

### OCaml

- OCaml Manual, Parallelism (domains): https://ocaml.org/manual/5.4/parallelism.html
- `Semaphore.Counting`: https://ocaml.org/manual/5.4/api/Semaphore.Counting.html
- `Atomic`: https://ocaml.org/manual/5.4/api/Atomic.html

### Swift

- Actors: https://docs.swift.org/swift-book/documentation/the-swift-programming-language/concurrency/#Actors
- `TaskGroup`: https://developer.apple.com/documentation/swift/taskgroup
- SE-0306, Actors: https://github.com/apple/swift-evolution/blob/main/proposals/0306-actors.md

### Kotlin (JVM)

- `java.util.concurrent.Semaphore`: https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/concurrent/Semaphore.html
- `synchronized`: https://kotlinlang.org/docs/functions.html
- `LockSupport`: https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/concurrent/locks/LockSupport.html
