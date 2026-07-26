(* The Module 00 fixture in OCaml: the same trivial computation as the Go and
   Rust sides, so the three demo binaries emit byte-identical output. *)

(* region:sum:start *)

(** [sum n] returns 1 + 2 + ... + n. Identical on every 64-bit target and in
    every language, which makes it a clean cross-toolchain fixture. *)
let sum n =
  let total = ref 0 in
  for i = 1 to n do
    total := !total + i
  done;
  !total

(** Native word size in bytes: 8 on any 64-bit platform. Note OCaml's [int] is
    63-bit (one tag bit), but the machine word it lives in is still 8 bytes. *)
let word_size_bytes () = Sys.word_size / 8

(* region:sum:end *)
