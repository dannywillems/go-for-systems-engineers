// Module 10 in Kotlin/JVM: measuring allocation with the JVM's own accounting.
// com.sun.management.ThreadMXBean.getThreadAllocatedBytes reports the bytes a
// thread has allocated, so the delta across an operation is what it allocated --
// the JVM parallel to Go's AllocsPerRun and Rust's counting allocator.

import com.sun.management.ThreadMXBean
import java.lang.management.ManagementFactory

// concatPlus builds with String += in a loop. JVM strings are immutable, so
// each += allocates a new String and copies: O(n) allocations, O(n^2) bytes.
fun concatPlus(parts: List<String>): String {
    var s = ""
    for (p in parts) s += p
    return s
}

// builder pre-sizes a StringBuilder, so the build is a couple of allocations.
fun builder(parts: List<String>): String {
    val total = parts.sumOf { it.length }
    val sb = StringBuilder(total)
    for (p in parts) sb.append(p)
    return sb.toString()
}

private val bean = ManagementFactory.getThreadMXBean() as ThreadMXBean

// bytesAllocated returns the bytes f allocates on the current thread.
fun bytesAllocated(f: () -> Unit): Long {
    val id = Thread.currentThread().threadId()
    val before = bean.getThreadAllocatedBytes(id)
    f()
    return bean.getThreadAllocatedBytes(id) - before
}
