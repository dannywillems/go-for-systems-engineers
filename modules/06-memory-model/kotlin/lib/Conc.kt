// Kotlin/JVM shares Go's model: no compile-time data-race check. Safety comes
// from the java.util.concurrent primitives (AtomicInteger, volatile,
// synchronized) under the Java Memory Model. Coroutines add structured
// concurrency on top but do not change the memory model; they need
// kotlinx-coroutines (a Gradle dependency) and so are omitted from this
// self-contained kotlinc build.

import java.util.concurrent.atomic.AtomicInteger

// region:atomic:start

fun atomicCount(
    threads: Int,
    per: Int,
): Int {
    val counter = AtomicInteger(0)
    val ts =
        (1..threads).map {
            Thread {
                repeat(per) { counter.incrementAndGet() }
            }
        }
    ts.forEach { it.start() }
    ts.forEach { it.join() }
    return counter.get()
}

// region:atomic:end
