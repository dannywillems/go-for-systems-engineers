let () =
  let b = Bytes.of_string "hello" in
  Printf.printf "bytes_to_string = %s (no copy)\n" (Un.bytes_to_string b);
  Printf.printf "Obj.magic launder 42 = %d\n" (Un.launder 42)
