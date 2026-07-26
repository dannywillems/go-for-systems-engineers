fun main() {
    check(chain(3) == 7) { "chain(3) wrong" }
    check(chainResult(3).getOrNull() == 7) { "chainResult(3) wrong" }
    check(chainResult(60).isFailure) { "chainResult(60) should fail" }
    println("ok")
}
