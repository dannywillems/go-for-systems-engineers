(* OCaml's escape hatches. [Bytes.unsafe_to_string] is the sanctioned zero-copy
   bytes->string (in stdlib): no copy, but the caller must never mutate the bytes
   afterwards. [Obj.magic] is the nuclear option -- an unchecked cast between ANY
   two types that bypasses the type system entirely and is undefined behaviour if
   the representations do not match. It is used here only to demonstrate what the
   escape hatch is; real code must avoid it. *)

(* region:unsafe:start *)

(* Zero-copy: reinterpret bytes as an immutable string. Safe iff [b] is never
   mutated after this call. *)
let bytes_to_string (b : bytes) : string = Bytes.unsafe_to_string b

(* Obj.magic: an unchecked cast. Here it is a no-op int->int to show the shape;
   applied to mismatched representations it is undefined behaviour. *)
let launder (x : int) : int = Obj.magic x

(* region:unsafe:end *)
