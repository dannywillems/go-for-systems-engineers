import XCTest

@testable import Bootstrap

final class BootstrapTests: XCTestCase {
  func testSumMatchesClosedForm() {
    for n in [0, 1, 10, 1_000_000] {
      XCTAssertEqual(sum(n), n * (n + 1) / 2)
    }
  }

  func testWordSizeIsEight() {
    XCTAssertEqual(wordSizeBytes(), 8)
  }
}
