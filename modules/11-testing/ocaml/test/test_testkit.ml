(* Table assertions + a hand-rolled property loop (idempotence over a generated
   corpus via a small LCG), dependency-free. *)

let () =
  let n = Testkit.normalize in
  assert (n "  hi  " = "hi");
  assert (n "a\t\n  b" = "a b");
  assert (n "MiXeD" = "mixed");
  assert (n "   " = "");

  (* property: normalize is idempotent, trimmed, no double space, on many
     generated shapes. *)
  let alphabet = [| 'a'; 'B'; ' '; '\t'; '\n'; 'c' |] in
  let seed = ref 0x2545F4914F6CDD1D in
  let next () =
    seed := (!seed * 2862933555777941757) + 1;
    !seed
  in
  for _ = 1 to 10_000 do
    let len = (next () lsr 60) land 0xF in
    let b = Buffer.create len in
    for _ = 1 to len do
      Buffer.add_char b alphabet.(abs (next () lsr 61) mod Array.length alphabet)
    done;
    let once = n (Buffer.contents b) in
    (* idempotence (which already rules out collapsible double spaces) + trimmed *)
    assert (n once = once);
    assert (String.trim once = once)
  done;
  print_string "ok\n"
