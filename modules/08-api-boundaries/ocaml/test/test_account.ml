(* Exercises only the public signature -- the client never sees the record. *)

let () =
  let open Account in
  let a = open_ 100 in
  let a = deposit a 50 in
  assert (balance a = 150);
  (match withdraw a 200 with _ -> assert false | exception Overdraft -> ());
  (match open_ (-1) with _ -> assert false | exception Overdraft -> ());
  print_string "ok\n"
