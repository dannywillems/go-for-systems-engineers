fun main() {
    val capacity = 256
    val maxInflight = 32
    val keys = 256L
    val workers = 64
    val perWorker = 10_000

    val cache = Cache(capacity, maxInflight, 100)
    val all = LongArray(workers * perWorker)
    val t0 = System.nanoTime()
    val threads =
        (0 until workers).map { w ->
            Thread {
                var seed = w.toLong() * 2654435761L or 1L
                var idx = w * perWorker
                repeat(perWorker) {
                    seed = seed * 6364136223846793005L + 1442695040888963407L
                    val key = ((seed ushr 33) % keys).toInt()
                    val s = System.nanoTime()
                    cache.get(key)
                    all[idx++] = System.nanoTime() - s
                }
            }
        }
    threads.forEach { it.start() }
    threads.forEach { it.join() }
    val elapsedS = (System.nanoTime() - t0) / 1e9

    all.sort()
    val n = all.size

    fun pc(p: Double): Double = all[minOf((p / 100 * n).toInt(), n - 1)] / 1000.0
    val backendPct = 100.0 * cache.backendCalls() / n
    println(
        "Kotlin %dk gets/%dw: %.0f kops/s  p50=%.1fus p99=%.1fus p999=%.1fus  backend=%.1f%% of gets"
            .format(n / 1000, workers, n / elapsedS / 1000, pc(50.0), pc(99.0), pc(99.9), backendPct),
    )
}
