// Module 11 in Kotlin: the same Normalize subject. Tests use plain check() to
// stay dependency-free (kotlinc, no Gradle); the production tools are
// kotlin.test / JUnit (units), kotest property testing, and jazzer/JQF
// (fuzzing) -- see the README.

fun normalize(s: String): String =
    s
        .split(Regex("\\s+"))
        .filter { it.isNotEmpty() }
        .joinToString(" ")
        .lowercase()
