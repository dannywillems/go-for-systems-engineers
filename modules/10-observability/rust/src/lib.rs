//! Module 10 in Rust: proving allocation counts with a COUNTING GLOBAL
//! ALLOCATOR. Rust has no `testing.AllocsPerRun`, but it lets you install your
//! own `#[global_allocator]`, so counting allocations is a few lines and is
//! exact and deterministic — the same falsifiable spine as the Go module.
//!
//! Note the cross-language difference this exposes: Rust's `String` is a
//! GROWABLE buffer (amortized doubling), so naive `push_str` in a loop costs
//! O(log n) reallocations, not the O(n) that Go's IMMUTABLE-string `+=` pays.
//! Pre-sizing with `with_capacity` still collapses it to a single allocation.

use std::alloc::{GlobalAlloc, Layout, System};
use std::sync::atomic::{AtomicUsize, Ordering};

// region:counter:start

pub struct Counting;

static ALLOCS: AtomicUsize = AtomicUsize::new(0);

// Wrapping the System allocator and counting every alloc call turns "how many
// allocations did this do" into an exact, in-process number.
unsafe impl GlobalAlloc for Counting {
    unsafe fn alloc(&self, layout: Layout) -> *mut u8 {
        ALLOCS.fetch_add(1, Ordering::Relaxed);
        unsafe { System.alloc(layout) }
    }
    unsafe fn dealloc(&self, ptr: *mut u8, layout: Layout) {
        unsafe { System.dealloc(ptr, layout) }
    }
}

#[global_allocator]
static GLOBAL: Counting = Counting;

// region:counter:end

/// reset the allocation counter and return the count since the last reset.
pub fn take_allocs() -> usize {
    ALLOCS.swap(0, Ordering::Relaxed)
}

/// concat_plus grows a fresh String as it appends: amortized doubling, so
/// O(log n) reallocations.
pub fn concat_plus(parts: &[String]) -> String {
    let mut s = String::new();
    for p in parts {
        s.push_str(p);
    }
    s
}

/// with_cap pre-sizes the String once, so the whole build is a single
/// allocation.
pub fn with_cap(parts: &[String]) -> String {
    let total: usize = parts.iter().map(|p| p.len()).sum();
    let mut s = String::with_capacity(total);
    for p in parts {
        s.push_str(p);
    }
    s
}

#[cfg(test)]
mod tests {
    use super::*;

    fn parts() -> Vec<String> {
        (0..64).map(|i| format!("chunk-{i:02};")).collect()
    }

    // A single test function: the counting allocator is process-global, so
    // running it alongside a second #[test] (which cargo does in parallel)
    // would race the counter. One test keeps the measurement deterministic.
    #[test]
    fn build_correct_and_alloc_bounded() {
        let p = parts();
        assert_eq!(concat_plus(&p), with_cap(&p), "builders disagree");

        take_allocs(); // reset
        let _ = with_cap(&p);
        let cap = take_allocs();
        let _ = concat_plus(&p);
        let plus = take_allocs();
        assert!(cap < plus, "with_cap {cap} should be < concat_plus {plus}");
        assert_eq!(cap, 1, "with_cap should be a single allocation");
    }
}
