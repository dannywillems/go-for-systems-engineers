// Prints bytes allocated by each builder plus a wall-clock timing for
// measured.txt. Byte counts are stable; timings are not (and the JVM needs
// warmup, so these are indicative, not JMH-grade).

fun main() {
    val parts = (0 until 64).map { "chunk-%02d;".format(it) }

    val plusBytes = bytesAllocated { concatPlus(parts) }
    val builderBytes = bytesAllocated { builder(parts) }
    println("concatPlus (64 parts): $plusBytes bytes allocated")
    println("builder    (64 parts): $builderBytes bytes allocated")

    // warm up, then time
    repeat(100_000) {
        concatPlus(parts)
        builder(parts)
    }
    val iters = 1_000_000
    var t = System.nanoTime()
    repeat(iters) { concatPlus(parts) }
    val plusNs = (System.nanoTime() - t) / iters
    t = System.nanoTime()
    repeat(iters) { builder(parts) }
    val builderNs = (System.nanoTime() - t) / iters
    println(
        "Kotlin concatPlus: $plusNs ns/op ($plusBytes B)  builder: $builderNs ns/op ($builderBytes B)",
    )
}
