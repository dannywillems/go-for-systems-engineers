// DOES NOT COMPILE. `zip_add` unifies the length N of both arrays in the type,
// so passing a [i64; 3] and a [i64; 2] is a type error -- the length lives in
// the type and the compiler checks it:
//
//   error[E0308]: mismatched types
//   expected an array with a size of 3, found one with a size of 2

fn zip_add<const N: usize>(a: [i64; N], _b: [i64; N]) -> [i64; N] {
    a
}

fn main() {
    let _ = zip_add([1, 2, 3], [4, 5]); // lengths 3 and 2 disagree
}
