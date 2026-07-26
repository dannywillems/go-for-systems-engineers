// DOES NOT COMPILE. `balance` is a private field of `bank::Account`, so reading
// it from outside the `bank` module is rejected:
//
//   error[E0616]: field `balance` of struct `Account` is private

mod bank {
    pub struct Account {
        balance: i64,
    }
    impl Account {
        pub fn open() -> Account {
            Account { balance: 0 }
        }
    }
}

fn main() {
    let a = bank::Account::open();
    println!("{}", a.balance); // private field: rejected
}
