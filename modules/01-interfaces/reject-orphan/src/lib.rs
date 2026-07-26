//! DOES NOT COMPILE, on purpose. Rust's coherence rule (the "orphan rule")
//! forbids implementing a foreign trait for a foreign type: at least one of the
//! trait or the type must be local to this crate. Here both `Display` and
//! `Vec<T>` are foreign, so rustc rejects it with E0117.
//!
//! Go has no orphan rule because a method is declared syntactically WITH its
//! receiver type and interface satisfaction is structural, computed at each use
//! site; there is no global coherence invariant to protect. The trade-off:
//! Rust guarantees at most one impl of a trait for a type program-wide (so
//! trait dispatch is unambiguous), which Go cannot guarantee for, e.g., method
//! promotion via embedding.

use std::fmt;

impl fmt::Display for Vec<u8> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{} bytes", self.len())
    }
}
