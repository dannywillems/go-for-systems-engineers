let () =
  (* (2 * 3) + (-4) = 2 *)
  let e = Expr.(Add (Mul (Lit 2, Lit 3), Neg (Lit 4))) in
  Printf.printf "eval((2*3) + -4) = %d\n" (Expr.eval e);
  Printf.printf "default 0 None = %d\n" (Expr.default 0 None);
  Printf.printf "default 0 (Some 7) = %d\n" (Expr.default 0 (Some 7))
