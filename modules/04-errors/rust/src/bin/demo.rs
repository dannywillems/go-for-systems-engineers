use errs::{chain, use_question};

fn main() {
    println!(
        "Ok(3).map(*2).and_then(+1) = {:?}  (method chaining)",
        chain(3)
    );
    println!("chain(60) = {:?}", chain(60));
    println!("use_question(3) with ? = {:?}", use_question(3));
}
