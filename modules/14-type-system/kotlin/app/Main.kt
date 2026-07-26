fun main() {
    // Covariance: a Producer<Cat> flows where a Producer<Animal> is expected.
    println("firstName(catProducer) = ${firstName(catProducer)} (covariance via out)")
}
