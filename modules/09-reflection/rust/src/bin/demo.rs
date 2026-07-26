// Shows the compiler-generated impls at work: Debug prints, PartialEq compares,
// Clone duplicates -- all generated at compile time by #[derive]. Deterministic.

use reflectgen::Person;

fn main() {
    let a = Person::new("Ada", 36);
    println!("Debug (derived):      {a:?}");
    println!("PartialEq (derived):  {}", a == Person::new("Ada", 36));
    let mut b = a.clone(); // Clone (derived)
    b.age = 37;
    println!("after clone+edit:     {b:?}");
    println!("still equal?          {}", a == b);
}
