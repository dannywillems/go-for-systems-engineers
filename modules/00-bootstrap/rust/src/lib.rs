//! The Module 00 fixture in Rust: the same trivial computation as the Go and
//! OCaml sides, so the three demo binaries emit byte-identical output.

// region:sum:start
/// `sum(n)` returns 1 + 2 + ... + n. The result is identical on every 64-bit
/// target and in every language, which makes it a clean cross-toolchain fixture.
pub fn sum(n: u64) -> u64 {
    (1..=n).sum()
}

/// Size of a pointer-width integer in bytes: 8 on any 64-bit platform.
pub fn word_size_bytes() -> usize {
    std::mem::size_of::<usize>()
}
// region:sum:end

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn sum_matches_closed_form() {
        for n in [0u64, 1, 10, 1_000_000] {
            assert_eq!(sum(n), n * (n + 1) / 2);
        }
    }

    #[test]
    fn word_size_is_eight() {
        assert_eq!(word_size_bytes(), 8);
    }
}
