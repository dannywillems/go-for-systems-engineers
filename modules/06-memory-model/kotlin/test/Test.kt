fun main() {
    check(atomicCount(8, 100000) == 800000) { "atomicCount wrong" }
    check(atomicCount(4, 1000) == 4000) { "atomicCount small wrong" }
    println("ok")
}
