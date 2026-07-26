// DOES NOT COMPILE: a `when` over a sealed type, used as an expression, must be
// exhaustive. The missing `Blue` branch is a compile error.

sealed interface Color

object Red : Color

object Green : Color

object Blue : Color

fun name(c: Color): String =
    when (c) {
        is Red -> "red"
        is Green -> "green"
    }
