//! Rust HAS coproducts: `enum` is a real tagged union and `match` is a total
//! eliminator the compiler checks. A missing arm is a compile error (E0004),
//! not a silent fall-through. See the `reject-rust` crate for the rejection.

// region:enum:start

pub enum Expr {
    Lit(i64),
    Add(Box<Expr>, Box<Expr>),
    Mul(Box<Expr>, Box<Expr>),
    Neg(Box<Expr>),
}

impl Expr {
    pub fn eval(&self) -> i64 {
        match self {
            Expr::Lit(v) => *v,
            Expr::Add(l, r) => l.eval() + r.eval(),
            Expr::Mul(l, r) => l.eval() * r.eval(),
            Expr::Neg(x) => -x.eval(),
        }
    }
}

// region:enum:end

/// Niche optimization: `Option<T>` reuses an invalid bit pattern of `T` as the
/// `None` tag when one exists, so `Option<&T>` and `Option<Box<T>>` are the
/// SAME size as the pointer. This is why Rust needs no nullable-pointer type.
pub fn niche_sizes() -> [(&'static str, usize); 4] {
    [
        ("&u8", size_of::<&u8>()),
        ("Option<&u8>", size_of::<Option<&u8>>()),
        ("Option<Box<u8>>", size_of::<Option<Box<u8>>>()),
        ("Option<u8>", size_of::<Option<u8>>()), // no niche: 1 tag + 1 byte
    ]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn eval_works() {
        // (2 * 3) + (-4) = 2
        let e = Expr::Add(
            Box::new(Expr::Mul(Box::new(Expr::Lit(2)), Box::new(Expr::Lit(3)))),
            Box::new(Expr::Neg(Box::new(Expr::Lit(4)))),
        );
        assert_eq!(e.eval(), 2);
    }

    #[test]
    fn option_of_pointer_is_free() {
        assert_eq!(size_of::<Option<&u8>>(), size_of::<&u8>());
        assert_eq!(size_of::<Option<Box<u8>>>(), size_of::<Box<u8>>());
        // But a value type with no spare bit pattern pays a tag byte:
        assert_eq!(size_of::<Option<u8>>(), 2);
    }
}
