import Calc

let e: Expr = .add(.mul(.lit(2), .lit(3)), .neg(.lit(4)))
print("eval((2*3) + -4) = \(eval(e))")
