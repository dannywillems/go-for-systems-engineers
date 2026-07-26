//! Allocate 50M short-lived boxed 8-int arrays. Rust has NO GC: each Box is a
//! malloc on creation and a free at end of scope. There is no collector pause,
//! but also no bump-allocation amortization; the allocator does the work inline.
//! black_box prevents the optimizer from eliding the allocation.

use std::hint::black_box;
use std::time::Instant;

fn main() {
    let n: u64 = 50_000_000;
    let mut acc: i64 = 0;
    let t0 = Instant::now();
    for i in 0..n {
        let mut b = Box::new([0i64; 8]);
        b[0] = i as i64;
        acc = acc.wrapping_add(black_box(&b)[0]);
    }
    println!(
        "Rust alloc 50M (Box, no GC): {} ms (acc={})",
        t0.elapsed().as_millis(),
        acc
    );
}
