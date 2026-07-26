import XCTest

@testable import Ts

final class TsTests: XCTestCase {
  func testTransitionsCompile() {
    let _door: Door<Closed> = Door<Closed>().open().close()
    XCTAssertTrue(true)
  }
}
