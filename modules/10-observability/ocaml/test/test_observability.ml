(* Red/green gate: the two builders agree, and the pre-sized Buffer allocates
   strictly fewer minor-heap words than folding (^). *)

let () =
  let parts = List.init 64 (fun i -> Printf.sprintf "chunk-%02d;" i) in
  assert (Observability.concat_caret parts = Observability.buffer_build parts);
  let caret_w = Observability.words_allocated (fun () -> Observability.concat_caret parts) in
  let buf_w = Observability.words_allocated (fun () -> Observability.buffer_build parts) in
  assert (buf_w < caret_w);
  print_string "ok\n"
