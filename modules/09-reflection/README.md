# 09 — Reflection & code generation

**Thesis.** "Operate generically over a type you did not hand-write for"
(serialize, compare, print, hash) is resolved at one of two times. Go resolves
it at RUN TIME with reflection: the type descriptor travels inside the value and
`reflect` walks it dynamically, so any type works with no build step, but every
call pays the walk and every mistake is deferred to run time. The other four
resolve it at COMPILE TIME by generating specialized code per type — Rust's
`#[derive]` (a procedural macro), Swift's protocol synthesis, Kotlin's
`data class`, and OCaml's ppx (or, lacking one, a hand-written impl). The
consequence that matters most is not speed but ERROR TIMING: codegen rejects a
bad type at compile time; reflection accepts anything and fails, or silently does
nothing, at run time.

## Contents

- [Go: runtime reflection](#go-runtime-reflection)
- [The four compile-time mechanisms](#the-four-compile-time-mechanisms)
- [Error timing: the real divergence](#error-timing-the-real-divergence)
- [The measured cost of reflection](#the-measured-cost-of-reflection)
- [Exercises](#exercises)
- [References](#references)

## Go: runtime reflection

A struct carries TAGS that reflection reads at run time; `encoding/json` has
never seen the type at compile time and discovers its fields by walking
`reflect.Type` when `Marshal` is called:

<!-- BEGIN:snippet go-tags -->
```go
// Person carries struct TAGS that reflection reads at run time. encoding/json
// never sees this type at compile time; it discovers Name/Age and their tags by
// walking reflect.Type when Marshal is called. secret is unexported, so
// reflection from another package cannot see it at all -- it is silently omitted
// from the JSON, with no error.
//
//nolint:govet // field order follows the demo/JSON narrative, not byte alignment
type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	secret string // unexported: invisible to reflection, silently dropped
}
```
<!-- END:snippet go-tags -->

The same mechanism powers a function that inspects ANY struct, including types it
was never written for — the thing reflection buys, and the thing compile-time
codegen cannot do without being asked for each type in advance:

<!-- BEGIN:snippet go-reflect -->
```go
// Describe is what reflection buys: one function that inspects the fields of ANY
// struct, including types it has never heard of, with no per-type code. It reads
// the reflect.Type at run time. This is impossible with compile-time codegen
// alone -- the derive/synthesis approaches must be asked for each type in
// advance -- but every call pays the reflection walk, and a non-struct argument
// is a run-time error, not a compile-time one.
func Describe(v any) []string {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if t.Kind() != reflect.Struct {
		return []string{fmt.Sprintf("<not a struct: %s>", t.Kind())}
	}
	var out []string
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() { // reflection sees the field but cannot read its value
			out = append(out, fmt.Sprintf("%s %s (unexported)", f.Name, f.Type))
			continue
		}
		tag := f.Tag.Get("json")
		out = append(out, fmt.Sprintf("%s %s json=%q = %v",
			f.Name, f.Type, tag, val.Field(i).Interface()))
	}
	return out
}
```
<!-- END:snippet go-reflect -->

<!-- BEGIN:output go-demo -->
```text
$ go run ./cmd/demo
json.Marshal (reflection): {"name":"Ada","age":36}
secret field still set:    "top-secret" (dropped from JSON, no error)
Describe (reflection walk):
  Name string json="name" = Ada
  Age int json="age" = 36
  secret string (unexported)
round-trip: Ada is 36
```
<!-- END:output go-demo -->

Note the silently dropped `secret` field: reflection from `encoding/json` cannot
see an unexported field, so it is omitted with no error — the first taste of
run-time deferral.

## The four compile-time mechanisms

Each of the other four generates the operation at compile time, and produces the
same kind of result through a different machine. Rust's `#[derive]` is a
procedural macro that emits a concrete `impl` for exactly this type:

<!-- BEGIN:snippet rust-derive -->
```rust
// The compiler generates Debug (a formatter), Clone, and PartialEq (a
// field-by-field ==) for this exact type. Each is real generated code, checked
// against every field's type at compile time.
#[derive(Debug, Clone, PartialEq)]
pub struct Person {
    pub name: String,
    pub age: u32,
}
```
<!-- END:snippet rust-derive -->

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
Debug (derived):      Person { name: "Ada", age: 36 }
PartialEq (derived):  true
after clone+edit:     Person { name: "Ada", age: 37 }
still equal?          false
```
<!-- END:output rust-demo -->

Swift's compiler SYNTHESIZES a protocol conformance (here `Equatable`)
member-by-member when you declare it:

<!-- BEGIN:output swift-demo -->
```text
describe:              Ada is 36
Equatable (synth):     true
not equal after edit:  false
```
<!-- END:output swift-demo -->

Kotlin's `data class` generates `equals`, `hashCode`, `toString`, `copy`, and
`componentN` from the primary-constructor properties:

<!-- BEGIN:output kotlin-demo -->
```text
toString (generated): Person(name=Ada, age=36)
equals (generated):   true
copy (generated):     Person(name=Ada, age=37)
not equal after copy: false
```
<!-- END:output kotlin-demo -->

OCaml is the outlier: its standard library has NEITHER runtime reflection NOR a
built-in derive, so you hand-write the operation (below) or reach for a ppx —
`[@@deriving show, eq, yojson]` — a syntactic macro that rewrites the AST at
compile time, the same family as Rust's derive:

<!-- BEGIN:output ocaml-demo -->
```text
$ dune exec bin/main.exe
to_json (hand-written): {"name":"Ada","age":36}
show (hand-written):    { name = "Ada"; age = 36 }
equal (hand-written):   true
not equal after edit:   false
```
<!-- END:output ocaml-demo -->

## Error timing: the real divergence

Because the compile-time mechanisms check the type as they generate code, a type
they cannot handle is a COMPILE error. Reflection has no such gate — its
equivalent guard is a separate static analyzer (`go vet`) bolted on. The four
rejects:

**Go** — a malformed struct tag compiles fine (reflection reads tags as opaque
strings at run time); `go vet`'s structtag check is what catches it:

<!-- BEGIN:output go-reject -->
```text
bad.go:15:2: struct field tag `json:name` not compatible with reflect.StructTag.Get: bad syntax for struct tag value
```
<!-- END:output go-reject -->

**Rust** — `#[derive(Debug)]` requires every field to be `Debug`, checked when
the macro expands:

<!-- BEGIN:output rust-reject -->
```text
error[E0277]: `NotDebug` doesn't implement `Debug`
```
<!-- END:output rust-reject -->

**Swift** — `Equatable` synthesis requires every stored property to be
`Equatable`:

<!-- BEGIN:output swift-reject -->
```text
reject.swift:8:8: error: type 'Wrapper' does not conform to protocol 'Equatable'
```
<!-- END:output swift-reject -->

**Kotlin** — a `data class` exists to generate members from its properties, so
one with none is rejected:

<!-- BEGIN:output kotlin-reject -->
```text
reject.kt:6:17: error: data class must have at least one primary constructor parameter.
```
<!-- END:output kotlin-reject -->

The Go case is the telling one: the bad tag is not an error to the compiler at
all. Without `go vet` it ships, and `json.Marshal` simply ignores the tag at run
time. That is the price of resolving genericity dynamically — the language cannot
check what it never looks at until the value flows through.

## The measured cost of reflection

The naive claim "reflection is slow" is mostly wrong, and the measurement says
why: `encoding/json` does NOT re-walk the type on every call — it builds a
per-type encoder once and caches it (a `sync.Map` of type to encoder function),
amortizing the reflection walk. Even so, the cached reflective path is still
several times a specialized hand-written encoder, because it works through
reflect-derived closures and allocates more:

<!-- BEGIN:file measured -->
```text
machine: Apple M4 Pro (14 cores), macOS arm64
workload: marshal a 2-field struct to JSON, two ways producing the same
bytes. ReflectMarshal = encoding/json (walks reflect.Type per call);
ManualMarshal = a hand-written encoder. allocs/op is deterministic.

# Go: reflection vs hand-written marshaler (go test -bench -benchmem)
BenchmarkReflectMarshal-14    	12301630	        88.21 ns/op	      96 B/op	       3 allocs/op
BenchmarkManualMarshal-14     	75106011	        14.78 ns/op	      24 B/op	       1 allocs/op
```
<!-- END:file measured -->

Read it honestly: reflection here is ~6x the hand-written encoder and allocates
3x as often — a real, measurable cost, but a constant-factor one that is fine for
most call sites and only matters on a hot path. Compile-time codegen (Rust
`derive`, Swift synthesis) generates the specialized path directly, so it pays
neither the walk nor the indirection; that, plus compile-time error checking, is
what you buy with a build step. `benchstat` turns repeated runs into a mean with
a confidence interval; a single ns/op is noise.

## Exercises

[`exercises/go`](exercises/go) is red until you implement a reflection-based
`ToMap(any) map[string]any` that works on any struct (exported fields only, honor
`json` tags for the key). [`solutions/go`](solutions/go) is the verified corrigé:

```
make exercises M=09   # red
make solutions M=09   # green
```

## References

Official sources first, grouped by language.

### Go

- `reflect`: https://pkg.go.dev/reflect
- The Laws of Reflection (the Go blog): https://go.dev/blog/laws-of-reflection
- `encoding/json` (reflection-based, cached encoders): https://pkg.go.dev/encoding/json
- `go generate` (codegen as external tooling): https://go.dev/blog/generate
- `go vet` structtag check: https://pkg.go.dev/cmd/vet

### Rust

- The Rust Reference, "Derive": https://doc.rust-lang.org/reference/attributes/derive.html
- Procedural macros: https://doc.rust-lang.org/reference/procedural-macros.html
- serde (the canonical derive-based serializer): https://serde.rs/

### OCaml

- ppx_deriving (compile-time deriving via PPX): https://github.com/ocaml-ppx/ppx_deriving
- ppxlib (the PPX infrastructure): https://github.com/ocaml-ppx/ppxlib

### Swift

- The Swift Programming Language, "Protocols" (synthesized conformances): https://docs.swift.org/swift-book/documentation/the-swift-programming-language/protocols/
- SE-0185, Synthesizing Equatable and Hashable: https://github.com/apple/swift-evolution/blob/main/proposals/0185-synthesize-equatable-hashable.md
- Codable: https://developer.apple.com/documentation/swift/codable

### Kotlin (JVM)

- Kotlin docs, "Data classes": https://kotlinlang.org/docs/data-classes.html
- kotlinx.serialization (compiler-plugin codegen): https://github.com/Kotlin/kotlinx.serialization
