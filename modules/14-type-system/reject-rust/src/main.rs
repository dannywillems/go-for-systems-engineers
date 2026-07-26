//! DOES NOT COMPILE: `open()` exists only for `Door<Closed>`. Calling it on a
//! `Door<Open>` is a compile error (no such method), so an invalid state
//! transition is caught statically.

use std::marker::PhantomData;

struct Open;
struct Closed;

struct Door<State> {
    _state: PhantomData<State>,
}

impl Door<Closed> {
    fn closed() -> Self {
        Door { _state: PhantomData }
    }
    fn open(self) -> Door<Open> {
        Door { _state: PhantomData }
    }
}

fn main() {
    let open_door: Door<Open> = Door::closed().open();
    let _ = open_door.open(); // error: no method `open` on Door<Open>
}
