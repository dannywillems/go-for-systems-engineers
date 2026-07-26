// DOES NOT COMPILE under the Swift 6 language mode: a mutable global is
// "nonisolated global shared mutable state", and reading/writing it from a
// concurrent task is a data-race-safety error. This is the direct analog of the
// Go RacyInc counter, which Go accepts. Compile with `swiftc -swift-version 6`.

var counter = 0

func bump() {
  Task.detached {
    counter += 1
  }
}
