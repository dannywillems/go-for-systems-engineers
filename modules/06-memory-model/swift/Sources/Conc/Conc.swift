// Swift 6 enforces data-race safety in the type system via `Sendable` and actor
// isolation: mutable state shared across concurrency domains must be an `actor`
// (or otherwise Sendable). Sharing a plain mutable class across tasks is a
// COMPILE error (see reject-swift). An actor serializes its mutations.

// region:actor:start

public actor Counter {
  private var n = 0
  public init() {}
  public func inc() { n += 1 }
  public func value() -> Int { n }
}

/// Increments an actor from `threads` concurrent child tasks. The actor's
/// isolation makes every `inc()` mutually exclusive without an explicit lock.
public func count(_ threads: Int, _ per: Int) async -> Int {
  let c = Counter()
  await withTaskGroup(of: Void.self) { group in
    for _ in 0..<threads {
      group.addTask {
        for _ in 0..<per { await c.inc() }
      }
    }
  }
  return await c.value()
}

// region:actor:end
