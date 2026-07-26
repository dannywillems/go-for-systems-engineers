// The capstone cache in Kotlin/JVM: a synchronized HashMap over a slow backend,
// with a java.util.concurrent.Semaphore for backpressure and platform threads
// for the load. LockSupport.parkNanos gives the sub-millisecond backend latency.

import java.util.concurrent.Semaphore
import java.util.concurrent.atomic.AtomicLong
import java.util.concurrent.locks.LockSupport

class Cache(
    private val capacity: Int,
    maxInflight: Int,
    private val latencyMicros: Long,
) {
    private val entries = HashMap<Int, Int>(capacity)
    private val lock = Any()
    private val sem = Semaphore(maxInflight)
    private val backend = AtomicLong(0)

    fun get(key: Int): Int {
        synchronized(lock) { entries[key]?.let { return it } }
        sem.acquire()
        backend.incrementAndGet()
        LockSupport.parkNanos(latencyMicros * 1000)
        val v = key * key
        sem.release()
        synchronized(lock) {
            if (entries.size >= capacity) entries.remove(entries.keys.first())
            entries[key] = v
        }
        return v
    }

    fun backendCalls(): Long = backend.get()

    fun size(): Int = synchronized(lock) { entries.size }
}
