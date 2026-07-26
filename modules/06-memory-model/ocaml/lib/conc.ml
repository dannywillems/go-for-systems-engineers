(* OCaml 5 introduced DOMAINS (true parallelism) and [Atomic] for lock-free
   shared state. Like Go, and unlike Rust/Swift 6, there is no compile-time
   Send/Sync check: unsynchronized sharing across domains is a race with
   undefined behavior. The idiomatic safe primitive is [Atomic]; the runtime and
   its data-race tooling are new (OCaml 5.0, 2022). *)

(* region:atomic:start *)

let atomic_count threads per =
  let counter = Atomic.make 0 in
  let domains =
    List.init threads (fun _ ->
        Domain.spawn (fun () ->
            for _ = 1 to per do
              Atomic.incr counter
            done))
  in
  List.iter Domain.join domains;
  Atomic.get counter

(* region:atomic:end *)
