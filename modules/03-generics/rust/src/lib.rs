//! Rust MONOMORPHIZES: each instantiation `sum::<T>` is compiled to its own
//! specialized machine code, with zero runtime dispatch and no dictionary — the
//! opposite of Go's shared `go.shape.*uint8` pointer stencil. The counter-cost
//! is compile time and binary size that grow with the number of instantiations
//! (measured in the README).

use std::ops::Add;

// region:sum:start

/// Generic over any addable, copyable type with a zero. Monomorphized per T.
pub fn sum<T>(xs: &[T]) -> T
where
    T: Add<Output = T> + Copy + Default,
{
    xs.iter().fold(T::default(), |acc, &x| acc + x)
}

// region:sum:end

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn sums() {
        assert_eq!(sum(&[1i64, 2, 3, 4, 5]), 15);
        assert_eq!(sum(&[1.5f64, 2.5, 3.0]), 7.0);
        assert_eq!(sum(&[1u32, 2, 3]), 6);
    }
}
