let () =
  let e = Expr.(Add (Mul (Lit 2, Lit 3), Neg (Lit 4))) in
  assert (Expr.eval e = 2);
  assert (Expr.default 0 None = 0);
  assert (Expr.default 0 (Some 7) = 7);
  print_string "ok\n"
