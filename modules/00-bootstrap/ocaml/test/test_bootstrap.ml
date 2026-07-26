(* Minimal assertion-based test; `dune runtest` fails if any assert trips. *)

let () =
  List.iter
    (fun n -> assert (Bootstrap.sum n = n * (n + 1) / 2))
    [ 0; 1; 10; 1_000_000 ];
  assert (Bootstrap.word_size_bytes () = 8);
  print_string "ok\n"
