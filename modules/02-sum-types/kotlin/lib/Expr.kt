// Kotlin sealed hierarchies are coproducts. A `when` over a sealed type used as
// an EXPRESSION must be exhaustive; a missing branch is a compile error. See
// reject-kotlin. (Nullability replaces Option: `T?` plus compiler null-checks.)

// region:sealed:start

sealed interface Expr

data class Lit(
    val v: Int,
) : Expr

data class Add(
    val l: Expr,
    val r: Expr,
) : Expr

data class Mul(
    val l: Expr,
    val r: Expr,
) : Expr

data class Neg(
    val x: Expr,
) : Expr

fun eval(e: Expr): Int =
    when (e) {
        is Lit -> e.v
        is Add -> eval(e.l) + eval(e.r)
        is Mul -> eval(e.l) * eval(e.r)
        is Neg -> -eval(e.x)
    }

// region:sealed:end
