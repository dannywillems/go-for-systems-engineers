// Wall-clock timing for measured.txt. Swift's std library has no in-process
// allocation counter, so the number reported is time, not allocations; use
// Instruments (Allocations / os_signpost) for allocation profiling.

import Foundation
import Observability

let parts = (0..<64).map { String(format: "chunk-%02d;", $0) }

func timeNs(_ iters: Int, _ f: () -> Void) -> Double {
  let clock = ContinuousClock()
  let elapsed = clock.measure {
    for _ in 0..<iters { f() }
  }
  let c = elapsed.components
  let seconds = Double(c.seconds) + Double(c.attoseconds) / 1e18
  return seconds * 1e9 / Double(iters)
}

let iters = 1_000_000
let plusNs = timeNs(iters) { _ = concatPlus(parts) }
let reserveNs = timeNs(iters) { _ = reserve(parts) }
print("Swift concatPlus: \(Int(plusNs)) ns/op  reserve: \(Int(reserveNs)) ns/op")
