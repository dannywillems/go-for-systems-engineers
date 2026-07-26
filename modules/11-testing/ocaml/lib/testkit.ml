(* Module 11 in OCaml: the same Normalize subject. Tests below use plain
   assertions; the production tools are alcotest (structured units), qcheck
   (property), and crowbar/afl (fuzzing) -- see the README. *)

let normalize s =
  s
  |> String.split_on_char ' '
  |> List.concat_map (String.split_on_char '\t')
  |> List.concat_map (String.split_on_char '\n')
  |> List.concat_map (String.split_on_char '\r')
  |> List.filter (fun w -> w <> "")
  |> String.concat " "
  |> String.lowercase_ascii
