# 02 — Comparison: coproducts and their eliminators

**Environment.** Apple M4 Pro, macOS arm64. Go 1.26.5 (+ golang.org/x/tools
v0.48.0 for the analyzer), Rust 1.92.0, OCaml 5.4.0, Swift 6.2.3, Kotlin 2.4.10.

## Where exhaustiveness lives

| Language | Sum type            | Eliminator      | Exhaustiveness enforced by      | Missing case is |
| -------- | ------------------- | --------------- | ------------------------------- | --------------- |
| Go       | sealed interface (idiom) | `switch x.(type)` | an EXTERNAL analyzer (opt-in) | legal; silent wrong answer |
| Rust     | `enum`              | `match`         | the compiler (E0004)            | compile error   |
| OCaml    | variant             | `match`         | the compiler (warning 8)        | warning, error under `-warn-error` |
| Swift    | `enum`              | `switch`        | the compiler                    | compile error   |
| Kotlin   | `sealed` + `when`   | `when` (as expr)| the compiler                    | compile error   |

The one structural difference: in the other four, the coproduct and its total
eliminator are the same language feature, and the totality check is a theorem
the type checker proves. In Go the coproduct is *encoded* (an interface plus a
naming discipline) and the eliminator is a control-flow statement with no
totality obligation, so the check has to be reconstructed by a separate program
walking the AST and the type graph — which is exactly what
[`go/exhaustive`](go/exhaustive) does.

## The five definitions

<!-- BEGIN:snippet go-variants -->
```go
type Lit struct{ V int }

type Add struct{ L, R Expr }

type Mul struct{ L, R Expr }

type Neg struct{ X Expr }

func (Lit) exprNode() {}
func (Add) exprNode() {}
func (Mul) exprNode() {}
func (Neg) exprNode() {}
```
<!-- END:snippet go-variants -->

<!-- BEGIN:snippet rust-enum -->
```rust
pub enum Expr {
    Lit(i64),
    Add(Box<Expr>, Box<Expr>),
    Mul(Box<Expr>, Box<Expr>),
    Neg(Box<Expr>),
}

impl Expr {
    pub fn eval(&self) -> i64 {
        match self {
            Expr::Lit(v) => *v,
            Expr::Add(l, r) => l.eval() + r.eval(),
            Expr::Mul(l, r) => l.eval() * r.eval(),
            Expr::Neg(x) => -x.eval(),
        }
    }
}
```
<!-- END:snippet rust-enum -->

<!-- BEGIN:snippet ocaml-variant -->
```ocaml
type expr =
  | Lit of int
  | Add of expr * expr
  | Mul of expr * expr
  | Neg of expr

let rec eval = function
  | Lit v -> v
  | Add (l, r) -> eval l + eval r
  | Mul (l, r) -> eval l * eval r
  | Neg x -> -eval x
```
<!-- END:snippet ocaml-variant -->

<!-- BEGIN:snippet swift-enum -->
```swift
public indirect enum Expr {
  case lit(Int)
  case add(Expr, Expr)
  case mul(Expr, Expr)
  case neg(Expr)
}

public func eval(_ e: Expr) -> Int {
  switch e {
  case .lit(let v): return v
  case .add(let l, let r): return eval(l) + eval(r)
  case .mul(let l, let r): return eval(l) * eval(r)
  case .neg(let x): return -eval(x)
  }
}
```
<!-- END:snippet swift-enum -->

<!-- BEGIN:snippet kotlin-sealed -->
```kotlin
sealed interface Expr

data class Lit(
    val v: Int,
) : Expr

data class Add(
    val l: Expr,
    val r: Expr,
) : Expr

data class Mul(
    val l: Expr,
    val r: Expr,
) : Expr

data class Neg(
    val x: Expr,
) : Expr

fun eval(e: Expr): Int =
    when (e) {
        is Lit -> e.v
        is Add -> eval(e.l) + eval(e.r)
        is Mul -> eval(e.l) * eval(e.r)
        is Neg -> -eval(e.x)
    }
```
<!-- END:snippet kotlin-sealed -->

## Option / null: the one-variant sum

`Option<T> = None + Some T`. How each language spells "maybe a T", and the size
cost:

| Language | "maybe T"      | Representation of `maybe pointer`                        |
| -------- | -------------- | -------------------------------------------------------- |
| Go       | `*T` or `(T, bool)` | `*T` is 8 bytes; nil is a valid value the compiler never forces you to check |
| Rust     | `Option<T>`    | `Option<&T>` / `Option<Box<T>>` niche-packed to 8 bytes; `match` forces the None case |
| OCaml    | `'a option`    | a boxed `Some` cell (a real constructor); `None` is immediate |
| Swift    | `Optional<T>`  | `Optional` is an enum; pointer payloads niche-packed; `if let`/`guard` force the check |
| Kotlin   | `T?`           | nullable reference (no wrapper allocation); the compiler tracks nullability and forces a check |

The Rust niche numbers are measured in `rust-demo` above: `Option<&u8>` and
`Option<Box<u8>>` are 8 bytes (same as the pointer), while `Option<u8>` is 2
(no spare bit pattern, so it pays a tag byte). Go's `*T` is also 8 bytes, but
the crucial difference is not size: it is that `nil` is a legal inhabitant the
type system never makes you handle, whereas the other four turn "absent" into a
constructor a total match must cover. Kotlin is the interesting middle: `T?` has
no wrapper (like Go's nil) but the compiler enforces the null check (like the
sum-typed languages).

## Reading

Go's design trades a compile-time totality guarantee for a simpler type system
(no coproducts, one fewer kind of type to teach) and gets, in exchange, a class
of "forgot a case" bugs that are silent at runtime. The analyzer recovers most
of the guarantee at the cost of a build-time tool and a `//sumtype:decl`
convention; it cannot recover the parts that depend on whole-program knowledge
the way a closed `enum` gives the compiler for free.
