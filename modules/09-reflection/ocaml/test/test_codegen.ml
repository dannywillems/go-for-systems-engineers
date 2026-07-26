let () =
  let open Codegen in
  let a = { name = "Ada"; age = 36 } in
  assert (to_json a = {|{"name":"Ada","age":36}|});
  assert (equal a { name = "Ada"; age = 36 });
  assert (not (equal a { name = "Ada"; age = 37 }));
  print_string "ok\n"
