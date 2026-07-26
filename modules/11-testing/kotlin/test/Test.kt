fun main() {
    // table
    check(normalize("  hi  ") == "hi")
    check(normalize("a\t\n  b") == "a b")
    check(normalize("MiXeD") == "mixed")
    check(normalize("   ") == "")

    // hand-rolled property: idempotence over a generated corpus (LCG).
    val alphabet = charArrayOf('a', 'B', ' ', '\t', '\n', 'c')
    var seed = 0x9e3779b97f4a7c15uL
    repeat(10_000) {
        seed = seed * 6364136223846793005uL + 1uL
        val len = (seed shr 60).toInt()
        val sb = StringBuilder()
        var x = seed
        repeat(len) {
            x = x * 6364136223846793005uL + 1uL
            sb.append(alphabet[((x shr 61).toInt()) % alphabet.size])
        }
        val once = normalize(sb.toString())
        check(normalize(once) == once) { "not idempotent on $sb" }
        check(once == once.trim()) { "untrimmed: $once" }
    }
    println("ok")
}
