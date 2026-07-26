# 06 — Memory model and data races

**Thesis.** The Go Memory Model defines happens-before through channel ops,
mutexes, and `sync/atomic`; a read not ordered against a concurrent write is a
data race, and "a program with data races has undefined behavior" — including
torn reads of multiword values. Go does **not** reject racy programs at compile
time (the `-race` detector is *dynamic*: it only sees races that actually occur
during a run). Rust and Swift 6 push this into the type system and make the same
race a compile error; OCaml 5 and Kotlin/JVM, like Go, rely on runtime
primitives and discipline.

## A data race, and the dynamic detector

An unsynchronized shared counter is a race:

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

The synchronized versions are deterministic; the racy one is not, so its result
is never captured here — the `-race` detector is the falsifiable evidence
instead (run against the build-tagged `TestDataRace`, addresses and goroutine
ids normalized):

<!-- BEGIN:output go-race -->
```text
WARNING: DATA RACE
Read at 0xADDR by goroutine N:
Previous write at 0xADDR by goroutine N:
```
<!-- END:output go-race -->

The detector is dynamic. It reports a race only on the interleaving that
actually happened; a race on a code path not exercised in the run is **not**
reported. `-race` is a testing tool, not a proof of race-freedom. Under a race,
Go is not even memory-safe: a concurrent write to a two-word `interface` or
slice header can be read torn (type word from one value, data word from
another), which is undefined behavior, not just a wrong number.

The synchronized results, deterministic and captured:

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
AtomicInc() = 800000 (correct)
MutexInc()  = 800000 (correct)
ParallelSquares([1..5], limit=2) = [1 4 9 16 25] err=<nil>
```
<!-- END:output go-demo -->

## Goroutine leaks

A goroutine with no exit path is a permanent leak (goroutines are not GC'd).
`goleak` catches it in tests:

<!-- BEGIN:output go-leak -->
```text
    leak_test.go:23: found unexpected goroutines:
```
<!-- END:output go-leak -->

## Structured concurrency with errgroup

The idiomatic replacement for hand-rolled goroutine+channel plumbing:
`errgroup.WithContext` (first error cancels) plus `SetLimit` (bounded
concurrency), each worker writing a distinct index (no shared write):

<!-- BEGIN:snippet go-pipeline -->
```go
// ParallelSquares squares each input concurrently with bounded parallelism,
// cancelling all workers on the first error. errgroup.WithContext gives the
// first-error-cancels behavior and SetLimit bounds concurrency; each worker
// writes a DISTINCT index, so there is no shared-write race. This is the pattern
// that replaces hand-rolled goroutine + channel + error plumbing.
func ParallelSquares(ctx context.Context, in []int, limit int) ([]int, error) {
	out := make([]int, len(in))
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit)
	for i, v := range in {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			out[i] = v * v
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return out, nil
}
```
<!-- END:snippet go-pipeline -->

## The same race is a compile error elsewhere

Rust's `Send`/`Sync` and Swift 6's `Sendable` reject the racy program before it
runs. Neither of these compiles:

<!-- BEGIN:output rust-reject -->
```text
error[E0277]: `Rc<i32>` cannot be sent between threads safely
```
<!-- END:output rust-reject -->

<!-- BEGIN:output swift-reject -->
```text
reject.swift:10:5: error: main actor-isolated var 'counter' can not be mutated from a nonisolated context
```
<!-- END:output swift-reject -->

See [`COMPARISON.md`](COMPARISON.md) for the mutex-vs-atomic benchmark and the
five memory models (Go happens-before, Rust `Send`/`Sync`, OCaml 5 domains +
`Atomic`, Swift 6 actors/`Sendable`, JVM/JMM), and [`exercises/`](exercises).

## References

Official sources first.

- Go Memory Model (normative): https://go.dev/ref/mem
- Data Race Detector (official guide): https://go.dev/doc/articles/race_detector
- "Introducing the Go Race Detector" (Go blog): https://go.dev/blog/race-detector
- `sync/atomic` package docs: https://pkg.go.dev/sync/atomic
- `sync` package docs: https://pkg.go.dev/sync
- `golang.org/x/sync/errgroup`: https://pkg.go.dev/golang.org/x/sync/errgroup
- `context` package docs: https://pkg.go.dev/context
- goleak (Uber): https://github.com/uber-go/goleak
- Rust `Send`/`Sync` (Reference): https://doc.rust-lang.org/reference/special-types-and-traits.html#send
- The Rustonomicon, Send and Sync: https://doc.rust-lang.org/nomicon/send-and-sync.html
- OCaml Manual, Parallelism (domains): https://ocaml.org/manual/5.4/parallelism.html
- OCaml `Atomic` module: https://ocaml.org/manual/5.4/api/Atomic.html
- Swift `Sendable` protocol: https://developer.apple.com/documentation/swift/sendable
- Swift concurrency (The Swift Programming Language): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/concurrency/
- Java Language Specification, ch. 17 (Java Memory Model): https://docs.oracle.com/javase/specs/jls/se21/html/jls-17.html
- `java.util.concurrent.atomic` (used from Kotlin): https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/util/concurrent/atomic/package-summary.html
