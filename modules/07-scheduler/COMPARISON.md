# 07 — Comparison: five scheduling models

**Environment.** Apple M4 Pro (14 cores), macOS arm64. Go 1.26.5, Rust 1.92.0,
OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). Measured throughput in
[`measured.txt`](measured.txt).

## The models

| Language | Unit | Mapping | Scheduler | Blocking I/O | CPU offload |
| -------- | ---- | ------- | --------- | ------------ | ----------- |
| Go       | goroutine | M:N over P=GOMAXPROCS | work-stealing, async preemption, netpoller | netpoller parks the G, frees the M | not needed (any G can block) |
| Rust     | `async` task (tokio) / OS thread | M:N (tokio) or 1:1 (`std::thread`) | tokio work-stealing; cooperative (await points) | `spawn_blocking` for blocking calls | `spawn_blocking` / `block_in_place` |
| OCaml 5  | domain / Eio fiber | domains 1:1; Eio fibers M:N on domains | Eio is effects-based, cooperative | Eio async I/O; a blocking call blocks the domain | run on another domain |
| Swift    | `Task` | M:N on a cooperative pool sized to cores | cooperative (suspension points) | `async` I/O; a blocking call blocks a pool thread | dispatch to a queue / detached task |
| Kotlin   | coroutine / thread / virtual thread | coroutines M:N; virtual threads M:N (Loom) | coroutines cooperative; Loom continuations | suspend for async; a blocking call blocks a carrier | `Dispatchers.IO` / a thread pool |

The deep split is **cooperative vs preemptive**. Go preempts (async
preemption + safepoints), so a CPU-heavy goroutine cannot monopolize a P and a
buggy tight loop cannot wedge the program. tokio, Swift, and Kotlin coroutines
are **cooperative**: a task yields only at an `await`/suspension point, so a
CPU-bound task with no await points can starve its worker thread — the classic
"don't block the async runtime" rule, and why all three provide an explicit
blocking-offload escape hatch (`spawn_blocking`, a dispatch queue,
`Dispatchers.IO`). Go's preemption removes that whole footgun class, at the cost
of a more complex runtime.

The second split is **who blocks on I/O**. Go's netpoller parks a goroutine
blocked on the network and hands the M back, so blocking looking code is
non-blocking underneath; the async runtimes require you to call the async API
(or offload), and OCaml domains block the whole domain unless you use Eio.

## What the numbers said

For a pure CPU-bound parallel sweep (`measured.txt`), all five are within ~15%:
the scheduler is not the bottleneck when the work is embarrassingly parallel and
CPU-saturating. Go's scheduler earns its keep elsewhere — the latency section in
the [`README`](README.md#latency-under-load) shows tail latency swinging three
orders of magnitude with worker count, and the cheap-goroutine result (over-
subscription helping, not hurting) is what a thread-based model cannot match.

## References

Official sources first; see the module [`README`](README.md#references) for the
full per-language list (Go scheduler doc + 1.14/1.25 notes + automaxprocs; tokio;
OCaml domains + Eio; Swift concurrency; Kotlin coroutines + Loom virtual
threads).
