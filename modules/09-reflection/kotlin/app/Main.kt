// Shows the generated members: toString, equals, and copy. Deterministic.
fun main() {
    val a = Person("Ada", 36)
    println("toString (generated): $a")
    println("equals (generated):   ${a == Person("Ada", 36)}")
    val b = a.copy(age = 37) // copy (generated)
    println("copy (generated):     $b")
    println("not equal after copy: ${a == b}")
}
