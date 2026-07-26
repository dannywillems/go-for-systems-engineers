let () =
  Printf.printf "atomic_count 8 100000 = %d (correct)\n"
    (Conc.atomic_count 8 100000)
