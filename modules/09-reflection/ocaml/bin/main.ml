(* Shows the hand-written generic operations. Deterministic. *)
let () =
  let open Codegen in
  let a = { name = "Ada"; age = 36 } in
  Printf.printf "to_json (hand-written): %s\n" (to_json a);
  Printf.printf "show (hand-written):    %s\n" (show a);
  Printf.printf "equal (hand-written):   %b\n" (equal a { name = "Ada"; age = 36 });
  Printf.printf "not equal after edit:   %b\n" (equal a { name = "Ada"; age = 37 })
