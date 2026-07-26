# Verdict

Cross-cutting synthesis, written last, once the modules exist to support it.
This file is a stub until then; it will summarize, with pointers to the module
that measured each claim:

- where Go's semantics diverge from a Rust/OCaml/Swift/Kotlin engineer's
  intuition, and the concrete cost of each divergence;
- the runtime trade-offs that showed up in the numbers (dispatch, generics
  stenciling, allocation/escape, GC behavior, scheduler latency);
- what each language made easy and what it made painful in the capstone;
- the axes on which a choice between them actually turns, stated as measured
  trade-offs rather than advocacy.

No claim will appear here that is not backed by a module's captured output.
