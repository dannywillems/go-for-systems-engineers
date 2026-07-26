// A tiny demo of the subject under test.
fun main() {
    for (s in listOf("  Hello   World  ", "MiXeD\tCase", "   ", "a\n\nb")) {
        println("normalize(${s.replace("\n", "\\n").replace("\t", "\\t")}) = ${normalize(s)}")
    }
}
