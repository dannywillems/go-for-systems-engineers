import XCTest

@testable import Cache

final class CacheTests: XCTestCase {
  func testCorrectAndBounded() async {
    let c = Cache(capacity: 16, latencyMicros: 0)
    await withTaskGroup(of: Void.self) { group in
      for w in 0..<8 {
        group.addTask {
          for i in 0..<500 {
            let key = (w * 500 + i) % 100
            let v = await c.get(key)
            XCTAssertEqual(v, key * key)
          }
        }
      }
    }
    let count = await c.count()
    XCTAssertLessThanOrEqual(count, 16)
    let backend = await c.backendCalls
    XCTAssertLessThan(backend, 8 * 500)
  }
}
