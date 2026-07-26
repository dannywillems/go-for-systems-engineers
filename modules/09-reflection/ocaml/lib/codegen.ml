(* Module 09 in OCaml: the odd one out. OCaml's stdlib has NEITHER runtime
   reflection (no way to walk a value's type at run time) NOR a built-in derive.
   You either HAND-WRITE the generic operation, as below, or reach for a PPX
   (ppx_deriving: [@@deriving show, eq, yojson]) -- a syntactic macro that
   rewrites the AST at compile time, the same compile-time-codegen family as
   Rust's derive and Swift's synthesis, just supplied by a preprocessor. This
   module hand-writes to stay dependency-free; see the README for the ppx form. *)

type person = { name : string; age : int }

let to_json p = Printf.sprintf {|{"name":%S,"age":%d}|} p.name p.age
let equal a b = a.name = b.name && a.age = b.age
let show p = Printf.sprintf "{ name = %S; age = %d }" p.name p.age
