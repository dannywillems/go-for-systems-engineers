import Errs

let a = (try? chain(3)) ?? -1
print("try? chain(3) = \(a)")
print("chainResult(3) = \(chainResult(3))")
print("chainResult(60) = \(chainResult(60))")
