//! Module 11 in Rust: the same Normalize subject with Rust's test idioms.
//! Built-in `#[test]` units, plus a hand-rolled PROPERTY loop (idempotence over
//! a generated corpus) to stay dependency-free -- the production tools are
//! proptest / quickcheck (property) and cargo-fuzz (fuzzing); see the README.

pub fn normalize(s: &str) -> String {
    s.split_whitespace()
        .collect::<Vec<_>>()
        .join(" ")
        .to_lowercase()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn table() {
        assert_eq!(normalize("  hi  "), "hi");
        assert_eq!(normalize("a\t\n  b"), "a b");
        assert_eq!(normalize("MiXeD"), "mixed");
        assert_eq!(normalize("   "), "");
    }

    // Hand-rolled property: idempotence over a generated corpus. An LCG builds
    // varied strings from a small alphabet (letters, spaces, tabs) so the
    // invariant is checked on many shapes, not a fixed list.
    #[test]
    fn idempotent_property() {
        let alphabet = [b'a', b'B', b' ', b'\t', b'\n', b'c'];
        let mut seed: u64 = 0x9e37_79b9_7f4a_7c15;
        for _ in 0..10_000 {
            seed = seed.wrapping_mul(6364136223846793005).wrapping_add(1);
            let len = (seed >> 60) as usize; // 0..15
            let mut s = String::new();
            let mut x = seed;
            for _ in 0..len {
                x = x.wrapping_mul(6364136223846793005).wrapping_add(1);
                s.push(alphabet[(x >> 61) as usize % alphabet.len()] as char);
            }
            let once = normalize(&s);
            assert_eq!(normalize(&once), once, "not idempotent on {s:?}");
            assert!(!once.contains("  "), "double space in {once:?}");
            assert_eq!(once.trim(), once, "untrimmed {once:?}");
        }
    }
}
