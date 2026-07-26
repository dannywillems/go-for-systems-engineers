// Module 10 in Swift: Swift has no in-process allocation counter in its standard
// library (unlike Go's AllocsPerRun, Rust's global allocator, OCaml's
// Gc.minor_words, or the JVM's ThreadMXBean). The idiomatic Swift path is
// os_signpost + Instruments: you bracket a region so the profiler shows its
// duration and allocations on a timeline. So the falsifiable artifact here is a
// wall-clock timing (see the bench target); allocation profiling is external.

import os

let obsLog = OSLog(subsystem: "module10.observability", category: .pointsOfInterest)

// concatPlus builds with += . Swift's String is copy-on-write over a growable
// buffer, so += grows amortized (like Rust, unlike Go's immutable strings).
public func concatPlus(_ parts: [String]) -> String {
  var s = ""
  for p in parts { s += p }
  return s
}

// reserve pre-sizes the string once with reserveCapacity.
public func reserve(_ parts: [String]) -> String {
  let total = parts.reduce(0) { $0 + $1.utf8.count }
  var s = ""
  s.reserveCapacity(total)
  for p in parts { s += p }
  return s
}

// region:signpost:start

// buildInstrumented brackets the build with an os_signpost interval, the
// standard Swift way to make a region measurable in Instruments without a
// bespoke profiler.
public func buildInstrumented(_ parts: [String]) -> String {
  let id = OSSignpostID(log: obsLog)
  os_signpost(.begin, log: obsLog, name: "build", signpostID: id)
  defer { os_signpost(.end, log: obsLog, name: "build", signpostID: id) }
  return reserve(parts)
}

// region:signpost:end
