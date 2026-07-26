let () =
  assert (Conc.atomic_count 8 100000 = 800000);
  assert (Conc.atomic_count 4 1000 = 4000);
  print_string "ok\n"
