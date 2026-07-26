# 01 — Comparison: dispatch across five type systems

**Environment.** Apple M4 Pro, macOS arm64. Go 1.26.5, Rust 1.92.0, OCaml
5.4.0, Swift 6.2.3, Kotlin 2.4.10 (OpenJDK 26). The Go timings are in
[`go/bench.txt`](go/bench.txt); the other languages are compared on *mechanism*
here (a cross-language dispatch microbenchmark lives in the capstone, where the
workload is realistic rather than a single multiply).

## The two axes

Every language answers "how does one value stand for many types?" on two axes:
**nominal vs structural** (does satisfaction require a declaration?) and
**universal vs existential** (is the concrete type erased at runtime?).

| Language | Abstraction | Satisfaction | Static path | Dynamic path (existential) | Orphan rule |
| -------- | ----------- | ------------ | ----------- | -------------------------- | ----------- |
| Go       | `interface` | structural   | devirtualized calls | `interface` value: `(itab, data)`, 2 words | none |
| Rust     | `trait`     | nominal (`impl`) | `<T: Trait>` monomorphized | `dyn Trait`: fat ptr `(data, vtable)`, 2 words | yes (coherence) |
| OCaml    | object type / module | structural (objects) / nominal (modules) | inlined method / functor | object: `(obj, method table)`; `(module S)` packs an existential | none |
| Swift    | `protocol`  | nominal (+ retroactive) | `some`/`<S>` specialized | `any Shape`: boxed existential container | none (warns) |
| Kotlin   | `interface` | nominal      | monomorphic call site | JVM `invokeinterface` on the object header | n/a (no external impls) |

The Go, Rust, OCaml-object, and Swift dynamic paths are the *same idea* — an
existential quantifier reified as a value carrying its own witness table — and
Go/Rust make it exactly two machine words. OCaml objects are the other
structural system here and, like Go, have no orphan rule. Swift and Kotlin are
nominal like Rust, but neither enforces Rust's coherence: Swift allows
retroactive conformance (warning only), and Kotlin simply forbids adding an
interface to a type you do not own, so the ambiguity Rust's rule prevents cannot
arise a different way.

## The five interface definitions

<!-- BEGIN:snippet go-iface -->
```go
// Shape is satisfied structurally: any type with an Area() float64 method is a
// Shape, with no "implements" clause anywhere. This is the width-subtyping
// coercion X <: Shape done by structure, not by nominal declaration.
type Shape interface {
	Area() float64
}

type Circle struct{ R float64 }

func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

type Square struct{ S float64 }

func (s Square) Area() float64 { return s.S * s.S }
```
<!-- END:snippet go-iface -->

<!-- BEGIN:snippet rust-trait -->
```rust
/// A trait is satisfied only by an explicit `impl` (nominal), unlike Go's
/// structural satisfaction. Coherence (the orphan rule) further restricts WHERE
/// that impl may live: see the `reject-orphan` crate for the rejection.
pub trait Shape {
    fn area(&self) -> f64;
}

pub struct Circle {
    pub r: f64,
}
impl Shape for Circle {
    fn area(&self) -> f64 {
        std::f64::consts::PI * self.r * self.r
    }
}

pub struct Square {
    pub s: f64,
}
impl Shape for Square {
    fn area(&self) -> f64 {
        self.s * self.s
    }
}
```
<!-- END:snippet rust-trait -->

<!-- BEGIN:snippet ocaml-obj -->
```ocaml
type shape = < area : float >

let circle r : shape =
  object
    method area = Float.pi *. r *. r
  end

let square s : shape =
  object
    method area = s *. s
  end
```
<!-- END:snippet ocaml-obj -->

<!-- BEGIN:snippet swift-proto -->
```swift
public protocol Shape {
  func area() -> Double
}

public struct Circle: Shape {
  public let r: Double
  public init(r: Double) { self.r = r }
  public func area() -> Double { Double.pi * r * r }
}

public struct Square: Shape {
  public let s: Double
  public init(s: Double) { self.s = s }
  public func area() -> Double { s * s }
}
```
<!-- END:snippet swift-proto -->

<!-- BEGIN:snippet kotlin-iface -->
```kotlin
interface Shape {
    fun area(): Double
}

class Circle(
    private val r: Double,
) : Shape {
    override fun area(): Double = Math.PI * r * r
}

class Square(
    private val s: Double,
) : Shape {
    override fun area(): Double = s * s
}
```
<!-- END:snippet kotlin-iface -->

## Dynamic-dispatch outputs

<!-- BEGIN:output rust-demo -->
```text
$ cargo run --quiet --bin demo
a Circle used as &dyn Shape has area 3.1416
sizeof(&dyn Shape) = 16 bytes
sizeof(&Circle)    = 8 bytes
sum via dyn dispatch = 7.1416
```
<!-- END:output rust-demo -->

<!-- BEGIN:output ocaml-demo -->
```text
$ dune exec bin/main.exe
a circle object has area 3.1416
sum over structural shapes = 7.1416
```
<!-- END:output ocaml-demo -->

<!-- BEGIN:output swift-demo -->
```text
a Circle as any Shape has area 3.1416
sum via existential dispatch = 7.1416
```
<!-- END:output swift-demo -->

<!-- BEGIN:output kotlin-demo -->
```text
a Circle used as a Shape has area 3.1416
sum via interface dispatch = 7.1416
```
<!-- END:output kotlin-demo -->

## Reading

The theoretical unification: universal quantification (`∀X. ...`, Go's concrete
call, Rust's `<T>`, a specialized generic) is eliminated by the compiler and
costs nothing at runtime; existential quantification (`∃X. ...`, every "dynamic"
column) survives as a value pairing data with a witness table and pays an
indirect call plus, usually, lost inlining. Go's distinctive choice is to make
the *existential* implicit and structural — you never write `∃` or `impl`, the
`itab` appears wherever a concrete value flows into an interface — which is why
it has no coherence rule to enforce and no place to hang one.
