//! Rust makes shared mutable state across threads type-safe: `Send`/`Sync` are
//! auto-traits the compiler tracks, so the only way to share a counter is
//! `Arc<Mutex<_>>` or an atomic. The data race is unrepresentable, not merely
//! undetected. Contrast reject-race, which fails to compile.

use std::sync::Arc;
use std::sync::atomic::{AtomicU64, Ordering};
use std::thread;

// region:atomic:start

/// The synchronized counter: `Arc<AtomicU64>` is `Send + Sync`, so sharing it
/// across threads compiles, and the result is always deterministic.
pub fn atomic_count(threads: usize, per: u64) -> u64 {
    let counter = Arc::new(AtomicU64::new(0));
    let handles: Vec<_> = (0..threads)
        .map(|_| {
            let c = Arc::clone(&counter);
            thread::spawn(move || {
                for _ in 0..per {
                    c.fetch_add(1, Ordering::Relaxed);
                }
            })
        })
        .collect();
    for h in handles {
        h.join().unwrap();
    }
    counter.load(Ordering::SeqCst)
}

// region:atomic:end

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn counts_exactly() {
        assert_eq!(atomic_count(8, 100_000), 800_000);
    }
}
