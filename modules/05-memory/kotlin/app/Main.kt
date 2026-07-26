fun main() {
    val n = 50_000_000
    val t0 = System.nanoTime()
    val acc = allocBench(n)
    val ms = (System.nanoTime() - t0) / 1_000_000
    println("Kotlin alloc 50M (JVM GC): $ms ms (acc=$acc)")
}
