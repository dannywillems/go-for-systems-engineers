// Swift expresses typestate with phantom generics + CONSTRAINED extensions: a
// method exists only for a specific State, so an invalid transition is a compile
// error (the method is absent). Same idea as Rust's typestate. See reject-swift.

public enum Open {}
public enum Closed {}

public struct Door<State> {
  public init() {}
}

extension Door where State == Closed {
  public func open() -> Door<Open> { Door<Open>() }
}

extension Door where State == Open {
  public func close() -> Door<Closed> { Door<Closed>() }
}
