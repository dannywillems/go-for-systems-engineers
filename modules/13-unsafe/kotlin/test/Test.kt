fun main() {
    check(floatBits(1.0f) == 0x3F800000) { "floatBits wrong" }
    check(readIntLE(byteArrayOf(1, 0, 0, 0)) == 1) { "readIntLE wrong" }
    println("ok")
}
