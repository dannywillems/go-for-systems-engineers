// Swift concurrency runs `async` tasks on a cooperative thread pool sized to the
// core count (M:N, like Go and tokio). A TaskGroup fans a CPU sweep across child
// tasks; the runtime schedules them onto its worker threads.

import Foundation

public func chunkSum(_ lo: Int, _ hi: Int) -> Double {
  var s = 0.0
  var i = lo
  while i < hi {
    s += Double(i).squareRoot()
    i += 1
  }
  return s
}

public func parallelSqrtSum(_ total: Int, _ workers: Int) async -> Double {
  let chunk = total / workers
  return await withTaskGroup(of: Double.self) { group in
    for k in 0..<workers {
      group.addTask { chunkSum(k * chunk, (k + 1) * chunk) }
    }
    var acc = 0.0
    for await s in group { acc += s }
    return acc
  }
}
