//! DOES NOT COMPILE: `match` on an enum must be total. The missing `Blue` arm
//! is rejected at compile time with E0004 — the exact check Go omits and the
//! ./go/exhaustive analyzer has to reconstruct.

pub enum Color {
    Red,
    Green,
    Blue,
}

pub fn name(c: Color) -> &'static str {
    match c {
        Color::Red => "red",
        Color::Green => "green",
    }
}
