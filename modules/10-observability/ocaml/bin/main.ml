(* Prints minor-heap words allocated by each builder, plus a wall-clock timing
   for measured.txt. The word counts are deterministic; the timings are not. *)

let parts = List.init 64 (fun i -> Printf.sprintf "chunk-%02d;" i)

let time_ns f iters =
  let t0 = Unix.gettimeofday () in
  for _ = 1 to iters do
    ignore (Sys.opaque_identity (f ()))
  done;
  (Unix.gettimeofday () -. t0) *. 1e9 /. float_of_int iters

let () =
  let caret_w = Observability.words_allocated (fun () -> Observability.concat_caret parts) in
  let buf_w = Observability.words_allocated (fun () -> Observability.buffer_build parts) in
  Printf.printf "concat_caret (64 parts): %.0f minor words\n" caret_w;
  Printf.printf "buffer_build (64 parts): %.0f minor words\n" buf_w;
  let iters = 200_000 in
  let caret_ns = time_ns (fun () -> Observability.concat_caret parts) iters in
  let buf_ns = time_ns (fun () -> Observability.buffer_build parts) iters in
  Printf.printf
    "OCaml concat_caret: %.0f ns/op (%.0f words)  buffer_build: %.0f ns/op (%.0f words)\n"
    caret_ns caret_w buf_ns buf_w
