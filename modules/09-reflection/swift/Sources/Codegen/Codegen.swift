// Module 09 in Swift: the compiler SYNTHESIZES protocol conformances. Declaring
// `: Equatable` (or Hashable, or Codable) makes the compiler generate the
// implementation member-by-member at compile time -- no runtime reflection. Like
// Rust's derive, it is checked at compile time: synthesis fails to compile if a
// stored property does not itself conform (see reject-swift).

public struct Person: Equatable, Hashable {
  public let name: String
  public let age: Int

  public init(name: String, age: Int) {
    self.name = name
    self.age = age
  }

  public func describe() -> String { "\(name) is \(age)" }
}
