let () =
  (* chunk_sum over [0,4) = 0 + 1 + sqrt 2 + sqrt 3 *)
  let expected = 0. +. 1. +. sqrt 2. +. sqrt 3. in
  assert (Float.abs (Sched.chunk_sum 0 4 -. expected) < 1e-9);
  (* parallel over the same range equals the serial sum *)
  assert (
    Float.abs (Sched.parallel_sqrt_sum 400 4 -. Sched.chunk_sum 0 400) < 1e-6);
  print_string "ok\n"
