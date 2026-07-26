use sched::parallel_sqrt_sum;
use std::thread;
use std::time::Instant;

const TOTAL: u64 = 400_000_000;

fn main() {
    let w = thread::available_parallelism()
        .map(|n| n.get())
        .unwrap_or(4) as u64;
    let t0 = Instant::now();
    let acc = parallel_sqrt_sum(TOTAL, w);
    let _ = acc;
    println!(
        "Rust   sqrt-sum {}M / {} threads: {} ms",
        TOTAL / 1_000_000,
        w,
        t0.elapsed().as_millis()
    );
}
