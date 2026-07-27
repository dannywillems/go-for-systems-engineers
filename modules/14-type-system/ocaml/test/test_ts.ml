let () =
  assert (Ts.(eval (Add (Int 2, Int 3))) = 5);
  assert (Ts.(eval (Eq (Int 3, Int 3))) = true);
  assert (Ts.(eval (If (Bool false, Int 1, Int 2))) = 2);
  (* the type-equality witness casts int->int through a proof, no Obj.magic *)
  assert (Ts.(cast (Refl : (int, int) eq) 41) = 41);
  print_string "ok\n"
