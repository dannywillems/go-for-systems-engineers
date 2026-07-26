import XCTest

@testable import Errs

final class ErrsTests: XCTestCase {
  func testChain() throws {
    XCTAssertEqual(try chain(3), 7)
    XCTAssertThrowsError(try chain(60))
    XCTAssertEqual(chainResult(3), .success(7))
  }
}
