// Prints the two deterministic facts Module 00 checks. stdout is injected into
// the module README verbatim by the capture tool.

import Bootstrap

let n = 1_000_000
print("sum(1..\(n)) = \(sum(n))")
print("word size (bytes) = \(wordSizeBytes())")
