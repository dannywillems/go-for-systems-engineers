let capacity = 256
let max_inflight = 32
let keys = 256
let workers = 64
let per_worker = 10_000

let () =
  let c = Capstone.create capacity max_inflight 0.0001 in
  let t0 = Unix.gettimeofday () in
  let domains =
    List.init workers (fun w ->
        Domain.spawn (fun () ->
            let lat = Array.make per_worker 0.0 in
            let seed = ref (Int64.logor (Int64.mul (Int64.of_int w) 2654435761L) 1L) in
            for i = 0 to per_worker - 1 do
              seed :=
                Int64.add
                  (Int64.mul !seed 6364136223846793005L)
                  1442695040888963407L;
              let key =
                Int64.to_int
                  (Int64.rem (Int64.shift_right_logical !seed 33) (Int64.of_int keys))
              in
              let s = Unix.gettimeofday () in
              ignore (Capstone.get c (abs key));
              lat.(i) <- Unix.gettimeofday () -. s
            done;
            lat))
  in
  let all =
    List.concat_map (fun d -> Array.to_list (Domain.join d)) domains
    |> Array.of_list
  in
  let elapsed = Unix.gettimeofday () -. t0 in
  Array.sort compare all;
  let n = Array.length all in
  let pc p = all.(min (int_of_float (p /. 100. *. float_of_int n)) (n - 1)) in
  let us s = s *. 1_000_000. in
  Printf.printf
    "OCaml  %dk gets/%dw: %.0f kops/s  p50=%.1fus p99=%.1fus p999=%.1fus  backend=%.1f%% of gets\n"
    (n / 1000) workers
    (float_of_int n /. elapsed /. 1000.)
    (us (pc 50.)) (us (pc 99.)) (us (pc 99.9))
    (100. *. float_of_int (Capstone.backend_calls c) /. float_of_int n)
