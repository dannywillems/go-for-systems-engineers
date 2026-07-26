// Dependency-free assertion test: exits non-zero (via AssertionError) on any
// failure, so `java -jar test.jar` is the pass/fail signal. Kotlin's kotlin.test
// + JUnit would pull dependencies we deliberately avoid in this self-contained
// kotlinc build; Gradle + JUnit is the production choice.

fun main() {
    for (n in listOf(0, 1, 10, 1_000_000)) {
        check(sum(n) == n.toLong() * (n + 1) / 2) { "sum($n) wrong" }
    }
    check(wordSizeBytes() == 8) { "word size wrong" }
    println("ok")
}
