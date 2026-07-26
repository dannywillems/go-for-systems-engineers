import Un

let bits = String(floatBits(1.0), radix: 16, uppercase: true)
print("floatBits(1.0) = 0x\(bits) (withUnsafeBytes)")
print("bytesToString([104,105]) = \(bytesToString([104, 105]))")
