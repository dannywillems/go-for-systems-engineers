fun main() {
    // 1+2+...+999 with the low 4 bits rotating is just the sum; check a small n.
    val acc = allocBench(1000)
    check(acc == (0 until 1000).sumOf { it.toLong() }) { "allocBench sum wrong: $acc" }
    println("ok")
}
