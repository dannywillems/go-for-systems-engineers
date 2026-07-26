fun main() {
    println("sum(ints)   = ${sum(listOf(1, 2, 3, 4, 5), 0) { a, b -> a + b }}")
    println("sum(doubles)= ${sum(listOf(1.5, 2.5, 3.0), 0.0) { a, b -> a + b }}")
    println("reified typeName<Int>() = ${typeName<Int>()}")
}
