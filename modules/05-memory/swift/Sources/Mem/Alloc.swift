// Swift uses ARC (reference counting), not a tracing GC: each object is retained
// on creation and released (freed) deterministically when the last reference
// drops. No collector pause, but atomic refcount traffic per reference op. Each
// Node is stored into a rotating slot so it genuinely escapes.

public final class Node {
  var v: Int = 0
  var pad = (0, 0, 0, 0, 0, 0)
}

public func allocBench(_ n: Int) -> Int {
  var keep = [Node?](repeating: nil, count: 16)
  var acc = 0
  for i in 0..<n {
    let node = Node()
    node.v = i
    keep[i & 15] = node
    acc &+= node.v
  }
  return acc
}
