# 12 — Comparison: five concurrent-cache implementations

**Environment.** Apple M4 Pro (14 cores), macOS arm64. Go 1.26.5, Rust 1.92.0,
OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). Measured numbers in
[`measured.txt`](measured.txt).

## The same design, five primitive sets

Every implementation is the same program: a `key -> value` map guarded for
mutual exclusion, capped in size by eviction, fronted by a counting semaphore
that bounds concurrent backend fetches, driven by 64 concurrent workers. What
differs is which primitives the language hands you.

| Language | Mutual exclusion           | Backpressure semaphore                                | Concurrency unit        | Backend-call counter |
| -------- | -------------------------- | ----------------------------------------------------- | ----------------------- | -------------------- |
| Go       | `sync.Mutex`               | buffered channel `chan struct{}` (cap = maxInflight)  | goroutine               | `atomic.Int64`       |
| Rust     | `Mutex<HashMap>`           | hand-rolled `Mutex<usize>` + `Condvar` (std has none) | `std::thread` + `Arc`   | `AtomicI64`          |
| OCaml 5  | `Mutex` + `Hashtbl`        | `Semaphore.Counting`                                  | domain (`Domain.spawn`) | `Atomic.t`           |
| Swift    | the `actor` itself         | (none — actor serializes)                             | `Task` in a `TaskGroup` | actor-isolated `var` |
| Kotlin   | `synchronized` + `HashMap` | `java.util.concurrent.Semaphore`                      | platform `Thread`       | `AtomicLong`         |

Three of the five (Go, Rust, OCaml, Kotlin) share the same skeleton: lock,
check the map, unlock, acquire the semaphore, fetch, release, lock, evict-if-
full, insert, unlock. The lock is held only around the map, never across the
backend sleep — holding it across the fetch would serialize the whole cache and
erase the point of the semaphore.

## The actor is the outlier

Swift's `actor` replaces both the mutex and the semaphore with one mechanism:
only one task runs actor code at a time, and the `await` at the backend sleep
_suspends_ the actor, letting other gets enter. There is no explicit lock and no
explicit semaphore — the model gives you mutual exclusion for free and turns the
backpressure question into "how many tasks are suspended at the await".

The consequence is visible in the numbers: the actor serializes map access
through a single executor, so its throughput is the lowest of the five, but it
has no lock convoy, so its p999 tail is an order of magnitude tighter than the
mutex implementations. That is the capstone's sharpest single finding: the actor
model is not strictly better or worse, it moves the cost from tail latency to
throughput, and the benchmark makes the trade concrete.

## Semaphore availability

A small but telling ergonomics split: OCaml (`Semaphore.Counting`), Kotlin
(`java.util.concurrent.Semaphore`), and Go (a channel is idiomatically a
semaphore) give you backpressure directly. Rust's std library has **no** counting
semaphore, so the implementation hand-rolls one from `Mutex<usize>` + `Condvar`
— correct and small, but a reminder that Rust pushes you toward `tokio::sync::
Semaphore` (a dependency) or a manual build. Swift needs neither because the
actor subsumes the role.

## Graceful shutdown

Only the Go reference wires context cancellation through the backpressure wait
(the `select` on `c.sem <- struct{}{}` vs `<-ctx.Done()`), and its
`TestGracefulShutdown` proves a cancelled context unblocks a `Get` stuck behind
a full semaphore rather than deadlocking it. The other four omit the transport-
level cancellation to stay focused on the cache core; adding it would mean a
`CancellationToken` (Kotlin), an `async` cancellation check (Swift), a
`select!`/`recv_timeout` (Rust), or an Eio switch (OCaml) — each idiomatic, none
changing the throughput story.

## Bottom line

For a warm, contended, memory-bound cache, the runtime choice is a ~2–3x
throughput factor and a large tail-latency factor, not a 10x one. The
mutex-based four cluster together; the actor trades throughput for a predictable
tail. Binary size and compile time (the migration matrix in `measured.txt`)
separate them further, but none of those costs dominates — which is exactly why
the decision should rest on measured trade-offs of the real workload, the thing
this capstone builds.
