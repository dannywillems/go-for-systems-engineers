//! Prints the two deterministic facts Module 00 checks. stdout is injected into
//! the module README verbatim by the capture tool.

use bootstrap::{sum, word_size_bytes};

const N: u64 = 1_000_000;

fn main() {
    println!("sum(1..{N}) = {}", sum(N));
    println!("word size (bytes) = {}", word_size_bytes());
}
