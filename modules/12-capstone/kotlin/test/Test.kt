fun main() {
    val c = Cache(16, 4, 0)
    val threads =
        (0 until 8).map { w ->
            Thread {
                repeat(500) { i ->
                    val key = (w * 500 + i) % 100
                    check(c.get(key) == key * key) { "wrong value for $key" }
                }
            }
        }
    threads.forEach { it.start() }
    threads.forEach { it.join() }
    check(c.size() <= 16) { "cache exceeded capacity: ${c.size()}" }
    check(c.backendCalls() < 8 * 500) { "cache did nothing" }
    println("ok")
}
