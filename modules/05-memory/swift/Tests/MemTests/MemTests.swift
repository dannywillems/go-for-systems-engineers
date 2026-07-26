import XCTest

@testable import Mem

final class MemTests: XCTestCase {
  func testAllocBench() {
    XCTAssertEqual(allocBench(1000), (0..<1000).reduce(0, +))
  }
}
