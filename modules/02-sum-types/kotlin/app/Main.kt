fun main() {
    val e: Expr = Add(Mul(Lit(2), Lit(3)), Neg(Lit(4)))
    println("eval((2*3) + -4) = ${eval(e)}")
}
