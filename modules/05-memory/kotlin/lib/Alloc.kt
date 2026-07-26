// Allocate short-lived objects on the JVM (tracing GC). NOTE: the JVM's JIT does
// escape analysis and can SCALAR-REPLACE a non-escaping allocation, eliminating
// it entirely — so a naive "allocate in a loop" can measure the optimizer, not
// the GC. We store each object into a rotating array slot so it genuinely
// escapes and the allocation actually happens.

class Node {
    var v = 0
    val pad = LongArray(6)
}

fun allocBench(n: Int): Long {
    val keep = arrayOfNulls<Node>(16)
    var acc = 0L
    for (i in 0 until n) {
        val node = Node()
        node.v = i
        keep[i and 15] = node // escapes: defeats scalar replacement
        acc += node.v
    }
    return acc
}
