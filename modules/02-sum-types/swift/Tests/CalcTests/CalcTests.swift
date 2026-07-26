import XCTest

@testable import Calc

final class CalcTests: XCTestCase {
  func testEval() {
    let e: Expr = .add(.mul(.lit(2), .lit(3)), .neg(.lit(4)))
    XCTAssertEqual(eval(e), 2)
  }
}
