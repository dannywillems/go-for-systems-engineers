# 04 — Errors

**Thesis.** Go's `(T, error)` is a **product** `T × error`, not a coproduct
`T + error`. Both components always exist; which parts are valid is per-function
convention, not a type. A Rust or OCaml engineer's `Result` reflex — "on error,
bail" — is therefore wrong on the readers that return a value AND an error
together. This module demonstrates the product/coproduct gap, the typed-nil
trap, the `#[must_use]` hole, and builds a lawful-but-unusable `Result[T]` to
show exactly which piece Go's generics are missing.

## `(T, error)` is a product, not a sum

`io.Reader` may return `n > 0` bytes and `io.EOF` on the same `Read`:

<!-- BEGIN:snippet go-reader -->
```go
// eofReader yields its data in chunks and reports io.EOF ON THE SAME Read that
// returns the FINAL chunk. This is explicitly permitted by the io.Reader
// contract, and many real readers do it. It is the shape that breaks sum-type
// intuition: the last Read is (n > 0, io.EOF) together.
type eofReader struct {
	chunks [][]byte
	i      int
}

func (r *eofReader) Read(p []byte) (int, error) {
	if r.i >= len(r.chunks) {
		return 0, io.EOF
	}
	n := copy(p, r.chunks[r.i])
	r.i++
	if r.i == len(r.chunks) {
		return n, io.EOF // last chunk: n > 0 AND err != nil, together
	}
	return n, nil
}
```
<!-- END:snippet go-reader -->

Reading the error first (the sum-type reflex) drops those bytes:

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
CopyBuggy   got "hi"  (dropped the EOF chunk)
CopyCorrect got "hi!"
BuggyTypedNil() == nil: false  (phantom error)
CorrectNil()    == nil: true
errors.Is(Lookup(""), ErrNotFound): true
errors.As(Lookup("k"), *MyError):    true
errors.Is(Join(...), ErrNotFound):   true
AndThen(Map(Ok(3), *2), +1) = 7  (no method chaining)
```
<!-- END:output go-demo -->

`CopyBuggy` returns `"hi"`; `CopyCorrect` returns `"hi!"`. The difference is one
line: consume `n` before checking `err`. Because the two are a pair, only
convention — not the type — tells you both may be set.

## The typed-nil trap and the `#[must_use]` hole

A nil `*MyError` returned through the `error` interface is a **non-nil** error
(non-nil itab, nil data), so `err != nil` fires on success:

<!-- BEGIN:snippet go-typednil -->
```go
// BuggyTypedNil returns a *MyError typed nil. Assigning a nil *MyError to the
// error interface yields a NON-nil error (non-nil itab, nil data), so a caller's
// `if err != nil` fires even though "there is no error". The value is returned
// through the interface at the boundary, which is exactly where the trap hides.
func BuggyTypedNil() error {
	var e *MyError // nil
	return e       // becomes a non-nil error
}

// CorrectNil returns the interface nil explicitly.
func CorrectNil() error { return nil }
```
<!-- END:snippet go-typednil -->

Separately, Go's `error` carries no `#[must_use]`: ignoring it is legal. This
compiles, and only the external `errcheck` (in `golangci-lint`) objects:

<!-- BEGIN:output go-errcheck -->
```text
ignored.go:11:11: Error return value of `os.Remove` is not checked (errcheck)
```
<!-- END:output go-errcheck -->

## A lawful Result monad Go cannot use

`Result[T]` below is a genuine monad — the property tests
(`testing/quick`) confirm left identity, right identity, and associativity — and
yet it is unusable, for one precise reason:

<!-- BEGIN:snippet go-result -->
```go
type Result[T any] struct {
	val T
	err error
}

func Ok[T any](v T) Result[T]      { return Result[T]{val: v} }
func Err[T any](e error) Result[T] { return Result[T]{err: e} }

func (r Result[T]) Unwrap() (T, error) { return r.val, r.err }

// Map and AndThen are free functions, NOT methods, because they introduce a new
// type parameter U. This is what forces `AndThen(Map(r, f), g)` instead of the
// `r.map(f).and_then(g)` a Rust or OCaml engineer expects.
func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	return Ok(f(r.val))
}

func AndThen[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}
	return f(r.val)
}
```
<!-- END:snippet go-result -->

`Map` and `AndThen` change the type parameter (`T → U`), and **Go methods may
not have type parameters**. So they must be free functions, and you write
`AndThen(Map(r, f), g)` inside-out instead of `r.map(f).and_then(g)`. The
restriction is not a style choice; a generic method is a compile error:

<!-- BEGIN:output go-genericmethod-reject -->
```text
./badmethod.go:10:23: syntax error: method must have no type parameters
```
<!-- END:output go-genericmethod-reject -->

The other four languages get the chaining for free because `Result`/`throws` are
first-class and methods can be generic. See [`COMPARISON.md`](COMPARISON.md) for
`?` / `let*` / `try` / `runCatching` side by side, and [`exercises/`](exercises).
Go's own idiom — the sticky-error pattern (`bufio.Scanner`, an `errWriter`) — is
the manual `bind`: accumulate the first error in a field and check once at the
end, because the language will not thread it for you.

## References

Official sources first, grouped by language.

### Go

- `io.Reader` contract (n>0 with EOF): https://pkg.go.dev/io#Reader
- The Go Blog, "Working with Errors in Go 1.13" (`%w`, `Is`/`As`): https://go.dev/blog/go1.13-errors
- `errors` package: https://pkg.go.dev/errors
- The Go Blog, "An Introduction To Generics": https://go.dev/blog/intro-generics
- Go spec, Method declarations (methods take no type parameters): https://go.dev/ref/spec#Method_declarations
- errcheck: https://github.com/kisielk/errcheck

### Rust

- The Rust Book, Error handling (`Result`, `?`): https://doc.rust-lang.org/book/ch09-00-error-handling.html
- thiserror: https://docs.rs/thiserror

### OCaml

- OCaml Manual, Binding operators (`let*`): https://ocaml.org/manual/5.4/bindingops.html

### Swift

- Swift, Error handling (`throws`/`try`): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/errorhandling/
- Swift, `Result`: https://developer.apple.com/documentation/swift/result

### Kotlin

- Kotlin, Exceptions: https://kotlinlang.org/docs/exceptions.html
