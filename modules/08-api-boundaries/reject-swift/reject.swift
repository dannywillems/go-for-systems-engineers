// DOES NOT COMPILE. `balance` is private to Account, so touching it from outside
// the struct is rejected:
//
//   error: 'balance' is inaccessible due to 'private' protection level

struct Account {
    private var balance: Int64 = 0
}

let a = Account()
print(a.balance)
