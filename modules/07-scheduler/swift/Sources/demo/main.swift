import Foundation
import Sched

let total = 400_000_000
let w = ProcessInfo.processInfo.activeProcessorCount
let t0 = Date()
let acc = await parallelSqrtSum(total, w)
_ = acc
let ms = Int(Date().timeIntervalSince(t0) * 1000)
print("Swift  sqrt-sum \(total / 1_000_000)M / \(w) tasks: \(ms) ms")
