//! DOES NOT COMPILE. The Go RacyInc shares an unguarded counter across threads;
//! the exact translation is rejected by Rust's type system. `Rc` is not `Send`,
//! so moving it into a spawned thread is E0277 ("cannot be sent between threads
//! safely"). To share mutable state you are FORCED to use `Arc<Mutex<_>>` or an
//! atomic — the synchronization Go leaves optional, Rust makes mandatory.

use std::rc::Rc;
use std::thread;

fn main() {
    let counter = Rc::new(0);
    let handles: Vec<_> = (0..4)
        .map(|_| {
            let c = Rc::clone(&counter);
            thread::spawn(move || {
                let _ = *c; // Rc<i32> is !Send -> cannot cross the thread boundary
            })
        })
        .collect();
    for h in handles {
        h.join().unwrap();
    }
}
