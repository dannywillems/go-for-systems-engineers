use un::{bytes_to_str, f32_bits};

fn main() {
    println!(
        "bytes_to_str(b\"hello\") = {} (unchecked, no copy)",
        bytes_to_str(b"hello")
    );
    println!("f32_bits(1.0) = 0x{:08X}", f32_bits(1.0));
}
