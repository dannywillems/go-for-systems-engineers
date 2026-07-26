use capstone::Cache;
use std::sync::Arc;
use std::time::{Duration, Instant};

const CAPACITY: usize = 256;
const MAX_INFLIGHT: usize = 32;
const KEYS: i64 = 256;
const WORKERS: usize = 64;
const PER_WORKER: usize = 10_000;

fn main() {
    let cache = Arc::new(Cache::new(
        CAPACITY,
        MAX_INFLIGHT,
        Duration::from_micros(100),
    ));
    let start = Instant::now();
    let handles: Vec<_> = (0..WORKERS)
        .map(|w| {
            let cache = Arc::clone(&cache);
            std::thread::spawn(move || {
                let mut lat = Vec::with_capacity(PER_WORKER);
                let mut seed = ((w as u64) * 2654435761) | 1;
                for _ in 0..PER_WORKER {
                    seed = seed
                        .wrapping_mul(6364136223846793005)
                        .wrapping_add(1442695040888963407);
                    let key = ((seed >> 33) % KEYS as u64) as i64;
                    let s = Instant::now();
                    cache.get(key);
                    lat.push(s.elapsed());
                }
                lat
            })
        })
        .collect();

    let mut all: Vec<Duration> = Vec::with_capacity(WORKERS * PER_WORKER);
    for h in handles {
        all.extend(h.join().unwrap());
    }
    let elapsed = start.elapsed();
    all.sort_unstable();
    let n = all.len();
    let pc = |p: f64| all[((p / 100.0 * n as f64) as usize).min(n - 1)];
    println!(
        "Rust   {}k gets/{}w: {:.0} kops/s  p50={:?} p99={:?} p999={:?}  backend={:.1}% of gets",
        n / 1000,
        WORKERS,
        n as f64 / elapsed.as_secs_f64() / 1000.0,
        pc(50.0),
        pc(99.0),
        pc(99.9),
        100.0 * cache.backend_calls() as f64 / n as f64
    );
}
