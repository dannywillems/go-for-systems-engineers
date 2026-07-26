(* The capstone cache in OCaml 5: a bounded concurrent cache over a slow backend,
   using domains for parallelism, a Mutex-guarded Hashtbl, and
   Semaphore.Counting for backpressure. *)

type t = {
  tbl : (int, int) Hashtbl.t;
  mutex : Mutex.t;
  capacity : int;
  backend_calls : int Atomic.t;
  latency : float; (* seconds *)
  sem : Semaphore.Counting.t;
}

let create capacity max_inflight latency =
  {
    tbl = Hashtbl.create capacity;
    mutex = Mutex.create ();
    capacity;
    backend_calls = Atomic.make 0;
    latency;
    sem = Semaphore.Counting.make max_inflight;
  }

let get t key =
  Mutex.lock t.mutex;
  match Hashtbl.find_opt t.tbl key with
  | Some v ->
      Mutex.unlock t.mutex;
      v
  | None ->
      Mutex.unlock t.mutex;
      Semaphore.Counting.acquire t.sem;
      Atomic.incr t.backend_calls;
      if t.latency > 0. then Unix.sleepf t.latency;
      let v = key * key in
      Semaphore.Counting.release t.sem;
      Mutex.lock t.mutex;
      if Hashtbl.length t.tbl >= t.capacity then (
        match Hashtbl.to_seq_keys t.tbl () with
        | Seq.Cons (k, _) -> Hashtbl.remove t.tbl k
        | Seq.Nil -> ());
      Hashtbl.replace t.tbl key v;
      Mutex.unlock t.mutex;
      v

let backend_calls t = Atomic.get t.backend_calls
let length t = Hashtbl.length t.tbl
