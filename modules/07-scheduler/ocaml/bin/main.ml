let () =
  let total = 400_000_000 in
  let workers = Domain.recommended_domain_count () in
  let t0 = Unix.gettimeofday () in
  let acc = Sched.parallel_sqrt_sum total workers in
  ignore (Sys.opaque_identity acc);
  let ms = int_of_float ((Unix.gettimeofday () -. t0) *. 1000.) in
  Printf.printf "OCaml  sqrt-sum %dM / %d domains: %d ms\n"
    (total / 1_000_000) workers ms
