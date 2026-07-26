import Cache
import Foundation

let capacity = 256
let keys: UInt64 = 256
let workers = 64
let perWorker = 10_000

func micros(_ d: Duration) -> Double {
  Double(d.components.seconds) * 1e6 + Double(d.components.attoseconds) / 1e12
}

let cache = Cache(capacity: capacity, latencyMicros: 100)
let clock = ContinuousClock()
let t0 = clock.now

let latencies = await withTaskGroup(of: [Double].self) { group -> [Double] in
  for w in 0..<workers {
    group.addTask {
      var lat = [Double]()
      lat.reserveCapacity(perWorker)
      var seed = UInt64(w) &* 2_654_435_761 | 1
      for _ in 0..<perWorker {
        seed = seed &* 6_364_136_223_846_793_005 &+ 1_442_695_040_888_963_407
        let key = Int((seed >> 33) % keys)
        let s = clock.now
        _ = await cache.get(key)
        lat.append(micros(s.duration(to: clock.now)))
      }
      return lat
    }
  }
  var all = [Double]()
  for await l in group { all.append(contentsOf: l) }
  return all
}

let elapsedS = micros(t0.duration(to: clock.now)) / 1e6
let all = latencies.sorted()
let n = all.count
func pc(_ p: Double) -> Double { all[min(Int(p / 100 * Double(n)), n - 1)] }
let backend = await cache.backendCalls
print(
  String(
    format:
      "Swift  %dk gets/%dw: %.0f kops/s  p50=%.1fus p99=%.1fus p999=%.1fus  backend=%.1f%% of gets",
    n / 1000, workers, Double(n) / elapsedS / 1000,
    pc(50), pc(99), pc(99.9),
    100 * Double(backend) / Double(n)))
