(* Module 10 in OCaml: measuring allocation with the runtime's own counters.
   Gc.minor_words () returns the total words allocated in the minor heap so far,
   so the difference across an operation is exactly what it allocated -- the same
   falsifiable spine as the Go module, no external profiler.

   Like Rust (and unlike Go's immutable strings), OCaml's Buffer is a growable
   sink, so the contrast is naive intermediate strings via (^) versus one
   pre-sized Buffer. *)

(* concat_caret folds (^) over the parts, allocating a fresh intermediate string
   at every step: O(n) intermediate strings, O(n^2) bytes copied. *)
let concat_caret parts = List.fold_left (fun acc p -> acc ^ p) "" parts

(* buffer_build appends into one pre-sized Buffer: a handful of allocations. *)
let buffer_build parts =
  let total = List.fold_left (fun n p -> n + String.length p) 0 parts in
  let b = Buffer.create total in
  List.iter (Buffer.add_string b) parts;
  Buffer.contents b

(* words_allocated returns the number of minor-heap words f allocates. *)
let words_allocated f =
  let before = Gc.minor_words () in
  ignore (Sys.opaque_identity (f ()));
  Gc.minor_words () -. before
