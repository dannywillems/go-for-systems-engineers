// The Module 00 fixture in Swift: the same trivial computation as the other four
// languages, so all demo binaries emit byte-identical output.

// region:sum:start

/// `sum(n)` returns 1 + 2 + ... + n. Swift's `Int` is 64-bit on a 64-bit
/// target, so the value fits and matches the Go/Rust/OCaml/Kotlin sides.
public func sum(_ n: Int) -> Int {
  if n < 1 { return 0 }
  var total = 0
  for i in 1...n { total += i }
  return total
}

/// Size of the native word in bytes: 8 on any 64-bit platform.
public func wordSizeBytes() -> Int {
  MemoryLayout<Int>.size
}

// region:sum:end
