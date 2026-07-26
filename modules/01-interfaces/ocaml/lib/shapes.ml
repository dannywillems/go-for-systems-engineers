(* OCaml has TWO answers to Go interfaces, at opposite ends of the
   nominal/structural axis:

   1. Objects (below) are STRUCTURAL, exactly like Go interfaces: any object
      with an [area : float] method has type [< area : float >], with no
      declared "implements". Row polymorphism generalizes this.
   2. First-class modules are EXPLICIT existentials: [(module S : SHAPE)] packs
      a module behind a signature; you unpack to use it. This is the honest
      existential Go's interface hides behind a two-word value.

   Objects are shown here because they are the direct structural analog. *)

(* region:obj:start *)

type shape = < area : float >

let circle r : shape =
  object
    method area = Float.pi *. r *. r
  end

let square s : shape =
  object
    method area = s *. s
  end

(* region:obj:end *)

let sum_objs (xs : shape list) =
  List.fold_left (fun acc s -> acc +. s#area) 0. xs
