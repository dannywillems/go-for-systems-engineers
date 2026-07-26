let () =
  assert (Ts.(eval (Add (Int 2, Int 3))) = 5);
  assert (Ts.(eval (Eq (Int 3, Int 3))) = true);
  assert (Ts.(eval (If (Bool false, Int 1, Int 2))) = 2);
  print_string "ok\n"
