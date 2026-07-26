// Red/green gate: the builders agree, and the pre-sized StringBuilder allocates
// strictly fewer bytes than String +=.

fun main() {
    val parts = (0 until 64).map { "chunk-%02d;".format(it) }
    check(concatPlus(parts) == builder(parts)) { "builders disagree" }
    val plusBytes = bytesAllocated { concatPlus(parts) }
    val builderBytes = bytesAllocated { builder(parts) }
    check(builderBytes < plusBytes) {
        "builder ($builderBytes B) should allocate less than concatPlus ($plusBytes B)"
    }
    println("ok")
}
