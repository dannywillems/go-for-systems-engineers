import XCTest

@testable import Gen

final class GenTests: XCTestCase {
  func testSum() {
    XCTAssertEqual(sum([1, 2, 3, 4, 5]), 15)
    XCTAssertEqual(sum([1.5, 2.5, 3.0]), 7.0)
  }
}
