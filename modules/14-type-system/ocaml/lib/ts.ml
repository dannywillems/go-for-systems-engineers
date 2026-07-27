(* The star of OCaml's type system: GADTs (generalized algebraic data types). The
   constructor's return type is INDEXED by the value's type, so a single [eval]
   returns the right type per constructor -- no tags, no runtime type test, and
   the match is checked exhaustive against the index. This is expressible in
   Rust only with substantial machinery and not at all in Go's type system. *)

(* region:gadt:start *)

type _ expr =
  | Int : int -> int expr
  | Bool : bool -> bool expr
  | Add : int expr * int expr -> int expr
  | If : bool expr * 'a expr * 'a expr -> 'a expr
  | Eq : int expr * int expr -> bool expr

(* [eval : 'a expr -> 'a]. The return type varies with the constructor, tracked
   by the GADT index; the [If] branches are forced to agree, statically. *)
let rec eval : type a. a expr -> a = function
  | Int n -> n
  | Bool b -> b
  | Add (x, y) -> eval x + eval y
  | If (c, t, e) -> if eval c then eval t else eval e
  | Eq (x, y) -> eval x = eval y

(* region:gadt:end *)

(* region:eq:start *)

(* Propositional TYPE equality as a value: a term of type [(a, b) eq] is EVIDENCE
   that the types a and b are equal. Matching on [Refl] refines a = b in that
   branch, so [cast] coerces an [a] to a [b] with no Obj.magic -- the identity-
   type fragment of a dependent theory, which Go/Rust/Swift/Kotlin cannot state. *)
type (_, _) eq = Refl : ('a, 'a) eq

let cast : type a b. (a, b) eq -> a -> b = fun Refl x -> x

(* region:eq:end *)
