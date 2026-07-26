// Swift's unsafe surface is explicit and scoped: `withUnsafeBytes`,
// `UnsafePointer`, `unsafeBitCast`. The safe API is preferred; these are the
// controlled hatch for reinterpreting memory. Reads the raw IEEE-754 bits of a
// Float via a pointer, mirroring the Go/Rust bit-punning demos.

public func floatBits(_ f: Float) -> UInt32 {
  withUnsafeBytes(of: f) { raw in
    raw.load(as: UInt32.self)
  }
}

// Zero-copy-ish: interpret UTF-8 bytes as a String (this one copies in Swift's
// model, but shows the withUnsafeBufferPointer hatch producing the contents).
public func bytesToString(_ bytes: [UInt8]) -> String {
  bytes.withUnsafeBufferPointer { buf in
    String(decoding: buf, as: UTF8.self)
  }
}
