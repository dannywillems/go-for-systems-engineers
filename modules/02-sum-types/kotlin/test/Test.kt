fun main() {
    val e: Expr = Add(Mul(Lit(2), Lit(3)), Neg(Lit(4)))
    check(eval(e) == 2) { "eval wrong: ${eval(e)}" }
    println("ok")
}
