import Testing

@testable import Account

@Test func throughPublicAPI() throws {
  var a = try Account(100)
  a.deposit(50)
  #expect(a.value() == 150)
  #expect(throws: Overdraft.self) { try a.withdraw(200) }
  #expect(a.value() == 150)
  #expect(throws: Overdraft.self) { try Account(-1) }
}
