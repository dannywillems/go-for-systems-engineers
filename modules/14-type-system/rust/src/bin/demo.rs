use ts::Door;

fn main() {
    // The only legal transition sequence; each step returns a different type.
    let _door = Door::closed().open().close();
    println!("Door: closed -> open -> close typechecks (typestate)");
}
