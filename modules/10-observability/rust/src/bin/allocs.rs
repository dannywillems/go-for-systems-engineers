// Prints the deterministic allocation counts (via the counting global
// allocator) and a wall-clock timing for each builder. The counts are
// machine-independent; the timing feeds measured.txt.

use std::time::Instant;

use observability::{concat_plus, take_allocs, with_cap};

fn main() {
    let parts: Vec<String> = (0..64).map(|i| format!("chunk-{i:02};")).collect();

    take_allocs();
    let _ = concat_plus(&parts);
    let plus = take_allocs();
    let _ = with_cap(&parts);
    let cap = take_allocs();

    println!("concat_plus (64 parts): {plus} allocs");
    println!("with_cap    (64 parts): {cap} allocs");

    // timing (non-portable)
    let iters = 1_000_000u32;
    let t = Instant::now();
    for _ in 0..iters {
        std::hint::black_box(concat_plus(std::hint::black_box(&parts)));
    }
    let plus_ns = t.elapsed().as_nanos() / iters as u128;
    let t = Instant::now();
    for _ in 0..iters {
        std::hint::black_box(with_cap(std::hint::black_box(&parts)));
    }
    let cap_ns = t.elapsed().as_nanos() / iters as u128;
    println!(
        "Rust concat_plus: {plus_ns} ns/op ({plus} allocs)  with_cap: {cap_ns} ns/op ({cap} allocs)"
    );
}
