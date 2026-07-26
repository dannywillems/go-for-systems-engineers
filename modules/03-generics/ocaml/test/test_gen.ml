let () =
  assert (Gen.IntSum.sum [ 1; 2; 3; 4; 5 ] = 15);
  assert (Gen.FloatSum.sum [ 1.5; 2.5; 3. ] = 7.0);
  print_string "ok\n"
