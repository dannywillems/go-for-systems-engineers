(* Allocate 50M short-lived 8-int arrays. OCaml's generational GC bump-allocates
   on the minor heap and reclaims short-lived garbage cheaply, which is often
   faster than Go's non-generational GC for exactly this workload. Measured. *)

let () =
  let n = 50_000_000 in
  let acc = ref 0 in
  let t0 = Unix.gettimeofday () in
  for i = 1 to n do
    let a = Array.make 8 0 in
    a.(0) <- i;
    acc := !acc + a.(0)
  done;
  let ms = int_of_float ((Unix.gettimeofday () -. t0) *. 1000.) in
  Printf.printf "OCaml alloc 50M (minor GC): %d ms (acc=%d)\n" ms !acc
