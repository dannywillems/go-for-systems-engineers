import Testing

@testable import Codegen

@Test func synthesizedEquatable() {
  let a = Person(name: "Ada", age: 36)
  #expect(a == Person(name: "Ada", age: 36))
  #expect(a != Person(name: "Ada", age: 37))
}
