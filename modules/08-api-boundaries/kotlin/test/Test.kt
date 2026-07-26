fun main() {
    val a = Account.open(100)
    a.deposit(50)
    check(a.balance() == 150L) { "balance ${a.balance()}" }
    var threw = false
    try {
        a.withdraw(200)
    } catch (e: Overdraft) {
        threw = true
    }
    check(threw) { "withdraw over balance should throw" }
    check(a.balance() == 150L)
    println("ok")
}
