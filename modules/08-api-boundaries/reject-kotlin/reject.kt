// DOES NOT COMPILE. `balance` is private to Account, so reading it from outside
// the class is rejected:
//
//   error: cannot access 'var balance: Long': it is private in 'Account'.

class Account {
    private var balance: Long = 0
}

fun main() {
    val a = Account()
    println(a.balance)
}
