// The capstone cache in Swift: an ACTOR provides the mutual exclusion (no
// explicit mutex). An `await` at the backend fetch suspends the actor, letting
// other gets proceed, so fetches overlap while map access stays serialized. The
// actor is the interesting variable: it can become a serialization point for a
// hot shared cache -- the numbers show whether it does.

import Foundation

public actor Cache {
  private var entries: [Int: Int] = [:]
  private let capacity: Int
  private let latencyNanos: UInt64
  public private(set) var backendCalls: Int = 0

  public init(capacity: Int, latencyMicros: UInt64) {
    self.capacity = capacity
    self.latencyNanos = latencyMicros * 1000
  }

  public func get(_ key: Int) async -> Int {
    if let v = entries[key] { return v }
    backendCalls += 1
    try? await Task.sleep(nanoseconds: latencyNanos)  // suspends the actor
    let v = key * key
    if entries.count >= capacity, let evict = entries.keys.first {
      entries.removeValue(forKey: evict)
    }
    entries[key] = v
    return v
  }

  public func count() -> Int { entries.count }
}
