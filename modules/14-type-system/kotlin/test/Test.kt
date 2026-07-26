fun main() {
    check(firstName(catProducer) == "cat") { "covariance failed" }
    println("ok")
}
