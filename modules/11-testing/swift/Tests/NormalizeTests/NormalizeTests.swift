import Foundation
import Testing

@testable import Normalize

// swift-testing's PARAMETERIZED test: one @Test runs once per argument tuple,
// the idiomatic Swift table-driven form.
@Test(arguments: [
  ("  hi  ", "hi"),
  ("a\t\n  b", "a b"),
  ("MiXeD", "mixed"),
  ("   ", ""),
])
func table(input: String, want: String) {
  #expect(normalize(input) == want)
}

// Hand-rolled property loop: idempotence over a generated corpus (dependency-free).
@Test func idempotentProperty() {
  let alphabet: [Character] = ["a", "B", " ", "\t", "\n", "c"]
  var seed: UInt64 = 0x9e37_79b9_7f4a_7c15
  for _ in 0..<10_000 {
    seed = seed &* 6_364_136_223_846_793_005 &+ 1
    let len = Int(seed >> 60)
    var s = ""
    var x = seed
    for _ in 0..<len {
      x = x &* 6_364_136_223_846_793_005 &+ 1
      s.append(alphabet[Int(x >> 61) % alphabet.count])
    }
    let once = normalize(s)
    #expect(normalize(once) == once)
    #expect(once == once.trimmingCharacters(in: .whitespaces))
  }
}
