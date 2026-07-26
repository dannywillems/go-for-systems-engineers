fun main() {
    println("floatBits(1.0f) = 0x%08X (ByteBuffer)".format(floatBits(1.0f)))
    println("readIntLE([1,0,0,0]) = ${readIntLE(byteArrayOf(1, 0, 0, 0))}")
}
