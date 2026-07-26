(* OCaml's [('a, 'e) result] is a coproduct (Ok | Error). The [let*] binding
   operator (a user-defined monadic bind) gives the same left-to-right,
   early-exit ergonomics as Rust's [?], recovering do-notation for errors. *)

(* region:result:start *)

let ( let* ) = Result.bind

(* The same computation as the Go and Rust demos, in let*-notation. *)
let chain x =
  let* v = Ok (x * 2) in
  if v > 100 then Error "too big" else Ok (v + 1)

(* let* threads the error through automatically; the success path is linear. *)
let use_bind x =
  let* doubled = chain x in
  Ok (doubled + 100)

(* region:result:end *)
