// Module 11 in Swift: the same Normalize subject. Tests use swift-testing
// (@Test, #expect, and parameterized @Test(arguments:)); the production tools
// for property and fuzz are swift-testing's arguments, SwiftCheck, and libFuzzer
// via SwiftPM -- see the README.

public func normalize(_ s: String) -> String {
  s.split(whereSeparator: { $0.isWhitespace })
    .joined(separator: " ")
    .lowercased()
}
