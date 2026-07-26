fun main() {
    val a = Person("Ada", 36)
    check(a == Person("Ada", 36)) { "generated equals" }
    check(a.copy(age = 37).age == 37) { "generated copy" }
    check(a.toString() == "Person(name=Ada, age=36)") { "generated toString: $a" }
    val (name, age) = a // generated componentN (destructuring)
    check(name == "Ada" && age == 36) { "generated componentN" }
    println("ok")
}
