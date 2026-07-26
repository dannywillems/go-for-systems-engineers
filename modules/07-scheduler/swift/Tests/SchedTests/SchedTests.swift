import XCTest

@testable import Sched

final class SchedTests: XCTestCase {
  func testChunkSum() {
    let expected = 0.0 + 1.0 + 2.0.squareRoot() + 3.0.squareRoot()
    XCTAssertEqual(chunkSum(0, 4), expected, accuracy: 1e-9)
  }

  func testParallelMatchesSerial() async {
    let par = await parallelSqrtSum(400, 4)
    XCTAssertEqual(par, chunkSum(0, 400), accuracy: 1e-6)
  }
}
