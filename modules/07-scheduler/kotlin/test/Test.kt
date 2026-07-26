import kotlin.math.abs
import kotlin.math.sqrt

fun main() {
    val expected = 0.0 + 1.0 + sqrt(2.0) + sqrt(3.0)
    check(abs(chunkSum(0, 4) - expected) < 1e-9) { "chunkSum wrong" }
    check(abs(parallelSqrtSum(400, 4) - chunkSum(0, 400)) < 1e-6) { "parallel wrong" }
    println("ok")
}
