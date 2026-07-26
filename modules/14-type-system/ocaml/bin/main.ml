let () =
  (* if (2 = 2) then (1 + 3) else 0  -- an int expr, so eval returns an int *)
  let e = Ts.(If (Eq (Int 2, Int 2), Add (Int 1, Int 3), Int 0)) in
  Printf.printf "eval (if 2=2 then 1+3 else 0) = %d\n" (Ts.eval e);
  (* eval on a bool expr returns a bool, from the SAME function *)
  Printf.printf "eval (3 = 4) = %b\n" Ts.(eval (Eq (Int 3, Int 4)))
