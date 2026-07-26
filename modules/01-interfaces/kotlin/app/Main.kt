fun main() {
    val shapes: List<Shape> = listOf(Circle(1.0), Square(2.0))
    println("a Circle used as a Shape has area %.4f".format(shapes[0].area()))
    println("sum via interface dispatch = %.4f".format(sumShapes(shapes)))
}
