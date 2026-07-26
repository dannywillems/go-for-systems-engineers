(* OCaml has TWO forms of genericity, neither monomorphized:
   1. Parametric polymorphism ('a list, 'a -> 'a): ONE uniform implementation,
      values boxed to a common representation. Like Go's interface path, not its
      stencil.
   2. FUNCTORS: modules parameterized by modules. A functor body is compiled
      once; applying it (Make(Int)) selects the operations. Compile-time
      genericity without code duplication. *)

(* region:functor:start *)

module type ADD = sig
  type t

  val zero : t
  val add : t -> t -> t
end

module Sum (M : ADD) = struct
  let sum xs = List.fold_left M.add M.zero xs
end

module IntSum = Sum (struct
  type t = int

  let zero = 0
  let add = ( + )
end)

module FloatSum = Sum (struct
  type t = float

  let zero = 0.0
  let add = ( +. )
end)

(* region:functor:end *)
