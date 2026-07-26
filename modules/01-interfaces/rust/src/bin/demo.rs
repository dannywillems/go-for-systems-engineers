//! Prints the deterministic representation facts, mirroring the Go demo.

use shapes::{Circle, Shape, Square};

fn main() {
    let c = Circle { r: 1.0 };
    let s = Square { s: 2.0 };
    let dyns: [&dyn Shape; 2] = [&c, &s];

    println!("a Circle used as &dyn Shape has area {:.4}", dyns[0].area());
    println!("sizeof(&dyn Shape) = {} bytes", size_of::<&dyn Shape>());
    println!("sizeof(&Circle)    = {} bytes", size_of::<&Circle>());
    println!("sum via dyn dispatch = {:.4}", shapes::sum_dynamic(&dyns));
}
