// Prints the two deterministic facts Module 00 checks. stdout is injected into
// the module README verbatim by the capture tool.

fun main() {
    val n = 1_000_000
    println("sum(1..$n) = ${sum(n)}")
    println("word size (bytes) = ${wordSizeBytes()}")
}
