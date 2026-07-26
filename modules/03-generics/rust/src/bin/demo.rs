use genr::sum;

fn main() {
    // Each call site instantiates (monomorphizes) sum for a distinct type.
    println!("sum::<i64>  = {}", sum(&[1i64, 2, 3, 4, 5]));
    println!("sum::<f64>  = {}", sum(&[1.5f64, 2.5, 3.0]));
    println!("sum::<u32>  = {}", sum(&[1u32, 2, 3]));
}
