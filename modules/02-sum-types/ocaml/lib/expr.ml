(* OCaml has coproducts as its bread and butter: a variant type IS a sum, and
   [match] is a total eliminator. A non-exhaustive match is warning 8, promoted
   to an error here (and in most CI). See the reject-ocaml project. *)

(* region:variant:start *)

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

(* region:variant:end *)

(* [option] is itself a sum type ([None | Some of 'a]), so OCaml needs no
   nullable-pointer hack: the absence case is a distinct constructor the match
   must handle. *)
let default d = function None -> d | Some x -> x
