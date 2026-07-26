// Module 09 in Kotlin: a `data class` makes the COMPILER GENERATE equals,
// hashCode, toString, copy, and componentN (destructuring) from the primary
// constructor properties -- compile-time codegen, no runtime reflection. It is
// the everyday Kotlin analogue of Rust's derive and Swift's synthesis.

data class Person(
    val name: String,
    val age: Int,
)
