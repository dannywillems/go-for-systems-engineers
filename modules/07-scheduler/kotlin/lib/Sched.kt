// Kotlin/JVM parallelism here uses platform threads (1:1 OS threads). Kotlin
// also has coroutines (M:N, cooperative, the closest analog to goroutines) and,
// on modern JDKs, virtual threads (Project Loom) — both need dependencies or a
// newer runtime and are noted in the comparison; a CPU sweep uses plain threads.

import kotlin.math.sqrt

fun chunkSum(
    lo: Long,
    hi: Long,
): Double {
    var s = 0.0
    var i = lo
    while (i < hi) {
        s += sqrt(i.toDouble())
        i++
    }
    return s
}

fun parallelSqrtSum(
    total: Long,
    workers: Int,
): Double {
    val chunk = total / workers
    val results = DoubleArray(workers)
    val threads =
        (0 until workers).map { k ->
            Thread { results[k] = chunkSum(k * chunk, (k + 1) * chunk) }
        }
    threads.forEach { it.start() }
    threads.forEach { it.join() }
    return results.sum()
}
