let () =
  let shapes = [ Shapes.circle 1.0; Shapes.circle 2.0 ] in
  let expected = (Float.pi *. 1.0) +. (Float.pi *. 4.0) in
  assert (Float.abs (Shapes.sum_objs shapes -. expected) < 1e-9);
  (* Structural: an anonymous object with an area method is a shape too. *)
  let anon : Shapes.shape =
    object
      method area = 10.0
    end
  in
  assert (Float.abs (anon#area -. 10.0) < 1e-9);
  print_string "ok\n"
