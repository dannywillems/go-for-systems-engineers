let () =
  assert (Errs.chain 3 = Ok 7);
  assert (Errs.chain 60 = Error "too big");
  assert (Errs.use_bind 3 = Ok 107);
  print_string "ok\n"
