import XCTest

@testable import Shapes

final class ShapesTests: XCTestCase {
  func testDispatchPathsAgree() {
    let cs = [Circle(r: 1.0), Circle(r: 2.0)]
    let dyn: [any Shape] = cs
    XCTAssertEqual(sumStatic(cs), sumDynamic(dyn), accuracy: 1e-9)
  }
}
