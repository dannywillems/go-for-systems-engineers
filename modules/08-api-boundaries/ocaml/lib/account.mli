(* The .mli is OCaml's encapsulation boundary, and it is the strongest of the
   five: `type t` here is ABSTRACT -- the signature exposes the name but not the
   representation, so no client can see that t is a record, build one with a
   literal, or pattern-match its fields. Only the values listed below exist
   outside this module. This is signature ascription: the .ml is checked AGAINST
   this .mli, and anything not re-exported here is invisible. *)

(* region:mli:start *)

type t
(** an account; representation hidden *)

exception Overdraft

val open_ : int -> t
(** [open_ initial] raises Overdraft if initial < 0 *)

val deposit : t -> int -> t
val withdraw : t -> int -> t
val balance : t -> int

(* region:mli:end *)
