// Module 08 in Swift: the richest access-control ladder of the five --
// `private` (enclosing declaration), `fileprivate` (the file), `internal`
// (the module, the default), `public` (other modules), and `open` (other
// modules may subclass/override). Account is a public struct whose stored
// property is `private`, so consumers use the methods but cannot touch or
// construct balance directly.

public struct Overdraft: Error {}

public struct Account {
  private var balance: Int64

  public init(_ initial: Int64) throws {
    if initial < 0 { throw Overdraft() }
    balance = initial
  }

  public mutating func deposit(_ amount: Int64) { balance += amount }

  public mutating func withdraw(_ amount: Int64) throws {
    if amount > balance { throw Overdraft() }
    balance -= amount
  }

  public func value() -> Int64 { balance }
}
