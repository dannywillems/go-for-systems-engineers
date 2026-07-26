let () =
  let c = Capstone.create 16 4 0.0 in
  let domains =
    List.init 8 (fun w ->
        Domain.spawn (fun () ->
            for i = 0 to 499 do
              let key = (w * 500 + i) mod 100 in
              assert (Capstone.get c key = key * key)
            done))
  in
  List.iter Domain.join domains;
  assert (Capstone.length c <= 16);
  assert (Capstone.backend_calls c < 8 * 500);
  print_string "ok\n"
