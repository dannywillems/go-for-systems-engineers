// The Module 00 fixture in Kotlin (JVM). Same computation as the other four
// languages so the demo prints byte-identical output. Built with kotlinc
// directly (no Gradle) to keep the repo self-contained; Gradle is the
// production build tool.

// region:sum:start

/**
 * Returns 1 + 2 + ... + n. On the JVM `Int` is 32-bit, which overflows here, so
 * the accumulator is a 64-bit `Long`. This 32-bit-Int fact is itself a
 * representational divergence from Go/Rust/Swift (64-bit Int) worth noting.
 */
fun sum(n: Int): Long {
    var total = 0L
    for (i in 1..n) total += i
    return total
}

/**
 * Native word size in bytes: 8 on a 64-bit JVM. There is no `sizeof` on the
 * JVM, so this reads the data model the JVM runs under ("64" -> 8 bytes).
 */
fun wordSizeBytes(): Int = System.getProperty("sun.arch.data.model").toInt() / 8

// region:sum:end
