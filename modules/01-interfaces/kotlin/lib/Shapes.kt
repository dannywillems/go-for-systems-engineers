// Kotlin interfaces are NOMINAL (a class must declare `: Shape`) and dispatch is
// virtual on the JVM (invokeinterface). There is no structural interface
// satisfaction (except SAM / `fun interface` conversion for single-method
// interfaces), and extension functions cannot retroactively make a foreign
// class implement an interface, so the Go-style structural coercion is absent.

// region:iface:start

interface Shape {
    fun area(): Double
}

class Circle(
    private val r: Double,
) : Shape {
    override fun area(): Double = Math.PI * r * r
}

class Square(
    private val s: Double,
) : Shape {
    override fun area(): Double = s * s
}

// region:iface:end

fun sumShapes(xs: List<Shape>): Double = xs.sumOf { it.area() }
