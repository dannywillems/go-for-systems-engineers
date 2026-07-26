// Swift protocols are nominal like Rust traits (an explicit conformance), but
// Swift 6 makes the existential/universal split syntactic:
//   - `any Shape`  is the existential (a boxed value + witness table),
//   - `some Shape` / `<S: Shape>` is universal, resolved statically.
// Unlike Rust, Swift allows RETROACTIVE conformance (conforming a type you do
// not own to a protocol you do not own), so it has no hard orphan rule; it only
// warns, and the app is expected to avoid conflicts.

// region:proto:start

public protocol Shape {
  func area() -> Double
}

public struct Circle: Shape {
  public let r: Double
  public init(r: Double) { self.r = r }
  public func area() -> Double { Double.pi * r * r }
}

public struct Square: Shape {
  public let s: Double
  public init(s: Double) { self.s = s }
  public func area() -> Double { s * s }
}

// region:proto:end

/// Static dispatch: the generic is specialized, the call can inline.
public func sumStatic<S: Shape>(_ xs: [S]) -> Double {
  xs.reduce(0) { $0 + $1.area() }
}

/// Dynamic dispatch through the existential's witness table.
public func sumDynamic(_ xs: [any Shape]) -> Double {
  xs.reduce(0) { $0 + $1.area() }
}
