use conc::atomic_count;

fn main() {
    println!(
        "atomic_count(8, 100000) = {} (correct)",
        atomic_count(8, 100_000)
    );
}
