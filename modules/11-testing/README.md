# 11 — Testing

**Thesis.** This whole repository is an argument that claims should be
mechanically checked, so its closing module is about the checking itself. Go's
standard `testing` package covers four techniques with no third-party
dependency: table-driven tests, property-based tests (`testing/quick`), native
fuzzing (`go test -fuzz`), and — the newest and most consequential —
deterministic virtual time for concurrent code (`testing/synctest`, stabilized
in Go 1.25). The other four languages reach the same techniques through their
own frameworks. The point property and fuzz testing make, and the reason they
matter more than examples: a machine finds the violating input you never thought
to write a unit test for.

## Contents

- [The subject under test](#the-subject-under-test)
- [Property-based testing](#property-based-testing)
- [Native fuzzing](#native-fuzzing)
- [When the tooling catches a bug](#when-the-tooling-catches-a-bug)
- [Deterministic virtual time: testing/synctest](#deterministic-virtual-time-testingsynctest)
- [Exercises](#exercises)
- [References](#references)

## The subject under test

`Normalize` collapses whitespace, trims, and lowercases. Its invariant is
idempotence — `Normalize(Normalize(s)) == Normalize(s)` — which is exactly the
kind of property a machine can check on thousands of inputs:

<!-- BEGIN:snippet go-normalize -->
```go
// Normalize collapses all whitespace runs to a single space, trims the ends, and
// lowercases. It is idempotent: Normalize(Normalize(s)) == Normalize(s) for all
// s -- the invariant the property and fuzz tests check.
func Normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}
```
<!-- END:snippet go-normalize -->

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
Normalize("  Hello   World  ") = "hello world"
Normalize("MiXeD\tCase") = "mixed case"
Normalize("   ") = ""
Normalize("a\n\nb") = "a b"
```
<!-- END:output go-demo -->

## Property-based testing

A property test states an invariant and lets the framework generate inputs.
`testing/quick` is in the standard library; it manufactures random values for
the function's parameters and fails with a counterexample if the property is
ever false:

<!-- BEGIN:snippet go-property -->
```go
// Property-based test with the stdlib testing/quick: for ARBITRARY strings,
// Normalize is idempotent. quick.Check generates random inputs and fails with a
// counterexample if the property is ever violated.
func TestNormalizeIdempotent(t *testing.T) {
	idempotent := func(s string) bool {
		return Normalize(Normalize(s)) == Normalize(s)
	}
	if err := quick.Check(idempotent, nil); err != nil {
		t.Fatal(err)
	}
}
```
<!-- END:snippet go-property -->

This is strictly stronger than the table-driven test above: the table checks the
handful of cases you imagined; the property checks shapes you did not.

## Native fuzzing

Fuzzing is property testing with coverage-guided input mutation. Go's is native:
`func FuzzX(f *testing.F)`, seeded with `f.Add`, run either against the seed
corpus (plain `go test`) or as a mutation search (`go test -fuzz`). It shines on
parsers and any code that must not misbehave on arbitrary bytes, including
invalid UTF-8:

<!-- BEGIN:snippet go-fuzz -->
```go
// Native fuzz test: `go test -fuzz=FuzzNormalize` mutates the seed corpus to
// hunt for an input that breaks an invariant; plain `go test` runs it against
// the seeds. The invariants: idempotent, trimmed, and no double spaces -- for
// ANY input, including invalid UTF-8.
func FuzzNormalize(f *testing.F) {
	f.Add("  Hello   World  ")
	f.Add("")
	f.Add("\t\nMiXeD\r\n")
	f.Fuzz(func(t *testing.T, s string) {
		got := Normalize(s)
		if Normalize(got) != got {
			t.Errorf("not idempotent on %q", s)
		}
		if strings.TrimSpace(got) != got {
			t.Errorf("leading/trailing space in %q", got)
		}
		if strings.Contains(got, "  ") {
			t.Errorf("double space in %q", got)
		}
	})
}
```
<!-- END:snippet go-fuzz -->

## When the tooling catches a bug

The value is not that these tests pass — it is that they FAIL on a real defect.
[`reject-buggy`](reject-buggy) plants a bug (a `Normalize` that only lowercases,
forgetting to collapse and trim) and runs the same fuzz invariants. On the very
first seed the tooling reports the exact violating value:

<!-- BEGIN:output go-catch -->
```text
        buggy_test.go:15: invariant violated: double space in "  hello   world  "
```
<!-- END:output go-catch -->

No one wrote a unit test for "a string with a double space"; the invariant plus
the seed found it. That is the whole argument for property and fuzz testing over
hand-picked examples.

## Deterministic virtual time: testing/synctest

The hardest thing to test is time-dependent concurrent code: a rate limiter, a
retry with backoff, a cache TTL. The classic approaches are a real
`time.Sleep` (slow and flaky) or a hand-injected clock interface (invasive).
`testing/synctest` (Go 1.25) runs a goroutine "bubble" in which `time.Now`,
`time.Sleep`, and all timers use VIRTUAL time that advances instantly once every
goroutine in the bubble is blocked. The rate-limiter test sleeps for a window
and the test still runs in ~0 seconds, deterministically:

<!-- BEGIN:snippet go-synctest -->
```go
// TestLimiterRefill runs inside a synctest bubble: time is VIRTUAL, so the
// Sleep is instantaneous and the test is deterministic -- no real wall-clock
// wait, no flakiness. This is the Go 1.25 way to test time-dependent code.
func TestLimiterRefill(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := NewLimiter(2, 100*time.Millisecond)
		first, second := l.Allow(), l.Allow() // consume both tokens
		if !first || !second {
			t.Fatal("first two Allow should succeed")
		}
		if l.Allow() {
			t.Fatal("third Allow should fail (bucket empty)")
		}
		time.Sleep(100 * time.Millisecond) // virtual time: instant
		if !l.Allow() {
			t.Fatal("Allow should succeed after the window refills")
		}
	})
}
```
<!-- END:snippet go-synctest -->

This is the single most impactful recent addition to Go's test story: an entire
class of flaky, slow, sleep-riddled tests becomes fast and deterministic without
restructuring the code under test to accept a clock. Rust's `tokio::time::pause`
/ `advance` and Kotlin's `runTest` virtual time do the same for their async
runtimes; see [`COMPARISON.md`](COMPARISON.md).

## Exercises

[`exercises/go`](exercises/go) is red until you (a) implement `Dedup` (remove
consecutive duplicates) so its property test passes, and (b) fill in a
`testing/quick` property asserting `Dedup` is idempotent.
[`solutions/go`](solutions/go) is the verified corrigé:

```
make exercises M=11   # red
make solutions M=11   # green
```

## References

Official sources first, grouped by language.

### Go

- `testing` (table-driven, `T.Run`, `F.Fuzz`): https://pkg.go.dev/testing
- `testing/quick` (property-based): https://pkg.go.dev/testing/quick
- Go fuzzing (the tutorial + reference): https://go.dev/doc/security/fuzz/
- `testing/synctest` (virtual time): https://pkg.go.dev/testing/synctest
- Go 1.25 release notes (synctest stabilized): https://go.dev/doc/go1.25#testingsynctest

### Rust

- `#[test]` and the test harness: https://doc.rust-lang.org/book/ch11-01-writing-tests.html
- proptest (property testing): https://proptest-rs.github.io/proptest/
- quickcheck: https://github.com/BurntSushi/quickcheck
- cargo-fuzz (libFuzzer): https://rust-fuzz.github.io/book/
- `tokio::time::pause` / `advance` (virtual time): https://docs.rs/tokio/latest/tokio/time/fn.pause.html

### OCaml

- Alcotest (unit test framework): https://github.com/mirage/alcotest
- QCheck (property testing): https://github.com/c-cube/qcheck
- Crowbar / afl (fuzzing): https://github.com/stedolan/crowbar

### Swift

- Swift Testing (`@Test`, `#expect`, parameterized): https://developer.apple.com/documentation/testing
- SwiftCheck (property testing): https://github.com/typelift/SwiftCheck
- `Clock` protocol (injectable time): https://developer.apple.com/documentation/swift/clock

### Kotlin (JVM)

- kotlin.test: https://kotlinlang.org/api/latest/kotlin.test/
- Kotest property testing: https://kotest.io/docs/proptest/property-based-testing.html
- kotlinx-coroutines-test (`runTest`, virtual time): https://kotlinlang.org/api/kotlinx.coroutines/kotlinx-coroutines-test/
- Jazzer (JVM fuzzing): https://github.com/CodeIntelligenceTesting/jazzer
