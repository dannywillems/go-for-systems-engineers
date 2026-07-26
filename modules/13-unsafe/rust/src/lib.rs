//! Rust's `unsafe`: raw pointers, `*_unchecked`, and `transmute`. Crucially the
//! borrow checker and type checker STILL run inside `unsafe`; it only unlocks a
//! few extra operations (deref raw pointers, call unsafe fns, union field
//! access, etc.). Undefined behavior these enable is caught by Miri
//! (`cargo +nightly miri test`), not the normal compiler.

// region:unsafe:start

/// Zero-copy: reinterpret bytes as &str WITHOUT the UTF-8 validation `str::from_utf8`
/// does. Sound only if the caller guarantees valid UTF-8; otherwise UB.
pub fn bytes_to_str(b: &[u8]) -> &str {
    unsafe { std::str::from_utf8_unchecked(b) }
}

/// Read the raw bits of an f32 as a u32 by deref of a raw pointer.
pub fn f32_bits(f: f32) -> u32 {
    let p = &f as *const f32 as *const u32;
    unsafe { *p }
}

// region:unsafe:end

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn conversions() {
        assert_eq!(bytes_to_str(b"hello"), "hello");
        // 1.0f32 has IEEE-754 bit pattern 0x3F800000.
        assert_eq!(f32_bits(1.0), 0x3F80_0000);
    }
}
