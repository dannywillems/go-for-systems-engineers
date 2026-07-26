let () =
  let shapes = [ Shapes.circle 1.0; Shapes.square 2.0 ] in
  Printf.printf "a circle object has area %.4f\n" (List.hd shapes)#area;
  Printf.printf "sum over structural shapes = %.4f\n" (Shapes.sum_objs shapes)
