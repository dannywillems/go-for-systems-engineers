(* t is abstract: the representation is hidden from any user of Store. *)
type t

val make : int -> t
val get : t -> int
