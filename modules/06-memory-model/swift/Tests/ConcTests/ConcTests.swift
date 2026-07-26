import XCTest

@testable import Conc

final class ConcTests: XCTestCase {
  func testCount() async {
    let n = await count(4, 1000)
    XCTAssertEqual(n, 4000)
  }
}
