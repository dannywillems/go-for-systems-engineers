(* DOES NOT COMPILE. Store.t is abstract (see store.mli), so its record field is
   invisible outside the module and reaching for it is rejected:

     Error: Unbound record field Store.n

   The client cannot reach into an abstract type's representation. *)

let () =
  let s = Store.make 41 in
  Printf.printf "%d\n" s.Store.n
