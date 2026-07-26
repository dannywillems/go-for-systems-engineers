// DOES NOT COMPILE. A data class exists only to generate members FROM its
// primary constructor properties, so one with no properties is rejected:
//
//   error: data class must have at least one primary constructor parameter

data class Empty()

fun main() {
    println(Empty())
}
