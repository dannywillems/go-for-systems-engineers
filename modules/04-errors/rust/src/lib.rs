//! Rust's `Result<T, E>` is a coproduct: exactly one of Ok/Err, never both. It
//! carries `#[must_use]`, so ignoring it warns. Because it is a real type with
//! inherent methods, `map`/`and_then` chain, and `?` propagates. This is the
//! ergonomics Go's hand-rolled Result cannot reach (no generic methods).

// region:result:start

/// The SAME computation as the Go demo, but written as a left-to-right chain
/// (`map` then `and_then`), which Go's free-function Result cannot express.
pub fn chain(x: i32) -> Result<i32, String> {
    Ok(x).map(|v| v * 2).and_then(|v| {
        if v > 100 {
            Err("too big".into())
        } else {
            Ok(v + 1)
        }
    })
}

/// `?` propagates the Err early; the happy path stays at the left margin.
pub fn use_question(x: i32) -> Result<i32, String> {
    let doubled = chain(x)?;
    Ok(doubled + 100)
}

// region:result:end

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn chain_and_question() {
        assert_eq!(chain(3), Ok(7));
        assert_eq!(chain(60), Err("too big".to_string())); // 120 > 100
        assert_eq!(use_question(3), Ok(107));
        assert_eq!(use_question(60), Err("too big".to_string()));
    }
}
