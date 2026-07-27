//! Rust's showcase quirk: TYPESTATE. A phantom type parameter encodes the
//! object's state in its type, so the compiler makes an invalid transition
//! (calling `open` on an already-open door) a compile error -- the method
//! simply does not exist for that type. See reject-rust.

use std::marker::PhantomData;

// region:typestate:start

pub struct Open;
pub struct Closed;

pub struct Door<State> {
    _state: PhantomData<State>,
}

impl Door<Closed> {
    pub fn closed() -> Self {
        Door {
            _state: PhantomData,
        }
    }
    pub fn open(self) -> Door<Open> {
        Door {
            _state: PhantomData,
        }
    }
}

impl Door<Open> {
    pub fn close(self) -> Door<Closed> {
        Door {
            _state: PhantomData,
        }
    }
}

// region:typestate:end

// region:lenvec:start

// Const generics put a NATURAL NUMBER in the type. `zip_add` requires both
// arrays to have the SAME length N, so a length mismatch is a compile error, not
// a runtime panic -- a length-indexed vector. (See reject-rust-len for the
// rejected mismatched call.)
pub fn zip_add<const N: usize>(a: [i64; N], b: [i64; N]) -> [i64; N] {
    let mut out = [0i64; N];
    for i in 0..N {
        out[i] = a[i] + b[i];
    }
    out
}

// region:lenvec:end

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn valid_transitions_compile() {
        // closed -> open -> closed is the only legal sequence, and it typechecks.
        let _d: Door<Closed> = Door::closed().open().close();
    }

    #[test]
    fn zip_add_same_length() {
        // N is inferred as 3 for both; a [i64;2] second argument would not compile.
        assert_eq!(zip_add([1, 2, 3], [4, 5, 6]), [5, 7, 9]);
    }
}
