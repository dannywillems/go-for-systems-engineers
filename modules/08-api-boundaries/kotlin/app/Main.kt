// A tiny demo of the opaque Account: callers use only its public methods.
fun main() {
    val a = Account.open(100)
    a.deposit(50)
    a.withdraw(30)
    println("balance: ${a.balance()}")
}
