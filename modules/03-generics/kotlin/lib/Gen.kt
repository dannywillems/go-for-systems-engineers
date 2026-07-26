// Kotlin generics are ERASED on the JVM (Java-style type erasure): one shared
// bytecode implementation, all type arguments gone at runtime, elements boxed.
// The escape hatch is `inline`+`reified`, which inlines the generic at each call
// site and thus DOES have the concrete type available (a form of specialization,
// closer to Rust). Below, a normal erased generic and a reified one.

// region:generic:start

fun <T> sum(
    xs: List<T>,
    zero: T,
    add: (T, T) -> T,
): T = xs.fold(zero, add)

// reified: the type T survives because the function is inlined at each call.
inline fun <reified T> typeName(): String = T::class.simpleName ?: "?"

// region:generic:end
