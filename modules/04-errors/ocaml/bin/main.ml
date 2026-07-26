let show = function Ok v -> Printf.sprintf "Ok %d" v | Error e -> "Error " ^ e

let () =
  Printf.printf "chain 3 = %s  (let* notation)\n" (show (Errs.chain 3));
  Printf.printf "chain 60 = %s\n" (show (Errs.chain 60));
  Printf.printf "use_bind 3 = %s\n" (show (Errs.use_bind 3))
