import Foundation
import Mem

let n = 50_000_000
let t0 = Date()
let acc = allocBench(n)
let ms = Int(Date().timeIntervalSince(t0) * 1000)
print("Swift alloc 50M (class/ARC): \(ms) ms (acc=\(acc))")
