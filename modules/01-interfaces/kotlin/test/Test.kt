import kotlin.math.abs

fun main() {
    val shapes: List<Shape> = listOf(Circle(1.0), Circle(2.0))
    val expected = Math.PI * 1.0 + Math.PI * 4.0
    check(abs(sumShapes(shapes) - expected) < 1e-9) { "sum wrong" }
    println("ok")
}
