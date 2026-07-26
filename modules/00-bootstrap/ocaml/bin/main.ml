(* Prints the two deterministic facts Module 00 checks. stdout is injected into
   the module README verbatim by the capture tool. *)

let n = 1_000_000

let () =
  Printf.printf "sum(1..%d) = %d\n" n (Bootstrap.sum n);
  Printf.printf "word size (bytes) = %d\n" (Bootstrap.word_size_bytes ())
