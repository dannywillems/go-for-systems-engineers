import Testing

@testable import Observability

@Test func buildersAgree() {
  let parts = (0..<64).map { "chunk-\($0);" }
  #expect(concatPlus(parts) == reserve(parts))
  #expect(buildInstrumented(parts) == reserve(parts))
}
