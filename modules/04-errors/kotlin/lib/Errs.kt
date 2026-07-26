// Kotlin's default error channel is UNCHECKED exceptions (no checked-exception
// obligation; the compiler does not force you to handle or declare them, closer
// to Go's "you might forget" than to Rust's forced Result). It also offers a
// `Result<T>` via `runCatching`, with `map`/`fold` chaining for a value-based
// style.

// region:result:start

class CalcException(
    msg: String,
) : Exception(msg)

// Exception-based: nothing at the call site forces handling.
fun chain(x: Int): Int {
    val v = x * 2
    if (v > 100) throw CalcException("too big")
    return v + 1
}

// Value-based: runCatching turns exceptions into a Result you can map/fold.
fun chainResult(x: Int): Result<Int> = runCatching { chain(x) }

// region:result:end
