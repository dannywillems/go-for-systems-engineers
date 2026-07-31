# 07b — The theory of concurrency, parallelism, and scheduling

The [README](README.md) measures Go's scheduler and the [COMPARISON](COMPARISON.md)
tabulates the five models. This document explains the *theory* underneath: the
questions every language must answer to run many computations, the formal models
those answers come from, and where each of the five sits. Read it top-to-bottom
once; then the per-language sections stand alone.

Throughout, "measured" means a claim the repo verifies with a running program
(linked to Module 06/07); everything else is design/theory with a citation.

## Contents

- [First distinction: concurrency is not parallelism](#first-distinction-concurrency-is-not-parallelism)
- [The seven questions a concurrency model answers](#the-seven-questions-a-concurrency-model-answers)
  - [1. The execution model: 1:1, N:1, M:N](#1-the-execution-model-11-n1-mn)
  - [2. Stackful vs stackless: the coloring split](#2-stackful-vs-stackless-the-coloring-split)
  - [3. Cooperative vs preemptive scheduling](#3-cooperative-vs-preemptive-scheduling)
  - [4. Work-stealing: the scheduler's core algorithm](#4-work-stealing-the-schedulers-core-algorithm)
  - [5. The coordination paradigm: shared memory, CSP, actors, effects](#5-the-coordination-paradigm-shared-memory-csp-actors-effects)
  - [6. The memory model: what a data race means](#6-the-memory-model-what-a-data-race-means)
  - [7. Structured concurrency: the lifetime discipline](#7-structured-concurrency-the-lifetime-discipline)
- [The five languages, placed](#the-five-languages-placed)
  - [Go: CSP over a preemptive M:N runtime](#go-csp-over-a-preemptive-mn-runtime)
  - [Rust: 1:1 threads, or stackless futures with compile-time race freedom](#rust-11-threads-or-stackless-futures-with-compile-time-race-freedom)
  - [OCaml 5: domains for parallelism, effect handlers for concurrency](#ocaml-5-domains-for-parallelism-effect-handlers-for-concurrency)
  - [Swift: the actor model with compile-time isolation](#swift-the-actor-model-with-compile-time-isolation)
  - [Kotlin/JVM: stackless coroutines, and stackful virtual threads](#kotlinjvm-stackless-coroutines-and-stackful-virtual-threads)
  - [Aside: linear logic — why Rust and OCaml diverge on resources](#aside-linear-logic--why-rust-and-ocaml-diverge-on-resources)
- [The synthesis matrix](#the-synthesis-matrix)
- [The three deepest contrasts](#the-three-deepest-contrasts)
- [References](#references)

## First distinction: concurrency is not parallelism

Rob Pike's formulation is the one to keep: **concurrency is about *structure*,
parallelism is about *execution*.** Concurrency is composing a program out of
independently-executing computations (dealing with many things at once);
parallelism is running computations literally simultaneously on multiple
processors (doing many things at once). A concurrent program can run on one core
(the computations interleave) or many; parallelism is what happens when the
scheduler places concurrent work on multiple cores.

The distinction matters because each language draws the line differently, and
much of the confusion about "goroutines vs coroutines vs threads" dissolves once
you separate the two:

- **OCaml 5 separates them most cleanly**: *domains* are the unit of parallelism
  (one per core), *fibers/effects* are the unit of concurrency (many per domain).
- **Go fuses them behind one keyword**: `go f()` creates concurrency
  (a goroutine); `GOMAXPROCS` logical processors provide the parallelism.
- **Rust splits them by mechanism**: OS threads give parallelism *and*
  concurrency; `async` tasks give concurrency, and the executor's worker threads
  add parallelism.
- **Swift and Kotlin** put concurrency in a language feature (tasks/coroutines)
  and parallelism in a thread pool (the cooperative pool / dispatchers).

Everything below is really about **concurrency** — how a language represents,
suspends, resumes, and coordinates many logical computations — plus how it maps
them onto a smaller number of CPUs for **parallelism**.

## The seven questions a concurrency model answers

### 1. The execution model: 1:1, N:1, M:N

A "logical thread" (goroutine, task, coroutine, fiber) must eventually run on a
kernel/OS thread, the only thing an OS actually schedules onto a core. How
logical threads map to kernel threads is the first design axis:

- **1:1** — each logical thread *is* one kernel thread. Simple, and the OS does
  the scheduling (preemptively, via timer interrupts), but kernel threads are
  heavy (~1 MB stack, microsecond-scale creation, a syscall to switch), so you
  get thousands, not millions. _Rust `std::thread`, JVM platform threads, OCaml
  domains, C `pthread`._
- **N:1** — many logical threads on *one* kernel thread. Cheap logical threads,
  but no parallelism and one blocking call stalls everyone. _Classic
  single-threaded event loops._
- **M:N** — many logical threads multiplexed over K kernel threads by a
  user-space or runtime scheduler. Cheap logical threads (thousands to millions)
  *and* parallelism, at the cost of a scheduler you (or the runtime) must write,
  which in particular must handle **blocking**: when a logical thread blocks, the
  scheduler has to park it and free the kernel thread for other work.
  _Go goroutines, JVM virtual threads (Loom), Kotlin coroutines over a dispatcher
  pool, Rust `async` over tokio, Swift tasks over the cooperative pool, OCaml
  fibers over domains (Eio)._

The whole art of an M:N runtime is the blocking-handling: Go's netpoller parks a
goroutine blocked on I/O and reuses its thread; Loom *unmounts* a blocked virtual
thread's continuation from its carrier; an `async` task simply returns `Pending`
at the blocking point. All three solve the same problem — don't waste a kernel
thread on a logical thread that is waiting.

### 2. Stackful vs stackless: the coloring split

This is the deepest and least-understood axis, and it explains almost everything
else. **How is a *suspended* computation represented?**

- **Stackful** (a *fiber* or *green thread*): the logical thread owns its own
  heap-allocated, growable **call stack**. Suspending saves the whole stack (swap
  the stack pointer); resuming restores it. Because the entire stack is
  preserved, the computation can suspend **anywhere** — arbitrarily deep in a
  call tree — and the code looks like ordinary blocking code. There is **no
  function coloring**. The cost is a stack per logical thread (kilobytes, though
  growable). _Go goroutines (~2 KB growable stacks — measured in the README),
  OCaml effect fibers, JVM virtual threads (stack chunks on the heap)._
- **Stackless** (a *state machine*): the compiler transforms an `async`/`suspend`
  function into a state machine — a struct/enum holding the live locals plus a
  "resume label." Suspending means *returning that state machine to the caller*;
  it has no stack of its own. Consequently it can only suspend at **explicit,
  statically-visible points** (`await`/`suspend`), and those points must be
  visible *up the entire call chain* — which is exactly Bob Nystrom's **"what
  color is your function?"**: `async` ("red") functions can only be called from
  other `async` functions, so async infects the call graph. The upside is that a
  suspended task is tiny (only the captured locals) and needs no stack.
  _Rust `Future`s, Swift tasks, Kotlin coroutines, C++20 coroutines, JS `async`._

The unifying concept is the **continuation** — "the rest of the computation." A
stackful fiber reifies its continuation *as a stack*; a stackless coroutine
reifies it *as a state machine object* (Kotlin literally CPS-transforms the
function, threading an extra `Continuation` parameter). OCaml's **delimited
continuations** (below) are the general form: with them you can build stackful,
uncolored fibers out of a first-class control operator.

The consequences to remember:

| | Stackful (fiber) | Stackless (state machine) |
| --- | --- | --- |
| Suspend anywhere? | yes | only at `await`/`suspend` |
| Function coloring? | **no** | **yes** |
| Blocking-style code? | yes (direct style) | no (must thread `await`) |
| Memory per suspended unit | a (growable) stack | just the captured locals |
| Who pays | RAM for stacks | ergonomics (coloring) |

### 3. Cooperative vs preemptive scheduling

Once several logical threads share a kernel thread, *who decides when one yields?*

- **Cooperative**: a logical thread runs until it *voluntarily* yields — at an
  `await`, a `suspend`, a channel operation, or an effect. A CPU-bound loop with
  no yield point **starves** its worker thread and can wedge the runtime. This is
  why every async runtime ships a blocking-offload escape hatch
  (`spawn_blocking`, `Dispatchers.IO`, a detached task): _Rust tokio, Swift,
  Kotlin coroutines are all cooperative._ There is **no forward-progress
  guarantee** for a task that never suspends.
- **Preemptive**: the scheduler can forcibly suspend a running logical thread.
  OS threads (1:1) are preempted by the kernel on a timer. The interesting case
  is a **preemptive M:N runtime**: Go added **asynchronous preemption** in 1.14 —
  a monitor thread signals a long-running goroutine (`SIGURG`) to yield at an
  async-safe point, even inside a tight loop with no function calls (measured in
  the README: `main` progresses past a spinning goroutine). This eliminates an
  entire class of bug — a buggy tight loop cannot starve the program — at the
  cost of runtime machinery (safe points, signal handling, precise stack maps).

Go is the outlier: it gives cheap M:N green threads *and* full preemption. The
async runtimes trade that away for a simpler, stackless design; the JVM's virtual
threads are cooperative for CPU work (a pinned virtual thread does not
time-slice) but their carrier threads are OS-preemptible.

### 4. Work-stealing: the scheduler's core algorithm

An M:N runtime with several worker threads needs to balance work without a
central bottleneck. The standard answer, with provable bounds, is **randomized
work stealing** (Blumofe & Leiserson, *Scheduling Multithreaded Computations by
Work Stealing*, JACM 1999, from the Cilk project): each worker owns a
double-ended queue (deque) of ready tasks; it pushes/pops its *own* end (LIFO,
cache-friendly), and when it runs dry it **steals** from the *other* end of a
*randomly chosen* victim's deque. The theorem: a computation with work `T₁` and
span (critical path) `T∞` runs in expected time `T₁/P + O(T∞)` on `P` processors,
within a constant factor of sequential space — near-optimal, with no global lock.

This same algorithm is under **Go's** P-local run queues (an idle P steals half
of a random P's queue), **tokio's** multi-threaded executor, the JVM
**ForkJoinPool** that schedules **virtual threads** and Kotlin's `Default`
dispatcher, and OCaml **Domainslib**. When people say a runtime "keeps all cores
busy," this is the mechanism.

### 5. The coordination paradigm: shared memory, CSP, actors, effects

How do concurrent computations *coordinate* and share data? Four lineages, from
1970s theory:

- **Shared memory + locks.** Logical threads share an address space and
  synchronize with mutexes/atomics; semantics are the *interleavings* of their
  memory accesses. Universal and fast, but the source of data races. _All five
  have this substrate; Rust and OCaml make it safer (below)._
- **CSP** — *Communicating Sequential Processes* (Hoare, 1978). *Anonymous*
  processes communicate over synchronous **channels**; a send blocks until a
  receive (a rendezvous); `select` offers external choice among channel
  operations. **Go's channels and `select` are CSP** (via Pike's Newsqueak → Alef
  → Limbo lineage). "Do not communicate by sharing memory; share memory by
  communicating."
- **Actors** (Hewitt, 1973; Agha, 1986). *Named* actors own private state and a
  **mailbox**; you send a message *asynchronously* to an actor's identity, and an
  actor processes one message at a time, so its state is never concurrently
  accessed. **Swift's `actor` is the actor model.** The CSP-vs-actor difference:
  CSP couples processes to a *channel* and is synchronous and anonymous; actors
  couple senders to an *identity* and are asynchronous with buffered mailboxes
  (which is why actors scale naturally to distributed systems — Erlang).
- **Algebraic effects & delimited continuations.** Rather than fix a paradigm,
  provide a general control operator and *build* the paradigm as a library. This
  is **OCaml 5's** route: a scheduler is an effect *handler* that interprets
  `Fork`/`Yield`/`Await` effects by manipulating first-class continuations. Effect
  handlers subsume exceptions, generators, coroutines, and async/await, so CSP,
  actors, or futures are all expressible on top.

(Two more worth naming: **futures/promises/dataflow** — a value that will exist
later, composed by combinators; Rust `Future`, Kotlin `Deferred`, JS `Promise` —
and **software transactional memory**, optimistic shared-memory transactions, as
in Haskell/Clojure; none of the five use STM as their primary model.)

### 6. The memory model: what a data race means

When two threads touch the same memory without synchronization and at least one
writes, that is a **data race**. A language's *memory model* says what such a
program is allowed to do, and it is the difference between "subtle bug" and
"undefined behavior." The key theorem is **DRF-SC** (Adve & Hill): a program with
*no* data races behaves *sequentially consistently* (as some simple interleaving)
— so if you synchronize properly, you may reason with the intuitive model. The
models differ entirely in what happens when there *is* a race. There are three
regimes, and the five languages occupy all three:

- **Compile-time prevention (races cannot exist).** The type system rejects racy
  programs. _Rust (ownership + `Send`/`Sync`) and Swift 6 (actor isolation +
  `Sendable`) — Section below and Module 06._
- **Bounded, defined races.** OCaml 5's memory model (Dolan et al., *Bounding
  Data Races in Space and Time*, PLDI 2018) gives **local DRF-SC**: race-free
  regions stay sequentially consistent, and a race's effects are *bounded* — a
  racy read returns *some* valid prior value, never garbage, with no
  "out-of-thin-air" results and no undefined behavior. This is strictly stronger
  than C++/Go.
- **Undefined or memory-unsafe races.** _Go_: the memory model is happens-before
  based (a channel send happens-before the corresponding receive, an unlock
  before a lock), and a race on a *multiword* value (an interface or slice header)
  is **undefined** — it can tear (Module 06). The `-race` detector finds races
  *dynamically*, only on executed paths. _Kotlin/JVM_: the Java Memory Model is
  also happens-before, but the JVM stays *memory-safe* under a race (references
  are read/written atomically, no tearing, no UB) — you get stale or reordered
  values, not corruption. Neither prevents races statically.

### 7. Structured concurrency: the lifetime discipline

The last axis is not about the CPU but about *program structure*. In the
**unstructured** model, `go f()` / `thread::spawn` / `Task { }` starts a
computation whose lifetime is independent of its spawner: it can outlive the
function that launched it, and errors and cancellation do not propagate. Nathaniel
J. Smith's *"Notes on structured concurrency, or: Go statement considered
harmful"* (2018) argues this is the concurrency analogue of `goto`: it breaks the
black-box abstraction of a function, because calling a function might silently
leak background tasks.

**Structured concurrency** imposes the discipline that control flow which *splits*
must *rejoin*: concurrent children are bound to a lexical scope (a "nursery"), the
scope does not return until all children finish (or are cancelled), and errors and
cancellation propagate through the resulting task *tree*. This makes concurrency
compose like ordinary function calls. _Swift task groups / `async let`, Kotlin
`coroutineScope` + the `Job` hierarchy, OCaml `Eio.Switch`, Python Trio, and
Java's `StructuredTaskScope` all implement it; Go's `go` and Rust's `spawn` are
the unstructured primitives it critiques (`errgroup`, `tokio::JoinSet`, and
`std::thread::scope` recover parts of it)._

## The five languages, placed

### Go: CSP over a preemptive M:N runtime

Go's answer is a **runtime-provided M:N scheduler of stackful goroutines, with
CSP for coordination and asynchronous preemption**.

- **Execution (GMP).** A goroutine (**G**) is a stackful green thread with a
  ~2 KB growable stack. An **M** is an OS thread; a **P** is a logical processor
  (there are `GOMAXPROCS` of them) that holds a run queue. An M must hold a P to
  run Gs. Per-P local queues plus a global queue, balanced by **work stealing**
  (an idle P steals half of a random P's queue). This is pure M:N.
- **Blocking.** The **netpoller** (epoll/kqueue) parks a goroutine blocked on
  network I/O and *frees its M* to run other Gs; when the I/O is ready the G is
  re-queued — so blocking-looking code is non-blocking underneath. A blocking
  *syscall* blocks the M, and the P is handed to another M so other Gs proceed.
- **Preemption.** Cooperative yield points (channel ops, allocation, calls) plus
  **asynchronous preemption** (1.14): a tight loop cannot starve a P (README).
- **Paradigm.** CSP channels (unbuffered = rendezvous, buffered = bounded async)
  and `select` (external choice), *and* shared memory via `sync`.
- **Memory model.** Happens-before; multiword races are undefined; `-race` is a
  dynamic detector (Module 06).
- **Structure.** Unstructured `go`; the scheduler is *mandatory* (you cannot opt
  out of the runtime).

The trade: the smallest surface (one keyword, cheap uncolored goroutines,
preemption) at the cost of no compile-time race safety and no structured
concurrency in the language.

### Rust: 1:1 threads, or stackless futures with compile-time race freedom

Rust deliberately has **no runtime**, so it offers two stories and makes you
choose.

- **OS threads (`std::thread`).** 1:1 with kernel threads, preemptively scheduled
  by the OS. Real parallelism, zero runtime overhead, direct C interop. (Rust
  *had* M:N green threads before 1.0 and **removed** them — RFC 230 — precisely to
  avoid a mandatory runtime.)
- **`async`/`await`.** Stackless coroutines: an `async fn` compiles to a state
  machine implementing the **`Future`** trait. The model is **poll-based (pull)**:
  an executor calls `Future::poll`, which returns `Ready` or `Pending`; on
  `Pending` the future registers a **`Waker`**, and when progress is possible the
  waker re-schedules it. Futures are **lazy** (they do nothing until polled —
  unlike eager JS promises). Crucially, **the executor is a library**, not part of
  the language: tokio/async-std/smol provide the (work-stealing, M:N) runtime, and
  you pick one. Cooperative; `spawn_blocking` offloads CPU/blocking work. `async`
  is **colored**.
- **The crown jewel: compile-time data-race freedom.** Ownership (an *affine*
  type discipline: a value has one owner, and moving transfers it) plus borrowing
  (shared `&T` *xor* unique `&mut T` — never aliasing with mutation) plus the
  `Send`/`Sync` auto-traits (safe to move to / share across threads) let the
  compiler **mechanically prove the absence of data races**, at zero runtime cost.
  A `&mut` cannot be aliased, so two threads cannot mutate the same datum; a
  `!Send` type (`Rc`) cannot cross a thread boundary. This is the strongest static
  guarantee of the five, and it is *region/ownership typing*, not a runtime check
  (Module 06 shows the compile-rejects).

The trade: maximum control and the only zero-cost compile-time race freedom, at
the cost of the coloring problem, an ecosystem-level runtime choice, and the
borrow checker's learning curve.

### OCaml 5: domains for parallelism, effect handlers for concurrency

OCaml 5 is the most theoretically distinctive: it factors the problem exactly
along the concurrency/parallelism line and builds concurrency from a *general
control operator* rather than a fixed scheduler.

- **Parallelism = domains.** A **domain** maps 1:1 to an OS thread and is the unit
  of parallelism; you spawn roughly one per core (they are heavyweight).
- **Concurrency = effect handlers.** OCaml 5 retrofitted **algebraic effect
  handlers** (Sivaramakrishnan et al., PLDI 2021). A **fiber** is a small,
  heap-allocated, growable stack chunk implementing a **one-shot delimited
  continuation**. `perform E` raises an effect; the enclosing `match … with
  effect E, k -> …` handler receives the delimited continuation `k` (the rest of
  the fiber up to the handler) and decides whether/when to `continue k`. **A
  scheduler is just an effect handler** that interprets `Fork`/`Yield`/`Await` by
  juggling continuations — Eio, Domainslib, and Lwt-on-effects are all user-level
  schedulers written this way, in a few lines.
- **Uncolored, direct style.** Because fibers are *stackful* (the continuation is
  a real stack chunk), suspension is transparent: Eio code reads like ordinary
  blocking code with **no `async`/`await` keyword and no coloring** — the same
  ergonomic win as goroutines, obtained through a first-class control mechanism
  rather than a baked-in runtime. (The continuations are *one-shot*, which is
  enough for schedulers and cheaper than multi-shot; effects are currently
  *untyped* — absent from signatures — a known gap.)
- **Memory model.** **Local DRF-SC** with bounded races (PLDI 2018) — the
  strongest of the five short of static prevention (Section 6).

The trade: the most general and composable concurrency (write your own scheduler)
and an unusually strong memory model, at the cost of youth (the 5.x ecosystem)
and effects not yet appearing in types.

### Swift: the actor model with compile-time isolation

Swift's model is **the actor model plus structured concurrency, made data-race
safe by the compiler**.

- **Execution.** Tasks (stackless async units) run on a **cooperative thread
  pool** bounded to the core count — deliberately avoiding "thread explosion."
  `async`/`await` suspends by capturing a **continuation** (a partial task), not a
  full stack; cooperative; `async` is **colored**.
- **Actors.** An `actor` type has **isolated** mutable state that only one task
  may touch at a time — the actor model. A call into an actor from outside is
  `await` (an asynchronous message send); actors are *reentrant* (they may
  interleave at `await` points, unlike a plain mutex). `@MainActor` is the
  main-thread actor.
- **Structured concurrency (SE-0304).** Child tasks (`async let`, task groups) are
  scoped to the parent, which awaits them all; **cancellation propagates down** the
  task tree and errors propagate up. Unstructured `Task { }` exists but is the
  exception.
- **Compile-time data-race freedom, via *isolation* rather than *ownership*.**
  `Sendable` (SE-0302) bounds what may cross an isolation boundary (value types are
  `Sendable` by copy); actor isolation confines mutable state; and **region-based
  isolation** (SE-0414) tracks an object's *provenance* so a non-`Sendable` value
  can be *transferred* across a boundary when the compiler proves no other
  reference remains (a flow-sensitive, affine-flavored analysis). So Swift reaches
  the same "no data races at compile time" as Rust by a *different theory* —
  isolate state in actors and check the crossings, versus own and borrow every
  value.

The trade: a high-level, safe-by-construction actor model with first-class
structured concurrency, at the cost of coloring and a concurrency checker whose
rules (Sendable, isolation) are their own learning curve.

### Kotlin/JVM: stackless coroutines, and stackful virtual threads

Kotlin is the instructive case because the JVM offers **both** sides of the
stackful/stackless split, and Kotlin can use either.

- **Coroutines (stackless, colored, library).** A `suspend` function is
  **CPS-transformed** by the compiler: it gains a hidden `Continuation` parameter
  and its body becomes a finite state machine (locals → FSM fields, each
  suspension point → a labeled state; suspending returns `COROUTINE_SUSPENDED` and
  stores state on the heap; resuming calls `continuation.resumeWith` at the saved
  label). `suspend` is the **color**. **Dispatchers** (`Default` sized to cores for
  CPU work, `IO` larger for blocking, `Main`) are thread pools; they *intercept*
  continuations to switch threads. Cooperative; cancellation is cooperative
  (checked at suspension points).
- **Structured concurrency, first-class.** `CoroutineScope` + the `Job` hierarchy:
  `coroutineScope { }`, `launch`, and `async` create child coroutines whose `Job`s
  form a tree; the scope completes only when its children do, **cancellation
  propagates down**, and a child failure cancels its siblings (unless a
  `SupervisorJob`). This is Kotlin's realization of nurseries.
- **Virtual threads (stackful, uncolored, JVM runtime).** Project Loom (JEP 444,
  JDK 21) adds **virtual threads**: M:N green threads whose state is a
  **continuation** mounted on a **carrier thread** (a work-stealing ForkJoinPool).
  A virtual thread presents an ordinary **blocking** API (`Thread.sleep`, blocking
  I/O) but *unmounts* its continuation from the carrier when it blocks — so
  blocking code becomes cheap and there is **no coloring**. Kotlin can run
  coroutines on a virtual-thread dispatcher, or use virtual threads directly. The
  contrast within one platform is the clearest illustration of the whole axis:
  *coroutines* = stackless + colored + a library; *virtual threads* = stackful +
  uncolored + the runtime.
- **Memory model.** The Java Memory Model — happens-before, memory-safe under
  races (no tearing/UB, but stale/reordered values), no static race prevention.

The trade: mature, ergonomic structured coroutines and (via Loom) cheap uncolored
blocking, at the cost of JVM overheads (erasure, object headers — Module 05) and
no compile-time race safety.

### Aside: linear logic — why Rust and OCaml diverge on resources

The sharpest way to state the Rust/OCaml difference is through **linear logic**
(Girard, 1987) and the Curry–Howard reading of its **substructural** fragments.
Ordinary (intuitionistic) logic admits two *structural rules* on hypotheses:
**weakening** (a hypothesis may go *unused* — you can discard it) and
**contraction** (a hypothesis may be used *more than once* — you can duplicate
it). Dropping them gives a hierarchy, and under Curry–Howard each corresponds to
a resource discipline on values:

- **structural** (both rules): a value may be used *any number of times* — the
  normal, GC-friendly world;
- **affine** (weakening only, no contraction): *at most once* — you may drop it,
  not duplicate it;
- **relevant** (contraction only, no weakening): *at least once*;
- **linear** (neither): *exactly once*.

The `!` exponential ("of course") re-admits the structural rules for a marked
proposition `!A` — i.e. it marks a value as *freely* copyable and droppable.

**Rust makes affine linear logic the default discipline of the entire
language.** Ownership *is* this type system: a value is used **at most once** —
moving it consumes it (no implicit duplication = no contraction), and letting it
go out of scope drops it (weakening is allowed). `Copy`/`Clone` are exactly the
`!` exponential: a type opts *back into* structural, freely-duplicable use.
Borrows split along the same seam — a shared `&T` is duplicable (structural), a
unique `&mut T` is a *linear* capability that cannot be aliased — so the
"shared **xor** mutable" rule is the linear/`!` distinction itself, and lifetimes
add the region layer on top. This is not decoration: the affine/unique discipline
is *precisely* what buys **compile-time data-race freedom** in the concurrency
sections above. A `&mut` is a unique capability, so no two threads can hold a
mutable path to the same cell; `Send`/`Sync` decide which capabilities may cross a
thread boundary. Race freedom is a *theorem of the linear type discipline*, paid
for at compile time.

**OCaml keeps structural (intuitionistic) logic as its default** — values are
freely duplicated and discarded, and *the garbage collector*, not the type
system, manages their lifetimes. That is why OCaml reaches for a **memory model**
(bounded LDRF-SC) rather than a type discipline to tame races: it never took the
linear turn in its core. Linearity appears only at the **edges**, and one of them
is squarely in this document's territory: OCaml's **effect continuations are
one-shot** — `continue k` may be invoked **at most once** (a second call raises
`Continuation_already_resumed`). A one-shot continuation is an **affine
resource**, and a scheduler written as an effect handler must respect that
affinity. (A fully *linear* system would demand resuming *exactly* once.) The
other edge is opt-in: **OxCaml's modes** (`unique`/`aliased`, `once`/`many` —
uniqueness and linearity) are a linear-logic-inspired extension bringing
Rust-style, stack-allocatable, no-GC resource control to OCaml *as a mode
annotation*, without giving up the GC-managed structural default.

So the one-line differentiation: **Rust is affine-by-default (ownership = the
type system, and the source of its static race freedom); OCaml is
structural-by-default (GC-managed), admitting linearity only where it must — the
affine one-shot continuation — or where you opt in (modes).** The rest of the
substructural story, with the two languages placed against the other three, is in
Module 14's [`PROPOSITIONS.md`](../14-type-system/PROPOSITIONS.md) (rows 27–29:
affine, linear, regions).

## The synthesis matrix

| Axis | Go | Rust | OCaml 5 | Swift | Kotlin/JVM |
| ---- | -- | ---- | ------- | ----- | ---------- |
| Logical-thread unit | goroutine | OS thread / `async` task | fiber (effect) / domain | task / actor | coroutine / virtual thread |
| Mapping | M:N | 1:1 (threads); M:N (async, lib) | 1:1 domains; M:N fibers | M:N | M:N |
| Stack model | **stackful** | **stackless** (async) | **stackful** (fibers) | **stackless** | **stackless** (coroutines) / **stackful** (Loom) |
| Function coloring | no | **yes** (async) | no | **yes** | **yes** (coroutines) / no (Loom) |
| Scheduling | **preemptive** (async preemption) | OS-preemptive (threads); cooperative (async) | cooperative (Eio) | cooperative | cooperative (coroutines) / OS-preemptive carriers |
| Paradigm | **CSP** channels | shared memory + ownership | **effects / continuations** | **actors** | coroutines + channels |
| Scheduler lives in | the runtime | a **library** (tokio…) | a **library** (Eio…) | the runtime | a **library** (kotlinx) / the JVM (Loom) |
| Data-race freedom | runtime (`-race`), UB on tear | **compile-time** (ownership) | **bounded** (LDRF-SC) | **compile-time** (isolation) | runtime (JMM), memory-safe |
| Structured concurrency | no (native) | no (native; scoped threads) | yes (`Eio.Switch`) | **yes** | **yes** |
| Work-stealing scheduler | yes | yes (tokio) | yes (Domainslib) | yes | yes (ForkJoinPool) |

(The runtime throughput of these models on a CPU-bound workload lands within
~15% across the five — measured in the [README](README.md) and
[COMPARISON](COMPARISON.md). The theory below is about *structure and safety*, not
raw speed; for this class of work the model is a low-single-digit factor, as
Module 12's capstone also finds.)

## The three deepest contrasts

If you remember only three things:

1. **Stackful vs stackless is the master axis.** It determines coloring (Nystrom),
   whether you can write blocking-style code, and the memory cost of a suspended
   unit. Go/OCaml/Loom are stackful and uncolored; Rust/Swift/Kotlin-coroutines are
   stackless and colored. Everything about "why is async 'infectious'?" reduces to
   this.

2. **Data-race freedom lives in three different places.** Rust and Swift move it
   to **compile time** (by ownership and by isolation, respectively — two distinct
   theories reaching the same guarantee). OCaml moves it into a **bounded memory
   model** (races are defined and contained). Go and the JVM leave it at
   **runtime** (a dynamic detector; memory-safe-but-racy values). This is the axis
   that decides whether "fearless concurrency" is a compiler promise or a testing
   activity.

3. **The scheduler is either the runtime's or yours.** Go and Swift and Loom bake
   a scheduler into the runtime (uniform, mandatory, preemptive-ish). Rust, OCaml,
   and Kotlin coroutines put it in a **library** — with OCaml's effect handlers
   making "write your own scheduler" a first-class, few-line exercise, the most
   general point on the axis. Runtime-provided buys uniformity; library-provided
   buys composability (and fragmentation).

None of these is strictly better; they are different cuts through the same
theory. Go optimizes for the smallest uncolored primitive with preemption; Rust
for zero-cost control and static safety; OCaml for generality via effects and a
strong memory model; Swift for a safe high-level actor model; Kotlin for ergonomic
structured coroutines with the JVM's two runtimes underneath.

## References

Foundations:

- Hoare, "Communicating Sequential Processes" (CACM 1978): https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf
- Hewitt, Bishop, Steiger, "A Universal Modular ACTOR Formalism…" (1973); Agha,
  *Actors* (1986). Overview of CSP vs actors: https://www.karanpratapsingh.com/blog/csp-actor-model-concurrency
- Blumofe & Leiserson, "Scheduling Multithreaded Computations by Work Stealing"
  (JACM 1999): https://www.csd.uwo.ca/~mmorenom/CS433-CS9624/Resources/Scheduling_multithreaded_computations_by_work_stealing.pdf
- Adve & Hill, "Weak Ordering — A New Definition" (1990, the DRF-SC theorem);
  readable overview across languages: Cox, "Programming Language Memory Models":
  https://research.swtch.com/plmm
- Nystrom, "What Color is Your Function?" (2015): https://journal.stuffwithstuff.com/2015/02/01/what-color-is-your-function/
- Smith, "Notes on structured concurrency, or: Go statement considered harmful"
  (2018): https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/
- Pike, "Concurrency is not Parallelism": https://go.dev/blog/waza-talk
- Girard, "Linear Logic" (1987): https://en.wikipedia.org/wiki/Linear_logic ;
  Walker, "Substructural Type Systems" (in *Advanced Topics in Types and
  Programming Languages*), overview: https://en.wikipedia.org/wiki/Substructural_type_system
- OxCaml modes (uniqueness and linearity for OCaml): https://oxcaml.org/

Go:

- The Go scheduler design (Vyukov): https://go.dev/s/go11sched
- The Go Memory Model: https://go.dev/ref/mem
- Go 1.14 asynchronous preemption (release notes): https://go.dev/doc/go1.14#runtime

Rust:

- Boats, "Why async Rust?": https://without.boats/blog/why-async-rust/
- Tokio, "Async in depth" (poll model, wakers): https://tokio.rs/tokio/tutorial/async
- Asynchronous Programming in Rust (the async book): https://rust-lang.github.io/async-book/
- Fearless Concurrency / `Send` and `Sync` (the Rust Book): https://doc.rust-lang.org/book/ch16-04-extensible-concurrency-sync-and-send.html

OCaml 5:

- Sivaramakrishnan et al., "Retrofitting Effect Handlers onto OCaml" (PLDI 2021): https://arxiv.org/abs/2104.00250
- Effect handlers (the manual): https://ocaml.org/manual/5.4/effects.html
- OCaml memory model ("The hard bits", LDRF-SC): https://ocaml.org/manual/5.4/memorymodel.html
- Dolan et al., "Bounding Data Races in Space and Time" (PLDI 2018): https://kcsrk.info/papers/pldi18-memory.pdf
- Eio (effects-based direct-style I/O): https://github.com/ocaml-multicore/eio

Swift:

- SE-0304, Structured Concurrency: https://github.com/apple/swift-evolution/blob/main/proposals/0304-structured-concurrency.md
- SE-0306, Actors: https://github.com/apple/swift-evolution/blob/main/proposals/0306-actors.md
- SE-0302, Sendable: https://github.com/apple/swift-evolution/blob/main/proposals/0302-concurrent-value-and-concurrent-closures.md
- SE-0414, Region-Based Isolation: https://github.com/apple/swift-evolution/blob/main/proposals/0414-region-based-isolation.md

Kotlin / JVM:

- Kotlin coroutines guide: https://kotlinlang.org/docs/coroutines-guide.html
- Elizarov, "How do you color your functions?" (Kotlin and coloring): https://elizarov.medium.com/how-do-you-color-your-functions-a6bb423d936d
- JEP 444, Virtual Threads (Project Loom): https://openjdk.org/jeps/444
- "The Basis of Virtual Threads: Continuations": https://foojay.io/today/the-basis-of-virtual-threads-continuations/
