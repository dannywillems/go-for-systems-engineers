import XCTest

@testable import Un

final class UnTests: XCTestCase {
  func testFloatBits() {
    XCTAssertEqual(floatBits(1.0), 0x3F80_0000)
  }

  func testBytesToString() {
    XCTAssertEqual(bytesToString([104, 105]), "hi")
  }
}
