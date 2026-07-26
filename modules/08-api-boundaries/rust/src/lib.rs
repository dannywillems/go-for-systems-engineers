//! Module 08 in Rust: encapsulation is by the MODULE tree, with graduated
//! visibility. A field is private by default; `pub` exposes it, `pub(crate)`
//! limits it to the crate, `pub(super)` to the parent module. Unlike Go's
//! by-case rule, Rust hides at field granularity: `Account` is `pub` but its
//! `balance` field is private, so other modules hold and use it but cannot read
//! or construct its representation. This is the same "abstract type" as Go's
//! all-unexported-fields struct, but with finer control over each item.

use std::fmt;

#[derive(Debug)]
pub struct Overdraft;

impl fmt::Display for Overdraft {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "overdraft")
    }
}
impl std::error::Error for Overdraft {}

// balance is private: no field on `pub` visibility, so it cannot be read,
// written, or built with a struct literal from another module.
pub struct Account {
    balance: i64,
}

impl Account {
    pub fn open(initial: i64) -> Result<Account, Overdraft> {
        if initial < 0 {
            return Err(Overdraft);
        }
        Ok(Account { balance: initial })
    }
    pub fn deposit(&mut self, amount: i64) {
        self.balance += amount;
    }
    pub fn withdraw(&mut self, amount: i64) -> Result<(), Overdraft> {
        if amount > self.balance {
            return Err(Overdraft);
        }
        self.balance -= amount;
        Ok(())
    }
    pub fn balance(&self) -> i64 {
        self.balance
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn through_public_api() {
        let mut a = Account::open(100).unwrap();
        a.deposit(50);
        assert_eq!(a.balance(), 150);
        assert!(a.withdraw(200).is_err());
        assert_eq!(a.balance(), 150);
        assert!(Account::open(-1).is_err());
    }
}
