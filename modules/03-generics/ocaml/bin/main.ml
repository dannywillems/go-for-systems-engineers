let () =
  Printf.printf "IntSum.sum [1;2;3;4;5] = %d\n" (Gen.IntSum.sum [ 1; 2; 3; 4; 5 ]);
  Printf.printf "FloatSum.sum [1.5;2.5;3.] = %g\n"
    (Gen.FloatSum.sum [ 1.5; 2.5; 3. ])
