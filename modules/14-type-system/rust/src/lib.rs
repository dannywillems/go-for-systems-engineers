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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn valid_transitions_compile() {
        // closed -> open -> closed is the only legal sequence, and it typechecks.
        let _d: Door<Closed> = Door::closed().open().close();
    }
}
