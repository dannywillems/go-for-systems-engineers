// DOES NOT COMPILE: `open()` exists only in the `where State == Closed`
// extension, so calling it on a Door<Open> is a compile error.

enum Open {}
enum Closed {}

struct Door<State> {}

extension Door where State == Closed {
  func open() -> Door<Open> { Door<Open>() }
}

func bad() {
  let openDoor: Door<Open> = Door<Closed>().open()
  _ = openDoor.open()  // error: 'Door<Open>' has no member 'open'
}
