// Kotlin's showcase quirk: DECLARATION-SITE VARIANCE. `out T` marks T covariant
// (a producer), `in T` contravariant (a consumer); the compiler then allows the
// safe subtype substitutions and rejects the unsafe ones. See reject-kotlin for
// the rejected case (using an `out` parameter in an `in` position).

open class Animal(
    val name: String,
)

class Cat : Animal("cat")

// Covariant: a Producer<Cat> is usable as a Producer<Animal>.
interface Producer<out T> {
    fun produce(): T
}

// Because of `out`, this accepts Producer<Cat> as well.
fun firstName(p: Producer<Animal>): String = p.produce().name

val catProducer =
    object : Producer<Cat> {
        override fun produce(): Cat = Cat()
    }
