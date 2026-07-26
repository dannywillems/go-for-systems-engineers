// Module 08 in Kotlin: visibility is private / internal / protected / public
// (the default). The unit for `internal` is the MODULE (compilation unit).
// Account is public but its backing field is private, so callers use the
// methods and cannot read or construct balance.

class Overdraft : Exception()

class Account private constructor(
    private var balance: Long,
) {
    companion object {
        fun open(initial: Long): Account {
            if (initial < 0) throw Overdraft()
            return Account(initial)
        }
    }

    fun deposit(amount: Long) {
        balance += amount
    }

    fun withdraw(amount: Long) {
        if (amount > balance) throw Overdraft()
        balance -= amount
    }

    fun balance(): Long = balance
}
