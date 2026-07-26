use expr::{Expr, niche_sizes};

fn main() {
    let e = Expr::Add(
        Box::new(Expr::Mul(Box::new(Expr::Lit(2)), Box::new(Expr::Lit(3)))),
        Box::new(Expr::Neg(Box::new(Expr::Lit(4)))),
    );
    println!("eval((2*3) + -4) = {}", e.eval());
    for (name, sz) in niche_sizes() {
        println!("size_of::<{name}> = {sz}");
    }
}
