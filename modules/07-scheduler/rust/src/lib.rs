//! CPU-bound work split across std threads (1:1 OS threads). Rust's async
//! (tokio) is an M:N work-stealing scheduler like Go's; for a pure CPU sweep,
//! OS threads are the direct comparison.

use std::thread;

pub fn chunk_sum(lo: u64, hi: u64) -> f64 {
    let mut s = 0.0f64;
    for i in lo..hi {
        s += (i as f64).sqrt();
    }
    s
}

pub fn parallel_sqrt_sum(total: u64, workers: u64) -> f64 {
    let chunk = total / workers;
    let handles: Vec<_> = (0..workers)
        .map(|k| thread::spawn(move || chunk_sum(k * chunk, (k + 1) * chunk)))
        .collect();
    handles.into_iter().map(|h| h.join().unwrap()).sum()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parallel_matches_serial() {
        assert!((parallel_sqrt_sum(400, 4) - chunk_sum(0, 400)).abs() < 1e-6);
    }
}
