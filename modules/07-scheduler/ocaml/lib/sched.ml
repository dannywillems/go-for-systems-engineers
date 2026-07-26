(* OCaml 5 parallelism is via DOMAINS: 1:1 with OS threads (unlike Go's M:N
   goroutines). Structured concurrency and cooperative I/O concurrency come from
   Eio (effect-based), a separate library; for a pure CPU sweep, domains are the
   direct comparison. *)

let chunk_sum lo hi =
  let s = ref 0.0 in
  for i = lo to hi - 1 do
    s := !s +. sqrt (float_of_int i)
  done;
  !s

let parallel_sqrt_sum total workers =
  let chunk = total / workers in
  let domains =
    List.init workers (fun k ->
        Domain.spawn (fun () -> chunk_sum (k * chunk) ((k + 1) * chunk)))
  in
  List.fold_left (fun acc d -> acc +. Domain.join d) 0.0 domains
