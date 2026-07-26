fun main() {
    check(sum(listOf(1, 2, 3, 4, 5), 0) { a, b -> a + b } == 15) { "int sum" }
    check(typeName<String>() == "String") { "reified failed" }
    println("ok")
}
