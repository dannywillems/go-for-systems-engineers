// The JVM is managed memory: there is no unsafe.Pointer. The historical
// `sun.misc.Unsafe` is deprecated and access-restricted; the sanctioned modern
// path is the Foreign Function & Memory API (java.lang.foreign, JEP 454, stable
// in JDK 22+). For structured reinterpretation of bytes without third-party
// deps, ByteBuffer is the portable tool, shown here reading raw bits.

import java.nio.ByteBuffer
import java.nio.ByteOrder

fun floatBits(f: Float): Int =
    ByteBuffer
        .allocate(4)
        .order(ByteOrder.LITTLE_ENDIAN)
        .putFloat(f)
        .getInt(0)

fun readIntLE(bytes: ByteArray): Int = ByteBuffer.wrap(bytes).order(ByteOrder.LITTLE_ENDIAN).int
