fun main() {
    val total = 400_000_000L
    val w = Runtime.getRuntime().availableProcessors()
    val t0 = System.nanoTime()
    val acc = parallelSqrtSum(total, w)
    if (acc.isNaN()) println("nan")
    val ms = (System.nanoTime() - t0) / 1_000_000
    println("Kotlin sqrt-sum ${total / 1_000_000}M / $w threads: $ms ms")
}
