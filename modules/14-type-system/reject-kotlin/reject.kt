// DOES NOT COMPILE: `T` is declared `out` (covariant, producer-only) but is used
// in an `in` position (a function parameter, a consumer). Kotlin rejects this
// because it would break covariance soundness.

interface Bad<out T> {
    fun consume(item: T) // error: type parameter T is declared as 'out' but occurs in 'in' position
}
