//! Rust's answer to Go interfaces: traits, with a hard nominal/coherence
//! discipline. `dyn Trait` is the existential (a fat pointer: data + vtable);
//! a trait bound `<T: Shape>` is universal quantification, monomorphized away.

// region:trait:start

/// A trait is satisfied only by an explicit `impl` (nominal), unlike Go's
/// structural satisfaction. Coherence (the orphan rule) further restricts WHERE
/// that impl may live: see the `reject-orphan` crate for the rejection.
pub trait Shape {
    fn area(&self) -> f64;
}

pub struct Circle {
    pub r: f64,
}
impl Shape for Circle {
    fn area(&self) -> f64 {
        std::f64::consts::PI * self.r * self.r
    }
}

pub struct Square {
    pub s: f64,
}
impl Shape for Square {
    fn area(&self) -> f64 {
        self.s * self.s
    }
}

// region:trait:end

/// Static dispatch: `T: Shape` is monomorphized per concrete type, so the call
/// inlines. This is the universally-quantified path (no vtable).
pub fn sum_static<T: Shape>(xs: &[T]) -> f64 {
    xs.iter().map(Shape::area).sum()
}

/// Dynamic dispatch: `dyn Shape` is the existential; each call goes through the
/// vtable in the fat pointer. Analogous to Go's itab path.
pub fn sum_dynamic(xs: &[&dyn Shape]) -> f64 {
    xs.iter().map(|s| s.area()).sum()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn dispatch_paths_agree() {
        let cs = [Circle { r: 1.0 }, Circle { r: 2.0 }];
        let dyns: [&dyn Shape; 2] = [&cs[0], &cs[1]];
        assert!((sum_static(&cs) - sum_dynamic(&dyns)).abs() < 1e-9);
    }

    #[test]
    fn dyn_pointer_is_fat() {
        // A &dyn Shape is two words (data + vtable); a &Circle is one.
        assert_eq!(size_of::<&dyn Shape>(), 16);
        assert_eq!(size_of::<&Circle>(), 8);
    }
}
