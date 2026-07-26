// DOES NOT COMPILE. #[derive(Debug)] generates a Debug impl that formats every
// field, so it requires every field to be Debug. `inner: NotDebug` is not, and
// the derive is rejected at COMPILE time -- reflection would have deferred this
// to run time (or silently produced nothing):
//
//   error[E0277]: `NotDebug` doesn't implement `Debug`

struct NotDebug;

#[derive(Debug)]
struct Wrapper {
    inner: NotDebug,
}

fn main() {
    let w = Wrapper { inner: NotDebug };
    println!("{w:?}");
}
