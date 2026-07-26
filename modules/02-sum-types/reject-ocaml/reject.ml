(* DOES NOT COMPILE: the match omits [Blue]. Warning 8 -> error. *)

type color = Red | Green | Blue

let name = function Red -> "red" | Green -> "green"

let () = print_string (name Red)
